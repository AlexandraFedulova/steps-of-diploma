import type { Core } from "cytoscape";
import type { GraphPayload } from "@/types";

/** Ключ ориентированного ребра. */
export function directedEdgeKey(source: string, target: string): string {
  return `${source}\t${target}`;
}

/** Ключ неориентированной пары (для взаимного выбора). */
export function undirectedPairKey(a: string, b: string): string {
  return a < b ? `${a}\t${b}` : `${b}\t${a}`;
}

/** Исходящие связи с ненулевым весом. */
export function buildOutgoingMap(graph: GraphPayload): Map<string, Set<string>> {
  const map = new Map<string, Set<string>>();
  for (const { source, target, weight } of graph.edges) {
    if (weight === 0) {
      continue;
    }
    if (!map.has(source)) {
      map.set(source, new Set());
    }
    map.get(source)!.add(target);
  }
  return map;
}

/** Пара респондентов взаимно выбирает друг друга (есть A→B и B→A, вес ≠ 0). */
export function isMutualPair(
  a: string,
  b: string,
  outgoing: Map<string, Set<string>>,
): boolean {
  return Boolean(outgoing.get(a)?.has(b) && outgoing.get(b)?.has(a));
}

/** Пары с взаимным выбором. */
export function buildMutualPairKeys(graph: GraphPayload): Set<string> {
  const outgoing = buildOutgoingMap(graph);
  const pairs = new Set<string>();
  const seen = new Set<string>();

  for (const { source, target } of graph.edges) {
    const pair = undirectedPairKey(source, target);
    if (seen.has(pair)) {
      continue;
    }
    seen.add(pair);
    if (isMutualPair(source, target, outgoing)) {
      pairs.add(pair);
    }
  }
  return pairs;
}

/** Компоненты связности по взаимным выборам (только группы из 2+ узлов). */
export function buildMutualComponents(graph: GraphPayload): string[][] {
  const outgoing = buildOutgoingMap(graph);
  const parent = new Map<string, string>();

  for (const node of graph.nodes) {
    parent.set(node.id, node.id);
  }

  const find = (x: string): string => {
    let root = parent.get(x)!;
    while (root !== parent.get(root)) {
      root = parent.get(root)!;
    }
    let current = x;
    while (parent.get(current) !== root) {
      const next = parent.get(current)!;
      parent.set(current, root);
      current = next;
    }
    return root;
  };

  const union = (a: string, b: string): void => {
    const ra = find(a);
    const rb = find(b);
    if (ra !== rb) {
      parent.set(ra, rb);
    }
  };

  const mergedPairs = new Set<string>();
  for (const { source, target } of graph.edges) {
    const pair = undirectedPairKey(source, target);
    if (mergedPairs.has(pair)) {
      continue;
    }
    mergedPairs.add(pair);
    if (isMutualPair(source, target, outgoing)) {
      union(source, target);
    }
  }

  const groups = new Map<string, string[]>();
  for (const node of graph.nodes) {
    const root = find(node.id);
    if (!groups.has(root)) {
      groups.set(root, []);
    }
    groups.get(root)!.push(node.id);
  }

  return [...groups.values()].filter((group) => group.length >= 2);
}

/** Упорядочить узлы: блоки взаимных кластеров рядом, затем одиночки. */
export function orderNodesByMutualClusters(
  nodeIds: string[],
  graph: GraphPayload,
): string[] {
  const components = buildMutualComponents(graph).filter((group) =>
    group.every((id) => nodeIds.includes(id)),
  );
  const clustered = new Set(components.flat());
  const singletons = nodeIds.filter((id) => !clustered.has(id));
  return [...components.flat(), ...singletons];
}

/**
 * Сжать взаимные кластеры в компактные группы (старт fcose).
 * Радиус подбирается по числу узлов, чтобы не накладывались.
 */
export function applyMutualClusterSeedPositions(cy: Core, graph: GraphPayload): void {
  const components = buildMutualComponents(graph);
  if (components.length === 0) {
    return;
  }

  for (const nodeIds of components) {
    if (nodeIds.length < 2) {
      continue;
    }

    let sumX = 0;
    let sumY = 0;
    for (const id of nodeIds) {
      const pos = cy.getElementById(id).position();
      sumX += pos.x;
      sumY += pos.y;
    }
    const centerX = sumX / nodeIds.length;
    const centerY = sumY / nodeIds.length;

    const spacing = 58 + Math.max(0, nodeIds.length - 2) * 14;
    const radius =
      nodeIds.length === 2 ? spacing * 0.55 : spacing * (0.45 + nodeIds.length * 0.06);

    nodeIds.forEach((id, index) => {
      const angle = (2 * Math.PI * index) / nodeIds.length - Math.PI / 2;
      cy.getElementById(id).position({
        x: centerX + radius * Math.cos(angle),
        y: centerY + radius * Math.sin(angle),
      });
    });
  }
}
