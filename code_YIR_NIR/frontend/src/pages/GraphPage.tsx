import { Link } from "react-router-dom";
import { GraphView } from "@/components/graph/GraphView";
import { useGraph } from "@/context/GraphContext";
import styles from "./GraphPage.module.css";

export function GraphPage() {
  const { graph } = useGraph();

  return (
    <article className={styles.page}>
      <header className={styles.header}>
        <h2 className={styles.title}>Граф</h2>
        <p className={styles.lead}>
          {graph
            ? `Визуализация: ${graph.nodeCount} узлов, ${graph.edgeCount} рёбер. Centrality и сообщества по кликам взаимного выбора получены при загрузке CSV.`
            : "Визуализация социоматрицы: вершины — участники, рёбра — выборы."}
        </p>
        {graph && (
          <p className={styles.meta}>
            <Link to="/upload">Загрузить другой файл</Link>
          </p>
        )}
      </header>

      {graph ? (
        <GraphView graph={graph} />
      ) : (
        <p className={styles.empty}>
          Граф ещё не загружен.{" "}
          <Link to="/upload">Загрузите CSV</Link> — после анализа на сервере граф отобразится здесь
          автоматически.
        </p>
      )}
    </article>
  );
}
