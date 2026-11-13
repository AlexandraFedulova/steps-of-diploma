import type { GraphPayload } from "@/types";

export function validateGraphPayload(value: unknown): GraphPayload {
  if (!value || typeof value !== "object") {
    throw new Error("Некорректный формат графа в ответе сервера");
  }

  const graph = value as GraphPayload;

  if (!Array.isArray(graph.nodes) || graph.nodes.length === 0) {
    throw new Error("Граф не содержит узлов");
  }

  if (!Array.isArray(graph.edges)) {
    throw new Error("Граф не содержит рёбер");
  }

  if (typeof graph.directed !== "boolean") {
    throw new Error("Отсутствует поле directed в ответе");
  }

  if (typeof graph.nodeCount !== "number" || typeof graph.edgeCount !== "number") {
    throw new Error("Отсутствуют счётчики узлов и рёбер");
  }

  if (!Array.isArray(graph.centrality) || !Array.isArray(graph.community)) {
    throw new Error("Ответ не содержит результатов анализа (centrality, community)");
  }

  return graph;
}
