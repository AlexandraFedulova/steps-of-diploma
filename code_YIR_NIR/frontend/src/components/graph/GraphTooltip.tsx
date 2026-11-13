import styles from "./GraphTooltip.module.css";

interface GraphTooltipProps {
  x: number;
  y: number;
  text: string;
}

export function GraphTooltip({ x, y, text }: GraphTooltipProps) {
  return (
    <div
      className={styles.tooltip}
      style={{ left: x, top: y }}
      role="tooltip"
    >
      {text}
    </div>
  );
}
