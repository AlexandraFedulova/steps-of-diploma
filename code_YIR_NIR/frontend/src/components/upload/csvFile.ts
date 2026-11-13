const CSV_EXTENSION = /\.csv$/i;

const CSV_MIME_TYPES = new Set([
  "text/csv",
  "application/csv",
  "text/comma-separated-values",
  "application/vnd.ms-excel",
]);

export function isCsvFile(file: File): boolean {
  if (CSV_EXTENSION.test(file.name)) {
    return true;
  }
  return file.type !== "" && CSV_MIME_TYPES.has(file.type);
}

export function formatFileSize(bytes: number): string {
  if (bytes < 1024) {
    return `${bytes} Б`;
  }
  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} КБ`;
  }
  return `${(bytes / (1024 * 1024)).toFixed(1)} МБ`;
}
