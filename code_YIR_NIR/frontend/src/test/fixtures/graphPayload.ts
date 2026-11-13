import type { GraphPayload } from "@/types";

/** Минимальный граф для UI-тестов (3 вершины, 2 ребра). */
export const sampleGraphPayload: GraphPayload = {
  directed: true,
  nodeCount: 3,
  edgeCount: 2,
  nodes: [
    { id: "A", label: "Алиса" },
    { id: "B", label: "Борис" },
    { id: "C", label: "Вера" },
  ],
  edges: [
    { source: "A", target: "B", weight: 1 },
    { source: "B", target: "C", weight: -1 },
  ],
  centrality: [
    { nodeId: "A", inWeightSum: 0, value: 0 },
    { nodeId: "B", inWeightSum: 1, value: 0.5 },
    { nodeId: "C", inWeightSum: -1, value: -0.5 },
  ],
  community: [
    { nodeId: "A", communityId: 0 },
    { nodeId: "B", communityId: 0 },
    { nodeId: "C", communityId: 1 },
  ],
  weights: {
    labels: ["A", "B", "C"],
    values: [
      [0, 1, 0],
      [0, 0, -1],
      [0, 0, 0],
    ],
  },
};
