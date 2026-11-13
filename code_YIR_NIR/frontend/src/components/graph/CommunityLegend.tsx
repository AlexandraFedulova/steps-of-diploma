import type { CommunityLegendItem } from "@/utils/communityVisual";
import styles from "./CommunityLegend.module.css";

interface CommunityLegendProps {
  items: CommunityLegendItem[];
}

export function CommunityLegend({ items }: CommunityLegendProps) {
  if (items.length === 0) {
    return null;
  }

  return (
    <div className={styles.root} role="list" aria-label="Сообщества">
      <span className={styles.title}>Сообщества</span>
      <ul className={styles.list}>
        {items.map((item) => (
          <li key={item.communityId} className={styles.item} role="listitem">
            <span
              className={styles.swatch}
              style={{ backgroundColor: item.color }}
              aria-hidden="true"
            />
            <span className={styles.label}>
              №{item.communityId}
              <span className={styles.count}> ({item.nodeCount})</span>
            </span>
          </li>
        ))}
      </ul>
    </div>
  );
}
