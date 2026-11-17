import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import type { ApiError } from "../api/client";
import { getProject, sendProject, updateProject, createMilestone, deleteMilestone, type Project } from "../api/projects";
import { formatMoney } from "../lib/money";

export default function ProjectDetail() {
  const { id } = useParams();
  const projectId = Number(id);
  const [project, setProject] = useState<Project | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [shareUrl, setShareUrl] = useState<string | null>(null);
  const [form, setForm] = useState({ title: "", clientName: "", clientEmail: "", scopeSummary: "" });
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);
  const [msForm, setMsForm] = useState({ title: "", amount: "", deliverables: "", acceptanceCriteria: "" });
  const [msError, setMsError] = useState<string | null>(null);

  useEffect(() => {
    getProject(projectId)
      .then(setProject)
      .catch((e) => setError((e as ApiError).error?.message ?? "failed to load"));
  }, [projectId]);

  useEffect(() => {
    if (!project) return;
    setForm({
      title: project.title,
      clientName: project.clientName,
      clientEmail: project.clientEmail,
      scopeSummary: project.scopeSummary,
    });
  }, [project]);

  async function handleSend() {
    setError(null);
    try {
      const res = await sendProject(projectId);
      setShareUrl(res.shareUrl);
      setProject(await getProject(projectId));
    } catch (e) {
      setError((e as ApiError).error?.message ?? "failed to send");
    }
  }

  async function handleSave() {
    setFormError(null);
    setSaving(true);
    try {
      setProject(await updateProject(projectId, form));
    } catch (e) {
      const apiErr = (e as ApiError).error;
      setFormError(Object.values(apiErr?.fields ?? {})[0] ?? apiErr?.message ?? "failed to save");
    } finally {
      setSaving(false);
    }
  }

  async function handleAddMilestone() {
    setMsError(null);
    try {
      await createMilestone(projectId, {
        title: msForm.title,
        deliverables: msForm.deliverables,
        acceptanceCriteria: msForm.acceptanceCriteria,
        amountCents: Math.round(parseFloat(msForm.amount) * 100),   // dollars -> cents
      });
      setMsForm({ title: "", amount: "", deliverables: "", acceptanceCriteria: "" });
      setProject(await getProject(projectId));
    } catch (e) {
      const apiErr = (e as ApiError).error;
      setMsError(Object.values(apiErr?.fields ?? {})[0] ?? apiErr?.message ?? "failed to add milestone");
    }
  }

  async function handleDeleteMilestone(mid: number) {
    try {
      await deleteMilestone(projectId, mid);
      setProject(await getProject(projectId));
    } catch (e) {
      setMsError((e as ApiError).error?.message ?? "failed to delete");
    }
  }

  if (error && !project) return <main className="page"><p className="form-error">{error}</p></main>;
  if (!project) return <main className="page"><p className="muted">Loading…</p></main>;

  return (
    <main className="page">
      <Link to="/projects" className="form-foot">← Back to projects</Link>
      <header className="page-head">
        <h1>{project.title || project.clientName}</h1>
        <div className="head-actions">
          <span className="status-chip">{project.status}</span>
          {project.status === "draft" && (
            <button className="btn-primary" onClick={handleSend}>Send to client</button>
          )}
        </div>
      </header>

      {shareUrl && (
        <div className="card ok">
          <strong>Sent</strong>
          <p className="form-foot">Share link: <code>{shareUrl}</code></p>
        </div>
      )}
      {error && <p className="form-error">{error}</p>}
      <section className="card">
        <h2>Details</h2>

        <div className="form-card">
          <label className="field">
            Title
            <input value={form.title} onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))} maxLength={200} />
          </label>
          <label className="field">
            Client name
            <input value={form.clientName} onChange={(e) => setForm((f) => ({ ...f, clientName: e.target.value }))} required />
          </label>
          <label className="field">
            Client email
            <input value={form.clientEmail} onChange={(e) => setForm((f) => ({ ...f, clientEmail: e.target.value }))} type="email" required />
          </label>
          <label className="field">
            Scope summary
            <textarea rows={4} value={form.scopeSummary} onChange={(e) => setForm((f) => ({ ...f, scopeSummary: e.target.value }))} />
          </label>
          {formError && <p className="form-error">{formError}</p>}
          <button className="btn-primary" onClick={handleSave} disabled={saving}>
            {saving ? "Saving…" : "Save changes"}
          </button>
        </div>
      </section>

      <section className="card">
        <h2>Original request</h2>
        <p className="muted">{project.rawRequest}</p>
      </section>

      <section className="card">
        <h2>Milestones</h2>
        {project.milestones?.length === 0 && <p className="muted">No milestones yet.</p>}
        <div className="card-container">
          {project.milestones?.map((m) => (
            <div className="ms-row" key={m.id}>
              <div>
                <strong>{m.order + 1}. {m.title}</strong>
                <p className="muted">{formatMoney(m.amountCents, m.currency)} · {m.status}</p>
              </div>
              {m.status === "pending" && (
                <button className="btn" onClick={() => handleDeleteMilestone(m.id)}>Delete</button>
              )}
            </div>
          ))}
        </div>

        <hr className="sep" />
        <h3>Add milestone</h3>
        <div className="form-card">
          <label className="field">
            Title
            <input value={msForm.title} onChange={(e) => setMsForm((f) => ({ ...f, title: e.target.value }))} required />
          </label>
          <label className="field">
            Amount ({project.milestones?.[0]?.currency ?? "USD"})
            <input value={msForm.amount} onChange={(e) => setMsForm((f) => ({ ...f, amount: e.target.value }))} type="number" step="0.01" min="0.01" required />
          </label>
          <label className="field">
            Deliverables
            <textarea rows={3} value={msForm.deliverables} onChange={(e) => setMsForm((f) => ({ ...f, deliverables: e.target.value }))} />
          </label>
          <label className="field">
            Acceptance criteria
            <textarea rows={3} value={msForm.acceptanceCriteria} onChange={(e) => setMsForm((f) => ({ ...f, acceptanceCriteria: e.target.value }))} />
          </label>
          {msError && <p className="form-error">{msError}</p>}
          <button className="btn-primary" onClick={handleAddMilestone}>Add milestone</button>
        </div>
      </section>
    </main>
  );
}
