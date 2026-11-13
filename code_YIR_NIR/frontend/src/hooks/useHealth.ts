import { useEffect, useState } from "react";
import type { HealthResponse } from "@/types";
import { toUserMessage } from "@/services/api";
import { fetchHealth } from "@/services/healthService";

interface UseHealthResult {
  health: HealthResponse | null;
  loading: boolean;
  error: string | null;
}

export function useHealth(): UseHealthResult {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    const load = async () => {
      try {
        const data = await fetchHealth();
        if (!cancelled) {
          setHealth(data);
          setError(null);
        }
      } catch (e) {
        if (!cancelled) {
          setError(toUserMessage(e));
          setHealth(null);
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    void load();

    return () => {
      cancelled = true;
    };
  }, []);

  return { health, loading, error };
}
