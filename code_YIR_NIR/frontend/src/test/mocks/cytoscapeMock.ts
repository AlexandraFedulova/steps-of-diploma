import { vi } from "vitest";
import type { EventObject, NodeSingular } from "cytoscape";

export type CytoscapeEventHandler = (event: EventObject) => void;

export interface MockNode {
  id: string;
  data: Record<string, unknown>;
  classes: Set<string>;
  selected: boolean;
  position: { x: number; y: number };
}

export interface MockLayout {
  on: ReturnType<typeof vi.fn>;
  run: ReturnType<typeof vi.fn>;
  stop: ReturnType<typeof vi.fn>;
  removeListener: ReturnType<typeof vi.fn>;
}

export interface MockCy {
  _nodes: MockNode[];
  _handlers: Map<string, CytoscapeEventHandler[]>;
  _destroyed: boolean;
  trigger: (eventName: string, targetId?: string) => void;
  destroy: ReturnType<typeof vi.fn>;
  layout: ReturnType<typeof vi.fn>;
  fit: ReturnType<typeof vi.fn>;
}

function wrapNode(node: MockNode): NodeSingular {
  return {
    id: () => node.id,
    data: (key?: string) => (key === undefined ? node.data : node.data[key]),
    addClass: (name: string) => {
      node.classes.add(name);
    },
    removeClass: (name: string) => {
      node.classes.delete(name);
    },
    hasClass: (name: string) => node.classes.has(name),
    select: () => {
      node.selected = true;
    },
    unselect: () => {
      node.selected = false;
    },
    renderedPosition: () => ({ x: 120, y: 80 }),
    position: (pos?: { x: number; y: number }) => {
      if (pos) {
        node.position = pos;
      }
      return node.position;
    },
  } as unknown as NodeSingular;
}

export function createMockCy(nodes: MockNode[]): MockCy {
  const handlers = new Map<string, CytoscapeEventHandler[]>();

  const getHandlers = (key: string) => {
    if (!handlers.has(key)) {
      handlers.set(key, []);
    }
    return handlers.get(key)!;
  };

  const emit = (key: string, event: EventObject) => {
    for (const fn of [...getHandlers(key)]) {
      fn(event);
    }
  };

  const cySelf = {} as MockCy & Record<string, unknown>;

  const layout: MockLayout = {
    on: vi.fn(),
    run: vi.fn(() => {
      const stopHandler = layout.on.mock.calls.find(([e]) => e === "layoutstop")?.[1] as
        | (() => void)
        | undefined;
      stopHandler?.();
    }),
    stop: vi.fn(),
    removeListener: vi.fn(),
  };

  Object.assign(cySelf, {
    _nodes: nodes,
    _handlers: handlers,
    _destroyed: false,

    on(event: string, selectorOrHandler: string | CytoscapeEventHandler, maybeHandler?: CytoscapeEventHandler) {
      const selector = typeof selectorOrHandler === "string" ? selectorOrHandler : undefined;
      const handler =
        typeof selectorOrHandler === "function" ? selectorOrHandler : maybeHandler!;
      const key = selector ? `${event} ${selector}` : event;
      getHandlers(key).push(handler);
    },

    off(event: string, selectorOrHandler: string | CytoscapeEventHandler, maybeHandler?: CytoscapeEventHandler) {
      const selector = typeof selectorOrHandler === "string" ? selectorOrHandler : undefined;
      const handler =
        typeof selectorOrHandler === "function" ? selectorOrHandler : maybeHandler!;
      const key = selector ? `${event} ${selector}` : event;
      const list = getHandlers(key);
      const idx = list.indexOf(handler);
      if (idx >= 0) {
        list.splice(idx, 1);
      }
    },

    nodes(selector?: string) {
      const collection = {
        length: nodes.length,
        forEach(fn: (n: NodeSingular, i: number) => void) {
          nodes.forEach((raw, i) => fn(wrapNode(raw), i));
        },
        filter(predicate: (n: NodeSingular) => boolean) {
          const filtered = nodes.filter((raw) => predicate(wrapNode(raw)));
          return {
            length: filtered.length,
            forEach(fn: (n: NodeSingular, i: number) => void) {
              filtered.forEach((raw, i) => fn(wrapNode(raw), i));
            },
          };
        },
        unselect() {
          nodes.forEach((n) => {
            n.selected = false;
          });
        },
        grabify: vi.fn(),
        removeClass(name: string) {
          nodes.forEach((n) => n.classes.delete(name));
        },
      };
      if (selector === ".hovered") {
        const hovered = nodes.filter((n) => n.classes.has("hovered"));
        return {
          ...collection,
          length: hovered.length,
          removeClass(name: string) {
            hovered.forEach((n) => n.classes.delete(name));
          },
        };
      }
      return collection;
    },

    layout: vi.fn(() => layout),
    fit: vi.fn(),
    resize: vi.fn(),
    zoom: vi.fn(),
    width: vi.fn(() => 800),
    height: vi.fn(() => 600),
    destroy: vi.fn(() => {
      cySelf._destroyed = true;
    }),
    userPanningEnabled: vi.fn(),
    userZoomingEnabled: vi.fn(),
    boxSelectionEnabled: vi.fn(),
    autoungrabify: vi.fn(),
    autounselectify: vi.fn(),

    trigger(eventName: string, targetId?: string) {
      const node = targetId ? nodes.find((n) => n.id === targetId) : undefined;
      const target = node ? wrapNode(node) : cySelf;
      const event = { target } as EventObject;

      if (eventName === "mouseover" && node) {
        emit("mouseover node", event);
      } else if (eventName === "mouseout" && node) {
        emit("mouseout node", event);
      } else if (eventName === "tap" && node) {
        emit("tap node", event);
      } else if (eventName === "tap") {
        emit("tap", event);
      }
    },
  });

  return cySelf;
}

export const cytoscapeConstructorMock = vi.fn();
