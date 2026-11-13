import type { GraphPayload } from "@/types";
import { stableCommunityColor } from "@/utils/communityVisual";
import { buildCentralityNormMap } from "@/utils/centralityVisual";
import { nodeNumber } from "@/utils/nodeNumbering";

export interface NodeInfo {
  id: string;
  label: string;
  number: number;
  communityId: number;
  communityColor: string;
  centrality: number;
  centralityNorm: number;
  inWeightSum: number;
  incomingCount: number;
  outgoingCount: number;
}

export function getNodeInfo(graph: GraphPayload, nodeId: string): NodeInfo | null {
  const node = graph.nodes.find((n) => n.id === nodeId);
  if (!node) {
    return null;
  }

  const centrality = graph.centrality.find((c) => c.nodeId === nodeId);
  const community = graph.community.find((c) => c.nodeId === nodeId);
  const incomingCount = graph.edges.filter((e) => e.target === nodeId).length;
  const outgoingCount = graph.edges.filter((e) => e.source === nodeId).length;

  const communityId = community?.communityId ?? 0;
  const centralityNorm = buildCentralityNormMap(graph).get(nodeId) ?? 0;

  return {
    id: node.id,
    label: node.label,
    number: nodeNumber(graph, nodeId),
    communityId,
    communityColor: stableCommunityColor(communityId),
    centrality: centrality?.value ?? 0,
    centralityNorm,
    inWeightSum: centrality?.inWeightSum ?? 0,
    incomingCount,
    outgoingCount,
  };
}

export function formatNodeTooltip(info: NodeInfo): string {
  return `№${info.number} · сообщество №${info.communityId} · centrality ${info.centrality.toFixed(3)}`;
}
