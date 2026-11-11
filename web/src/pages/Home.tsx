import { useCallback, useEffect, useState } from "react";
import { getHealth, type HealthResponse } from "../api/client";

export default function Home() {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      setHealth(await getHealth());
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  return (
    <main className="page">
      <h1>ScopePay</h1>
      <p className="muted">AI-powered freelance project payments</p>

      {loading && <p>Checking backend…</p>}

      {!loading && error && (
        <div className="card error">
          <p>{error}</p>
          <button onClick={() => void load()}>Retry</button>
        </div>
      )}

      {!loading && !error && health && (
        <div className="card ok">
          <strong>Backend healthy</strong>
          <p>status: {health.status}</p>
          <p>db: {health.db}</p>
        </div>
      )}
    </main>
  );
}
