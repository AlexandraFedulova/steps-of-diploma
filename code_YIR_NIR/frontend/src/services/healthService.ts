import type { HealthResponse } from "@/types";
import { apiFetch } from "./api";

export function fetchHealth(): Promise<HealthResponse> {
  return apiFetch<HealthResponse>("/health");
}
