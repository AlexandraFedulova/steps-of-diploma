import { CsvFileUpload } from "@/components/upload/CsvFileUpload";
import { useGraph } from "@/context/GraphContext";
import { useGraphUpload } from "@/hooks/useGraphUpload";
import { formatYesNo } from "@/utils/display";
import styles from "./UploadPage.module.css";

export function UploadPage() {
  const { file, loading, error, setFile } = useGraphUpload({ openGraphOnSuccess: true });
  const { graph } = useGraph();

  return (
    <article className={styles.page}>
      <h2 className={styles.title}>Загрузка CSV</h2>
      <p className={styles.lead}>
        Выберите или перетащите файл социоматрицы. Файл отправится на сервер для анализа; после
        успешного ответа откроется визуализация графа.
      </p>

      <CsvFileUpload file={file} onFileChange={setFile} disabled={loading} />

      {loading && (
        <p className={styles.status} role="status">
          Анализ на сервере…
        </p>
      )}

      {error && (
        <p className={styles.statusError} role="alert">
          {error}
        </p>
      )}

      {file && !loading && !error && graph && (
        <section className={styles.graphSummary} aria-live="polite">
          <h3 className={styles.graphTitle}>Анализ завершён</h3>
          <ul className={styles.graphStats}>
            <li>
              Файл: <strong>{file.name}</strong>
            </li>
            <li>
              Узлов: <strong>{graph.nodeCount}</strong>
            </li>
            <li>
              Рёбер: <strong>{graph.edgeCount}</strong>
            </li>
            <li>
              Сообществ (клики): <strong>{countCommunities(graph)}</strong>
            </li>
            <li>
              Ориентированный: <strong>{formatYesNo(graph.directed)}</strong>
            </li>
          </ul>
          <p className={styles.graphHint}>Переход на страницу графа…</p>
        </section>
      )}
    </article>
  );
}

function countCommunities(graph: { community: { communityId: number }[] }): number {
  return new Set(graph.community.map((c) => c.communityId)).size;
}
