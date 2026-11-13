import cytoscape, { type EdgeSingular, type NodeSingular } from "cytoscape";
import { fontSizeFromNorm, nodeSizeFromNorm } from "@/utils/centralityVisual";

export function buildGraphStylesheet(directed: boolean): cytoscape.StylesheetStyle[] {
  return [
    {
      selector: "node",
      style: {
        label: "data(displayLabel)",
        "text-valign": "center",
        "text-halign": "center",
        "font-size": (ele: NodeSingular) =>
          fontSizeFromNorm((ele.data("centralityNorm") as number) ?? 0),
        color: "#0f172a",
        "text-outline-width": 2,
        "text-outline-color": "#ffffff",
        width: (ele: NodeSingular) =>
          nodeSizeFromNorm((ele.data("centralityNorm") as number) ?? 0),
        height: (ele: NodeSingular) =>
          nodeSizeFromNorm((ele.data("centralityNorm") as number) ?? 0),
        "background-color": "data(communityFill)",
        "border-width": (ele: NodeSingular) => {
          const norm = (ele.data("centralityNorm") as number) ?? 0;
          return norm > 0.65 ? 3 : 2.5;
        },
        "border-color": (ele: NodeSingular) => {
          const fill = ele.data("communityFill") as string;
          const norm = (ele.data("centralityNorm") as number) ?? 0;
          if (norm > 0.65) {
            return fill;
          }
          return "#ffffff";
        },
        "z-index": (ele: NodeSingular) =>
          Math.round(((ele.data("centralityNorm") as number) ?? 0) * 100),
      },
    },
    {
      selector: "node.hovered",
      style: {
        "border-width": 4,
        "overlay-color": "#2563eb",
        "overlay-opacity": 0.12,
        "overlay-padding": 6,
      },
    },
    {
      selector: "node:selected",
      style: {
        "border-width": 4,
        "border-color": "#1d4ed8",
        "overlay-color": "#2563eb",
        "overlay-opacity": 0.18,
        "overlay-padding": 8,
      },
    },
    {
      selector: "node:active",
      style: {
        "overlay-opacity": 0.25,
      },
    },
    {
      selector: "edge",
      style: {
        width: (ele: EdgeSingular) =>
          Math.max(1, Math.min(6, Math.abs(ele.data("weight") as number))),
        "line-color": (ele: EdgeSingular) => {
          const w = ele.data("weight") as number;
          return w < 0 ? "#dc2626" : "#16a34a";
        },
        "target-arrow-color": (ele: EdgeSingular) => {
          const w = ele.data("weight") as number;
          return w < 0 ? "#dc2626" : "#16a34a";
        },
        "target-arrow-shape": directed ? "triangle" : "none",
        "curve-style": "bezier",
        label: "data(label)",
        "font-size": "9px",
        color: "#475569",
        "text-rotation": "autorotate",
        "text-margin-y": -8,
      },
    },
  ];
}
