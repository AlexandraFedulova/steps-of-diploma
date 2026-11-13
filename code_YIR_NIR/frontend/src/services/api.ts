const API_BASE = "/api";

export class ApiError extends Error {
  constructor(
    message: string,
    readonly status: number,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

interface ApiErrorBody {
  ok?: boolean;
  error?: string;
}

export async function readJsonBody<T>(response: Response): Promise<T> {
  const text = await response.text();

  if (!text.trim()) {
    throw new ApiError(
      response.ok ? "Пустой ответ сервера" : `Ошибка сервера (${response.status})`,
      response.status,
    );
  }

  try {
    return JSON.parse(text) as T;
  } catch {
    throw new ApiError(
      response.ok ? "Сервер вернул не JSON" : `Ошибка сервера (${response.status})`,
      response.status,
    );
  }
}

export function apiErrorFromBody(
  data: ApiErrorBody,
  response: Response,
  fallback: string,
): ApiError {
  return new ApiError(data.error ?? fallback, response.status);
}

export function toUserMessage(err: unknown): string {
  if (err instanceof ApiError) {
    if (err.status === 0) {
      return err.message;
    }
    return err.message;
  }

  if (err instanceof TypeError) {
    return "Не удалось связаться с сервером. Убедитесь, что backend запущен (порт 8080).";
  }

  if (err instanceof Error) {
    return err.message;
  }

  return "Неизвестная ошибка";
}

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  let response: Response;

  try {
    response = await fetch(`${API_BASE}${path}`, init);
  } catch {
    throw new ApiError(
      "Не удалось связаться с сервером. Убедитесь, что backend запущен (порт 8080).",
      0,
    );
  }

  if (!response.ok) {
    const data = await readJsonBody<ApiErrorBody>(response).catch(() => ({} as ApiErrorBody));
    throw apiErrorFromBody(data, response, `Ошибка запроса: ${response.status}`);
  }

  return readJsonBody<T>(response);
}
