import type { EdgeSingular, LayoutOptions, NodeSingular } from "cytoscape";
import type { GraphPayload } from "@/types";
import { sameCommunity } from "@/utils/communityVisual";
import { nodeRepulsionWithCentrality } from "@/utils/centralityVisual";

/**
 * Force-directed раскладка (fcose): связанные узлы ближе, но с достаточным зазором между кружками.
 */
export function buildForceDirectedLayout(graph: GraphPayload): LayoutOptions {
  const n = Math.max(graph.nodeCount, 1);
  const scale = Math.sqrt(n);
  const baseRepulsion = 9000 + scale * 1500;

  return {
    name: "fcose",
    quality: "default",
    randomize: false,
    animate: false,
    fit: true,
    padding: 56,
    nodeDimensionsIncluded: true,
    packComponents: true,
    tile: true,
    tilingPaddingVertical: 72,
    tilingPaddingHorizontal: 72,
    nodeSeparation: 85 + scale * 14,
    nodeRepulsion: (node: NodeSingular) => nodeRepulsionWithCentrality(node, baseRepulsion),
    idealEdgeLength: (edge: EdgeSingular) => {
      const w = Math.abs(edge.data("weight") as number) || 1;
      let base = Math.max(52, 145 / Math.sqrt(w));
      if (sameCommunity(edge)) {
        base *= 0.88;
      }
      if (edge.data("mutual")) {
        // Слегка сближаем взаимные пары, но не «стягиваем» весь граф в ком.
        base *= 0.85;
      }
      return base;
    },
    edgeElasticity: (edge: EdgeSingular) => {
      const w = Math.abs(edge.data("weight") as number) || 1;
      const base = 0.1 + 0.08 * w;
      return Math.min(0.55, base);
    },
    nestingFactor: 0.1,
    gravity: 0.2,
    gravityRangeCompound: 1.2,
    gravityCompound: 0.85,
    numIter: Math.min(5000, 1400 + n * 110),
  } as LayoutOptions;
}

/** Запасная встроенная раскладка cose. */
export function buildCoseLayout(graph: GraphPayload): LayoutOptions {
  const n = Math.max(graph.nodeCount, 1);
  const scale = Math.sqrt(n);

  return {
    name: "cose",
    animate: false,
    fit: true,
    padding: 56,
    randomize: true,
    componentSpacing: 100 + scale * 28,
    nodeRepulsion: 9000 + scale * 1600,
    nodeOverlap: 24,
    idealEdgeLength: (edge: EdgeSingular) => {
      const w = Math.abs(edge.data("weight") as number) || 1;
      return Math.max(55, 150 / Math.sqrt(w));
    },
    edgeElasticity: (edge: EdgeSingular) => {
      const w = Math.abs(edge.data("weight") as number) || 1;
      return Math.min(0.55, 0.12 + 0.08 * w);
    },
    nestingFactor: 0.1,
    gravity: 0.18,
    numIter: Math.min(4000, 1100 + n * 90),
    initialTemp: 220,
    coolingFactor: 0.95,
    minTemp: 1.0,
  } as LayoutOptions;
}
