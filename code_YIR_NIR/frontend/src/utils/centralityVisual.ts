import type { Core, NodeSingular } from "cytoscape";
import type { GraphPayload } from "@/types";

export const NODE_SIZE_MIN = 28;
export const NODE_SIZE_MAX = 76;

/** Нормализованная centrality [0, 1] по всему графу (линейная min–max). */
export function buildCentralityNormMap(graph: GraphPayload): Map<string, number> {
  const values = graph.centrality.map((c) => c.value);
  const min = values.length > 0 ? Math.min(...values) : 0;
  const max = values.length > 0 ? Math.max(...values) : 1;
  const span = max - min || 1;

  const map = new Map<string, number>();
  for (const node of graph.nodes) {
    const raw = graph.centrality.find((c) => c.nodeId === node.id)?.value ?? 0;
    map.set(node.id, (raw - min) / span);
  }
  return map;
}

/** Размер кружка по нормализованной centrality (линейная шкала). */
export function nodeSizeFromNorm(norm: number): number {
  const clamped = Math.max(0, Math.min(1, norm));
  return NODE_SIZE_MIN + clamped * (NODE_SIZE_MAX - NODE_SIZE_MIN);
}

export function fontSizeFromNorm(norm: number): number {
  const clamped = Math.max(0, Math.min(1, norm));
  return 8 + clamped * 7;
}

/**
 * Стартовые позиции: кольцо с лёгким разбросом по centrality (без стяжки в центр).
 */
export function applyCentralitySeedPositions(cy: Core): void {
  const nodes = cy.nodes();
  const count = nodes.length;
  if (count === 0) {
    return;
  }

  const baseRadius = 140 + count * 22;

  nodes.forEach((node, index) => {
    const norm = (node.data("centralityNorm") as number) ?? 0;
    const radius = baseRadius + norm * 28;
    const angle = (2 * Math.PI * index) / count;
    node.position({
      x: radius * Math.cos(angle),
      y: radius * Math.sin(angle),
    });
  });
}

/** Чем выше centrality, тем сильнее отталкивание — хабы не слипаются в центре. */
export function nodeRepulsionWithCentrality(
  node: NodeSingular,
  baseRepulsion: number,
): number {
  const norm = (node.data("centralityNorm") as number) ?? 0;
  return baseRepulsion * (1.05 + norm * 0.2);
}
