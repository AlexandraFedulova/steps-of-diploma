/** Узел графа */
export interface GraphNode {
  id: string;
  label: string;
}

/** Ребро графа (ориентированное) */
export interface GraphEdge {
  source: string;
  target: string;
  weight: number;
}

/** Degree centrality по входящим взвешенным связям */
export interface CentralityScore {
  nodeId: string;
  inWeightSum: number;
  value: number;
}

/** Назначение вершины в сообщество (клика взаимного выбора) */
export interface NodeCommunity {
  nodeId: string;
  communityId: number;
}

/** Матрица весов социоматрицы */
export interface WeightMatrix {
  labels: string[];
  values: number[][];
}

/** Полный payload графа для визуализации */
export interface GraphPayload {
  directed: boolean;
  nodeCount: number;
  edgeCount: number;
  nodes: GraphNode[];
  edges: GraphEdge[];
  centrality: CentralityScore[];
  community: NodeCommunity[];
  weights: WeightMatrix;
}

/** Ответ POST /upload */
export interface UploadGraphResponse {
  ok: boolean;
  filename: string;
  size: number;
  graph: GraphPayload;
  error?: string;
}

/** Ответ health-check */
export interface HealthResponse {
  status: "ok";
  service: string;
}
