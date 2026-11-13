import type cytoscape from "cytoscape";
import type { Core } from "cytoscape";
import type { GraphPayload } from "@/types";
import { orderNodesByMutualClusters } from "@/utils/mutualChoice";

/** Фиксированная палитра: communityId → всегда один и тот же цвет между рендерами. */
const COMMUNITY_PALETTE = [
  "#2563eb",
  "#16a34a",
  "#dc2626",
  "#9333ea",
  "#ea580c",
  "#0891b2",
  "#ca8a04",
  "#db2777",
  "#4f46e5",
  "#0d9488",
  "#b45309",
  "#7c3aed",
] as const;

export interface CommunityLegendItem {
  communityId: number;
  color: string;
  nodeCount: number;
}

/** Стабильный цвет по id сообщества (клика). */
export function stableCommunityColor(communityId: number): string {
  if (communityId >= 0 && communityId < COMMUNITY_PALETTE.length) {
    return COMMUNITY_PALETTE[communityId];
  }
  const hue = (communityId * 137.508) % 360;
  return `hsl(${hue} 62% 50%)`;
}

export function buildCommunityColorMap(graph: GraphPayload): Map<number, string> {
  const ids = [...new Set(graph.community.map((c) => c.communityId))].sort((a, b) => a - b);
  const map = new Map<number, string>();
  for (const id of ids) {
    map.set(id, stableCommunityColor(id));
  }
  return map;
}

export function buildCommunityLegend(graph: GraphPayload): CommunityLegendItem[] {
  const colorMap = buildCommunityColorMap(graph);
  const counts = new Map<number, number>();
  for (const { communityId } of graph.community) {
    counts.set(communityId, (counts.get(communityId) ?? 0) + 1);
  }

  return [...colorMap.entries()]
    .sort(([a], [b]) => a - b)
    .map(([communityId, color]) => ({
      communityId,
      color,
      nodeCount: counts.get(communityId) ?? 0,
    }));
}

/**
 * Начальные позиции: сектор на сообщество + радиус по centrality.
 * Узлы одного сообщества стартуют рядом → fcose визуально группирует кластеры.
 */
export function applyCommunitySeedPositions(cy: Core, graph: GraphPayload): void {
  const communityIds = [...new Set(graph.community.map((c) => c.communityId))].sort(
    (a, b) => a - b,
  );
  const communityCount = Math.max(communityIds.length, 1);
  const sectorAngle = (2 * Math.PI) / communityCount;
  const communityIndex = new Map(communityIds.map((id, index) => [id, index]));

  const nodesByCommunity = new Map<number, cytoscape.NodeCollection>();
  for (const id of communityIds) {
    nodesByCommunity.set(
      id,
      cy.nodes().filter((n) => (n.data("communityId") as number) === id),
    );
  }

  const totalNodes = cy.nodes().length;
  const baseRadius = 130 + totalNodes * 20;

  for (const [communityId, nodes] of nodesByCommunity) {
    const sectorIndex = communityIndex.get(communityId) ?? 0;
    const sectorCenter = sectorIndex * sectorAngle - Math.PI / 2;
    const spread = sectorAngle * 0.9;
    const count = nodes.length;

    const orderedIds = orderNodesByMutualClusters(
      nodes.map((n) => n.id()),
      graph,
    );
    const positionById = new Map(
      orderedIds.map((id, localIndex) => {
        const centralityNorm =
          (cy.getElementById(id).data("centralityNorm") as number) ?? 0;
        const radius = baseRadius + centralityNorm * 32;
        const angleOffset =
          count <= 1 ? 0 : (localIndex / (count - 1) - 0.5) * spread;
        const angle = sectorCenter + angleOffset;
        return [
          id,
          {
            x: radius * Math.cos(angle),
            y: radius * Math.sin(angle),
          },
        ] as const;
      }),
    );

    nodes.forEach((node) => {
      const pos = positionById.get(node.id());
      if (pos) {
        node.position(pos);
      }
    });
  }
}

export function sameCommunity(edge: {
  source: () => { data: (key: string) => unknown };
  target: () => { data: (key: string) => unknown };
}): boolean {
  return edge.source().data("communityId") === edge.target().data("communityId");
}
