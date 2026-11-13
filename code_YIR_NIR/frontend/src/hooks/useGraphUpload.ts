import { useCallback, useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useGraph } from "@/context/GraphContext";
import { toUserMessage } from "@/services/api";
import { uploadGraph } from "@/services/uploadService";

interface UseGraphUploadOptions {
  /** После успешного анализа открыть страницу графа */
  openGraphOnSuccess?: boolean;
}

interface UseGraphUploadResult {
  file: File | null;
  loading: boolean;
  error: string | null;
  setFile: (file: File | null) => void;
  reset: () => void;
}

export function useGraphUpload(options: UseGraphUploadOptions = {}): UseGraphUploadResult {
  const { openGraphOnSuccess = true } = options;
  const navigate = useNavigate();
  const { setGraph } = useGraph();
  const [file, setFileState] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const requestId = useRef(0);

  const reset = useCallback(() => {
    requestId.current += 1;
    setFileState(null);
    setGraph(null);
    setLoading(false);
    setError(null);
  }, [setGraph]);

  const setFile = useCallback(
    (next: File | null) => {
      setFileState(next);
      if (!next) {
        setGraph(null);
        setError(null);
        setLoading(false);
      }
    },
    [setGraph],
  );

  useEffect(() => {
    if (!file) {
      return;
    }

    const currentRequest = ++requestId.current;
    setLoading(true);
    setError(null);
    setGraph(null);

    uploadGraph(file)
      .then((payload) => {
        if (requestId.current !== currentRequest) {
          return;
        }
        setGraph(payload);
        if (openGraphOnSuccess) {
          navigate("/graph", { replace: false });
        }
      })
      .catch((err: unknown) => {
        if (requestId.current !== currentRequest) {
          return;
        }
        setError(toUserMessage(err));
        setGraph(null);
      })
      .finally(() => {
        if (requestId.current === currentRequest) {
          setLoading(false);
        }
      });
  }, [file, setGraph, navigate, openGraphOnSuccess]);

  return { file, loading, error, setFile, reset };
}
