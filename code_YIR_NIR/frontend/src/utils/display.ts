/** Отображение булева значения в интерфейсе */
export function formatYesNo(value: boolean): string {
  return value ? "да" : "нет";
}

/** Статус health-check с backend */
export function formatHealthStatus(status: string): string {
  if (status === "ok") {
    return "работает";
  }
  return status;
}
