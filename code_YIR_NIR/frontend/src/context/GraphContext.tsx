import { createContext, useContext, useMemo, useState, type ReactNode } from "react";
import type { GraphPayload } from "@/types";

interface GraphContextValue {
  graph: GraphPayload | null;
  setGraph: (graph: GraphPayload | null) => void;
}

const GraphContext = createContext<GraphContextValue | null>(null);

export function GraphProvider({ children }: { children: ReactNode }) {
  const [graph, setGraph] = useState<GraphPayload | null>(null);

  const value = useMemo(() => ({ graph, setGraph }), [graph]);

  return <GraphContext.Provider value={value}>{children}</GraphContext.Provider>;
}

export function useGraph(): GraphContextValue {
  const ctx = useContext(GraphContext);
  if (!ctx) {
    throw new Error("useGraph должен использоваться внутри GraphProvider");
  }
  return ctx;
}
