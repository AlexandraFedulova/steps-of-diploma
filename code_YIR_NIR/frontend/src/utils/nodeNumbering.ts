import type { GraphPayload } from "@/types";

/** Порядковый номер вершины на графе (1…n) в порядке узлов из payload. */
export function buildNodeNumberMap(graph: GraphPayload): Map<string, number> {
  const map = new Map<string, number>();
  graph.nodes.forEach((node, index) => {
    map.set(node.id, index + 1);
  });
  return map;
}

export function nodeNumber(graph: GraphPayload, nodeId: string): number {
  const index = graph.nodes.findIndex((n) => n.id === nodeId);
  return index >= 0 ? index + 1 : 0;
}
