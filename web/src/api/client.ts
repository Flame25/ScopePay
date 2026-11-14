const BASE_URL = import.meta.env.VITE_API_URL ?? "/api";

export type ApiError = { error: { code: string; message: string; fields?: Record<string, string> } };

export type HealthResponse = {
  status: string;
  db: string;
};

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, init);
  const body = await res.json().catch(() => null);
  if (!res.ok) {
    throw (body as ApiError) ?? { error: { code: "unknown", message: `HTTP ${res.status}` } };
  }
  return body as T
}

export async function getHealth(): Promise<HealthResponse> {
  return request<HealthResponse>("/health")
}
