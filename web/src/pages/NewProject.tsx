import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { createProject } from "../api/projects";
import type { ApiError } from "../api/client";

export default function NewProject() {
  const [clientName, setClientName] = useState("");
  const [clientEmail, setClientEmail] = useState("");
  const [rawRequest, setRawRequest] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const navigate = useNavigate();

  async function submit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      const p = await createProject({ clientName, clientEmail, rawRequest });
      navigate(`/projects/${p.id}`);
    } catch (err) {
      const apiErr = (err as ApiError).error;
      setError(Object.values(apiErr?.fields ?? {})[0] ?? apiErr?.message ?? "failed to create");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main className="page">
      <h1>New project</h1>
      <p className="muted">Describe what the client asked for</p>
      <form onSubmit={submit} className="card form-card">
        <label className="field">
          Client name
          <input value={clientName} onChange={(e) => setClientName(e.target.value)} required maxLength={120} />
        </label>
        <label className="field">
          Client email
          <input value={clientEmail} onChange={(e) => setClientEmail(e.target.value)} type="email" required />
        </label>
        <label className="field">
          Client request
          <textarea value={rawRequest} onChange={(e) => setRawRequest(e.target.value)} required minLength={10} rows={6} />
        </label>
        {error && <p className="form-error">{error}</p>}
        <button className="btn-primary" disabled={submitting}>{submitting ? "Creating…" : "Create project"}</button>
        <p className="form-foot"><Link to="/projects">Back to projects</Link></p>
      </form>
    </main>
  );
}
