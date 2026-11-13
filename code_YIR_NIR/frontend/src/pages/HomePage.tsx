import { useHealth } from "@/hooks/useHealth";
import { formatHealthStatus } from "@/utils/display";
import styles from "./HomePage.module.css";

export function HomePage() {
  const { health, loading, error } = useHealth();

  return (
    <article className={styles.page}>
      <h2 className={styles.title}>Главная</h2>
      <p className={styles.lead}>
        Веб-приложение для анализа социометрических графов: загрузка матрицы, визуализация связей
        и поиск сообществ.
      </p>
      {loading && <p className={styles.status}>Проверка сервера…</p>}
      {!loading && error && (
        <p className={styles.statusError} role="alert">
          {error}
        </p>
      )}
      {!loading && health && (
        <p className={styles.statusOk}>
          Сервер: <strong>{formatHealthStatus(health.status)}</strong>
        </p>
      )}
    </article>
  );
}
