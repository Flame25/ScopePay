const BASE_URL = import.meta.env.VITE_API_URL ?? "/api";

export type HealthResponse = {
  status: string;
  db: string;
};

export async function getHealth(): Promise<HealthResponse> {
  const res = await fetch(`${BASE_URL}/health`);
  if (!res.ok) {
    throw new Error(`health check failed: ${res.status}`);
  }
  return (await res.json()) as HealthResponse;
}
