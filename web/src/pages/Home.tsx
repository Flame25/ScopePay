import { useCallback, useEffect, useState, type FormEvent } from "react";
import { getHealth, type ApiError, type HealthResponse } from "../api/client";
import { useAuth } from "../auth/AuthContext";
import Login from "./Login";
import { useNavigate } from "react-router-dom";
import { logout } from "../api/auth";

export default function Home() {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const { user, logout, loading } = useAuth();
  const [submitting, setSubmitting] = useState(false);
  const navigate = useNavigate();

  async function handleLogout(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      logout();
      navigate("/");
    } catch (err) {
      setError((err as ApiError).error?.message ?? "logout failed");
    } finally {
      setSubmitting(false);
    }
  }

  const load = useCallback(async () => {
    setError(null);
    try {
      setHealth(await getHealth());
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
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

      <div className="card-container">
        {!loading && error && (
          <div className="card error">
            <p><b> Please retry & make backend running ... </b></p>
            <button id="retryBtn" onClick={() => void load()}> <b>Retry </b></button>
          </div>
        )}

        {!loading && !error && health && (
          <div className="card ok">
            <strong>Backend healthy</strong>
            <p>status: {health.status}</p>
            <p>db: {health.db}</p>
          </div>
        )}

        <div className="card">
          {loading ? <p>Checking session…</p> : <p>{user ? `Logged in as ${user.email}` : "Logged out"}</p>}
          {!loading && !user && (
            <button id="login-btn" onClick={() => navigate("/login")}>
              Login
            </button>
          )}
          {!loading && user && (
            <button className="logout-btn" onClick={handleLogout} disabled={submitting} >
              {submitting ? "Logging out…" : "Logout"}
            </button>
          )}
        </div>
      </div>
    </main >
  );
}
