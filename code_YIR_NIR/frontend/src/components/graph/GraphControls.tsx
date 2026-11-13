import type { Core } from "cytoscape";
import type { RefObject } from "react";
import styles from "./GraphControls.module.css";

interface GraphControlsProps {
  cyRef: RefObject<Core | null>;
}

export function GraphControls({ cyRef }: GraphControlsProps) {
  const withCy = (action: (cy: Core) => void) => {
    const cy = cyRef.current;
    if (cy) {
      action(cy);
    }
  };

  return (
    <div className={styles.controls} role="toolbar" aria-label="Управление графом">
      <button
        type="button"
        className={styles.button}
        title="Приблизить"
        onClick={() => withCy((cy) => cy.zoom({ level: cy.zoom() * 1.2, renderedPosition: { x: cy.width() / 2, y: cy.height() / 2 } }))}
      >
        +
      </button>
      <button
        type="button"
        className={styles.button}
        title="Отдалить"
        onClick={() => withCy((cy) => cy.zoom({ level: cy.zoom() / 1.2, renderedPosition: { x: cy.width() / 2, y: cy.height() / 2 } }))}
      >
        −
      </button>
      <button
        type="button"
        className={styles.button}
        title="Вписать в область"
        onClick={() => withCy((cy) => cy.fit(undefined, 48))}
      >
        ⊡
      </button>
      <span className={styles.hint}>Колёсико — масштаб · перетаскивание фона — панорама · узел — перетащить</span>
    </div>
  );
}
