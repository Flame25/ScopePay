import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../auth/AuthContext";
import type { ApiError } from "../api/client";

export default function Login() {
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const [submitting, setSubmitting] = useState(false);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      await login(email, password);
      navigate("/");
    } catch (err) {
      setError((err as ApiError).error?.message ?? "login failed");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main className="page">
      <h1>Sign in</h1>
      <p className="muted">Welcome back to ScopePay</p>

      <form onSubmit={submit} className="card form-card">
        <label className="field">
          Email
          <input value={email} onChange={(e) => setEmail(e.target.value)} type="email" required autoComplete="email" />
        </label>

        <label className="field">
          Password
          <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" required autoComplete="current-password" />
        </label>

        {error && <p className="form-error">{error}</p>}

        <button type="submit" className="btn-primary" disabled={submitting}>
          {submitting ? "Signing in…" : "Sign in"}
        </button>

        <p className="form-foot">
          No account? <Link to="/register">Create one</Link>
        </p>
      </form>
    </main>
  );

}
