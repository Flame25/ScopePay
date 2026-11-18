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
    <main className="ui-shell">
      <div className="ui-topbar">
        <Link to="/projects" className="ui-brand"><span className="ui-mark">S</span>ScopePay</Link>
        <Link to="/projects" className="btn btn-ghost btn-sm">Back to projects</Link>
      </div>

      <header className="page-hero">
        <div>
          <p className="eyebrow">New proposal</p>
          <h1 className="page-title">New project</h1>
          <p className="page-sub">Paste what the client asked for. Next step is breaking it into milestones.</p>
        </div>
      </header>

      <form onSubmit={submit} className="block">
        <div className="form-stack">
          <div className="field-grid">
            <label className="field">
              Client name
              <input
                value={clientName}
                onChange={(e) => setClientName(e.target.value)}
                required
                maxLength={120}
                placeholder="Andrew"
              />
            </label>
            <label className="field">
              Client email
              <input
                value={clientEmail}
                onChange={(e) => setClientEmail(e.target.value)}
                type="email"
                required
                placeholder="billing@scopepay.com"
              />
            </label>
          </div>

          <label className="field">
            Client request
            <textarea
              value={rawRequest}
              onChange={(e) => setRawRequest(e.target.value)}
              required
              minLength={10}
              rows={6}
              placeholder="Paste the client's message or brief here…"
            />
            <span className="field-hint">Keep it verbatim — the original wording is part of the proposal record.</span>
          </label>

          {error && <p className="alert">{error}</p>}

          <div className="form-actions">
            <Link to="/projects" className="btn btn-ghost">Cancel</Link>
            <button className="btn btn-primary" disabled={submitting}>
              {submitting ? "Creating…" : "Create project"}
            </button>
          </div>
        </div>
      </form>
    </main>
  );
}
