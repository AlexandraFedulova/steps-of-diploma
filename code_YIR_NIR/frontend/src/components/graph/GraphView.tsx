import cytoscape, { type Core } from "cytoscape";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import type { GraphPayload } from "@/types";
import { applyCommunitySeedPositions, buildCommunityLegend } from "@/utils/communityVisual";
import { buildForceDirectedLayout } from "@/utils/cytoscapeLayout";
import {
  bindGraphInteractivity,
  configureGraphInteraction,
} from "@/utils/graphInteractivity";
import { graphToCytoscapeElements } from "@/utils/graphToCytoscape";
import { formatNodeTooltip, getNodeInfo, type NodeInfo } from "@/utils/nodeInfo";
import { CommunityLegend } from "./CommunityLegend";
import { GraphControls } from "./GraphControls";
import { GraphTooltip } from "./GraphTooltip";
import { NodeInfoPanel } from "./NodeInfoPanel";
import { registerCytoscapeExtensions } from "./cytoscapeSetup";
import { buildGraphStylesheet } from "./graphStylesheet";
import styles from "./GraphView.module.css";

export interface GraphViewProps {
  graph: GraphPayload;
}

interface TooltipState {
  x: number;
  y: number;
  text: string;
}

export function GraphView({ graph }: GraphViewProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const canvasWrapRef = useRef<HTMLDivElement>(null);
  const cyRef = useRef<Core | null>(null);
  const communityLegend = useMemo(() => buildCommunityLegend(graph), [graph]);

  const [tooltip, setTooltip] = useState<TooltipState | null>(null);
  const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null);

  const selectedNodeInfo = useMemo((): NodeInfo | null => {
    if (!selectedNodeId) {
      return null;
    }
    return getNodeInfo(graph, selectedNodeId);
  }, [graph, selectedNodeId]);

  const handleNodeHover = useCallback(
    (nodeId: string | null, position?: { x: number; y: number }) => {
      if (!nodeId || !position) {
        setTooltip(null);
        return;
      }
      const info = getNodeInfo(graph, nodeId);
      if (!info) {
        setTooltip(null);
        return;
      }
      setTooltip({
        x: position.x,
        y: position.y,
        text: formatNodeTooltip(info),
      });
    },
    [graph],
  );

  const handleNodeSelect = useCallback((nodeId: string | null) => {
    setSelectedNodeId(nodeId);
  }, []);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) {
      return;
    }

    registerCytoscapeExtensions();

    const cy = cytoscape({
      container,
      elements: graphToCytoscapeElements(graph),
      style: buildGraphStylesheet(graph.directed),
      minZoom: 0.15,
      maxZoom: 4,
      wheelSensitivity: 0.25,
    });

    cyRef.current = cy;
    configureGraphInteraction(cy);

    const unbind = bindGraphInteractivity(cy, {
      onNodeHover: handleNodeHover,
      onNodeSelect: handleNodeSelect,
    });

    applyCommunitySeedPositions(cy, graph);

    const layoutOptions = buildForceDirectedLayout(graph);
    const layout = cy.layout(layoutOptions);

    const onLayoutStop = () => {
      cy.fit(undefined, 48);
    };

    layout.on("layoutstop", onLayoutStop);
    layout.run();

    return () => {
      unbind();
      layout.stop();
      layout.removeListener("layoutstop", onLayoutStop);
      cy.destroy();
      cyRef.current = null;
      setTooltip(null);
      setSelectedNodeId(null);
    };
  }, [graph, handleNodeHover, handleNodeSelect]);

  useEffect(() => {
    const wrap = canvasWrapRef.current;
    if (!wrap) {
      return;
    }

    const observer = new ResizeObserver(() => {
      const cy = cyRef.current;
      if (!cy) {
        return;
      }
      cy.resize();
    });

    observer.observe(wrap);
    return () => observer.disconnect();
  }, []);

  return (
    <div className={styles.wrapper}>
      <div className={styles.toolbar}>
        <span>
          Узлов: <strong>{graph.nodeCount}</strong>
        </span>
        <span>
          Рёбер: <strong>{graph.edgeCount}</strong>
        </span>
        <span className={styles.layoutHint}>
          На узлах — номера · клик — имя респондента и подробности · перетаскивание · колёсико и фон — навигация
        </span>
        <CommunityLegend items={communityLegend} />
        <div className={styles.legend} aria-hidden="true">
          <span className={styles.legendItem}>
            <span className={styles.legendNodeSmall} />
            низкая centrality
          </span>
          <span className={styles.legendItem}>
            <span className={styles.legendNodeLarge} />
            высокая centrality
          </span>
          <span className={styles.legendItem}>
            <span className={styles.legendLinePositive} />
            положительная связь
          </span>
          <span className={styles.legendItem}>
            <span className={styles.legendLineNegative} />
            отрицательная связь
          </span>
        </div>
      </div>

      <div className={styles.graphArea}>
        <div ref={canvasWrapRef} className={styles.canvasWrap}>
          <GraphControls cyRef={cyRef} />
          <div ref={containerRef} className={styles.container} role="application" aria-label="Интерактивная социограмма" />
          {tooltip && <GraphTooltip x={tooltip.x} y={tooltip.y} text={tooltip.text} />}
        </div>
        {selectedNodeInfo && (
          <NodeInfoPanel info={selectedNodeInfo} onClose={() => setSelectedNodeId(null)} />
        )}
      </div>
    </div>
  );
}
