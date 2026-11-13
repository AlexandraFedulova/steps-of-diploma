import { useId, useRef, useState, type DragEvent } from "react";
import { formatFileSize, isCsvFile } from "./csvFile";
import styles from "./CsvFileUpload.module.css";

export interface CsvFileUploadProps {
  file: File | null;
  onFileChange: (file: File | null) => void;
  disabled?: boolean;
}

export function CsvFileUpload({ file, onFileChange, disabled = false }: CsvFileUploadProps) {
  const inputId = useId();
  const inputRef = useRef<HTMLInputElement>(null);
  const [dragActive, setDragActive] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const acceptFile = (candidate: File | undefined) => {
    if (disabled) {
      return;
    }

    if (!candidate) {
      return;
    }

    if (!isCsvFile(candidate)) {
      setError("Допустимы только файлы с расширением .csv");
      return;
    }

    setError(null);
    onFileChange(candidate);
  };

  const onInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    acceptFile(event.target.files?.[0]);
    event.target.value = "";
  };

  const onDragEnter = (event: DragEvent) => {
    event.preventDefault();
    event.stopPropagation();
    setDragActive(true);
  };

  const onDragOver = (event: DragEvent) => {
    event.preventDefault();
    event.stopPropagation();
    setDragActive(true);
  };

  const onDragLeave = (event: DragEvent) => {
    event.preventDefault();
    event.stopPropagation();
    setDragActive(false);
  };

  const onDrop = (event: DragEvent) => {
    event.preventDefault();
    event.stopPropagation();
    setDragActive(false);
    acceptFile(event.dataTransfer.files[0]);
  };

  const openFileDialog = () => {
    inputRef.current?.click();
  };

  const clearFile = () => {
    setError(null);
    onFileChange(null);
    if (inputRef.current) {
      inputRef.current.value = "";
    }
  };

  const dropzoneClass = [
    styles.dropzone,
    dragActive ? styles.dropzoneActive : "",
    error ? styles.dropzoneError : "",
    disabled ? styles.dropzoneDisabled : "",
  ]
    .filter(Boolean)
    .join(" ");

  return (
    <div className={styles.root}>
      <div
        className={dropzoneClass}
        onDragEnter={disabled ? undefined : onDragEnter}
        onDragOver={disabled ? undefined : onDragOver}
        onDragLeave={disabled ? undefined : onDragLeave}
        onDrop={disabled ? undefined : onDrop}
        onClick={disabled ? undefined : openFileDialog}
        aria-label="Зона перетаскивания CSV"
        aria-disabled={disabled}
      >
        <UploadIcon className={styles.icon} />
        <p className={styles.title}>Перетащите CSV сюда</p>
        <p className={styles.hint}>или выберите файл на компьютере</p>
        <button
          type="button"
          className={styles.button}
          disabled={disabled}
          onClick={(event) => {
            event.stopPropagation();
            if (!disabled) {
              openFileDialog();
            }
          }}
        >
          Выбрать файл
        </button>
      </div>

      <input
        ref={inputRef}
        id={inputId}
        className={styles.input}
        type="file"
        accept=".csv,text/csv"
        onChange={onInputChange}
        disabled={disabled}
        tabIndex={-1}
      />

      {error && (
        <p className={styles.error} role="alert">
          {error}
        </p>
      )}

      {file && (
        <div className={styles.fileInfo}>
          <div>
            <p className={styles.fileName}>{file.name}</p>
            <p className={styles.fileMeta}>{formatFileSize(file.size)}</p>
          </div>
          <button
            type="button"
            className={styles.clearButton}
            onClick={clearFile}
            disabled={disabled}
          >
            Удалить
          </button>
        </div>
      )}
    </div>
  );
}

function UploadIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      width="40"
      height="40"
      viewBox="0 0 24 24"
      fill="none"
      aria-hidden="true"
    >
      <path
        d="M12 16V4m0 0l-4 4m4-4 4 4M4 17v2a2 2 0 002 2h12a2 2 0 002-2v-2"
        stroke="currentColor"
        strokeWidth="1.75"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}
