import type { GraphPayload, UploadGraphResponse } from "@/types";
import { validateGraphPayload } from "@/utils/validateGraphPayload";
import { apiErrorFromBody, ApiError, readJsonBody } from "./api";

const API_BASE = "/api";

export async function uploadGraph(file: File): Promise<GraphPayload> {
  const body = new FormData();
  body.append("file", file, file.name);

  let response: Response;

  try {
    response = await fetch(`${API_BASE}/upload`, {
      method: "POST",
      body,
    });
  } catch {
    throw new ApiError(
      "Не удалось связаться с сервером. Убедитесь, что backend запущен (порт 8080).",
      0,
    );
  }

  const data = await readJsonBody<UploadGraphResponse>(response);

  if (!response.ok || data.ok === false) {
    throw apiErrorFromBody(data, response, `Ошибка загрузки (${response.status})`);
  }

  if (!data.graph) {
    throw new ApiError("Сервер не вернул данные графа", response.status);
  }

  return validateGraphPayload(data.graph);
}
