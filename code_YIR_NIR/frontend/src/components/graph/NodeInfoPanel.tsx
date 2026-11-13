import type { NodeInfo } from "@/utils/nodeInfo";
import styles from "./NodeInfoPanel.module.css";

interface NodeInfoPanelProps {
  info: NodeInfo;
  onClose: () => void;
}

export function NodeInfoPanel({ info, onClose }: NodeInfoPanelProps) {
  return (
    <aside className={styles.panel} aria-label="Информация о респонденте">
      <div className={styles.header}>
        <h3 className={styles.title}>{info.label}</h3>
        <button type="button" className={styles.close} onClick={onClose} aria-label="Закрыть">
          ×
        </button>
      </div>
      <dl className={styles.list}>
        <div className={styles.row}>
          <dt>№ на графе</dt>
          <dd>{info.number}</dd>
        </div>
        <div className={styles.row}>
          <dt>Сообщество</dt>
          <dd>
            <span className={styles.communityBadge} style={{ backgroundColor: info.communityColor }}>
              №{info.communityId}
            </span>
          </dd>
        </div>
        <div className={styles.row}>
          <dt>Centrality</dt>
          <dd>{info.centrality.toFixed(4)}</dd>
        </div>
        <div className={styles.row}>
          <dt>Входящий вес</dt>
          <dd>{info.inWeightSum.toFixed(2)}</dd>
        </div>
        <div className={styles.row}>
          <dt>Входящие рёбра</dt>
          <dd>{info.incomingCount}</dd>
        </div>
        <div className={styles.row}>
          <dt>Исходящие рёбра</dt>
          <dd>{info.outgoingCount}</dd>
        </div>
      </dl>
    </aside>
  );
}
