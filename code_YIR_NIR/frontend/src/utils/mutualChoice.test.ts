import { describe, expect, it } from "vitest";
import type { GraphPayload } from "@/types";
import {
  buildMutualComponents,
  buildMutualPairKeys,
  isMutualPair,
  orderNodesByMutualClusters,
} from "./mutualChoice";

function miniGraph(
  edges: Array<{ source: string; target: string; weight: number }>,
): GraphPayload {
  const ids = [...new Set(edges.flatMap((e) => [e.source, e.target]))];
  return {
    directed: true,
    nodeCount: ids.length,
    edgeCount: edges.length,
    nodes: ids.map((id) => ({ id, label: id })),
    edges,
    centrality: ids.map((id) => ({ nodeId: id, inWeightSum: 1, value: 0.5 })),
    community: ids.map((id) => ({ nodeId: id, communityId: 0 })),
    weights: { labels: ids, values: ids.map(() => ids.map(() => 0)) },
  };
}

describe("mutualChoice", () => {
  it("detects mutual pairs", () => {
    const graph = miniGraph([
      { source: "A", target: "B", weight: 1 },
      { source: "B", target: "A", weight: 1 },
      { source: "A", target: "C", weight: 1 },
    ]);
    const outgoing = new Map([
      ["A", new Set(["B", "C"])],
      ["B", new Set(["A"])],
      ["C", new Set<string>()],
    ]);
    expect(isMutualPair("A", "B", outgoing)).toBe(true);
    expect(isMutualPair("A", "C", outgoing)).toBe(false);
    expect(buildMutualPairKeys(graph).size).toBe(1);
  });

  it("merges transitive mutual clusters", () => {
    const graph = miniGraph([
      { source: "A", target: "B", weight: 1 },
      { source: "B", target: "A", weight: 1 },
      { source: "B", target: "C", weight: 1 },
      { source: "C", target: "B", weight: 1 },
    ]);
    const components = buildMutualComponents(graph);
    expect(components).toHaveLength(1);
    expect(components[0]).toEqual(expect.arrayContaining(["A", "B", "C"]));
  });

  it("orders mutual cluster as a block", () => {
    const graph = miniGraph([
      { source: "A", target: "B", weight: 1 },
      { source: "B", target: "A", weight: 1 },
    ]);
    const ordered = orderNodesByMutualClusters(["C", "A", "B"], graph);
    expect(ordered.indexOf("A")).toBeLessThan(ordered.indexOf("C"));
    expect(ordered.indexOf("B")).toBeLessThan(ordered.indexOf("C"));
  });
});
