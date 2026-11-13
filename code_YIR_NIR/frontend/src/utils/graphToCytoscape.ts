import type { ElementDefinition } from "cytoscape";
import type { GraphPayload } from "@/types";
import { buildCommunityColorMap, stableCommunityColor } from "@/utils/communityVisual";
import { buildCentralityNormMap } from "@/utils/centralityVisual";
import { buildNodeNumberMap } from "@/utils/nodeNumbering";
import {
  buildMutualPairKeys,
  directedEdgeKey,
  undirectedPairKey,
} from "@/utils/mutualChoice";

export { stableCommunityColor as communityColor };

export function graphToCytoscapeElements(graph: GraphPayload): ElementDefinition[] {
  const communityByNode = new Map(graph.community.map((c) => [c.nodeId, c.communityId]));
  const communityColors = buildCommunityColorMap(graph);
  const centralityByNode = new Map(graph.centrality.map((c) => [c.nodeId, c.value]));
  const centralityNormByNode = buildCentralityNormMap(graph);
  const nodeNumbers = buildNodeNumberMap(graph);
  const mutualPairs = buildMutualPairKeys(graph);

  const nodes: ElementDefinition[] = graph.nodes.map((node) => {
    const communityId = communityByNode.get(node.id) ?? 0;
    const communityFill =
      communityColors.get(communityId) ?? stableCommunityColor(communityId);
    return {
      group: "nodes",
      data: {
        id: node.id,
        name: node.label,
        displayLabel: String(nodeNumbers.get(node.id) ?? ""),
        nodeNumber: nodeNumbers.get(node.id) ?? 0,
        communityId,
        communityFill,
        centrality: centralityByNode.get(node.id) ?? 0,
        centralityNorm: centralityNormByNode.get(node.id) ?? 0,
      },
    };
  });

  const edges: ElementDefinition[] = graph.edges.map((edge, index) => ({
    group: "edges",
    data: {
      id: `e-${edge.source}-${edge.target}-${index}`,
      source: edge.source,
      target: edge.target,
      weight: edge.weight,
      label: formatEdgeLabel(edge.weight),
      mutual: mutualPairs.has(undirectedPairKey(edge.source, edge.target)),
      edgeKey: directedEdgeKey(edge.source, edge.target),
    },
  }));

  return [...nodes, ...edges];
}

/** Подпись на ребре: только число при |весе| > 1 (знак — по цвету линии). */
export function formatEdgeLabel(weight: number): string {
  if (Math.abs(weight) <= 1) {
    return "";
  }
  const magnitude = Math.abs(weight);
  return Number.isInteger(magnitude) ? String(magnitude) : String(magnitude);
}
