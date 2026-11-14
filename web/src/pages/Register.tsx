import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../auth/AuthContext";
import { register } from "../api/auth";
import type { ApiError } from "../api/client";

export default function Register() {
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const navigate = useNavigate();

  async function submit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      await register(email, name, password);
      await login(email, password);
      navigate("/");
    } catch (err) {
      const apiErr = (err as ApiError).error;
      const fieldError = apiErr?.fields ? Object.values(apiErr.fields)[0] : undefined;
      setError(fieldError ?? apiErr?.message ?? "registration failed");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main className="page">
      <h1>Create account</h1>
      <p className="muted">Start structuring projects with ScopePay</p>

      <form onSubmit={submit} className="card form-card">
        <label className="field">
          Name
          <input value={name} onChange={(e) => setName(e.target.value)} required autoComplete="name" maxLength={120} />
        </label>

        <label className="field">
          Email
          <input value={email} onChange={(e) => setEmail(e.target.value)} type="email" required autoComplete="email" />
        </label>

        <label className="field">
          Password
          <input
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            type="password"
            required
            autoComplete="new-password"
            minLength={8}
            maxLength={72}
          />
        </label>

        {error && <p className="form-error">{error}</p>}

        <button type="submit" className="btn-primary" disabled={submitting}>
          {submitting ? "Creating…" : "Create account"}
        </button>

        <p className="form-foot">
          Already registered? <Link to="/login">Sign in</Link>
        </p>
      </form>
    </main>
  );
}
