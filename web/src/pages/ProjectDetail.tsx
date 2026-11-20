import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import type { ApiError } from "../api/client";
import { getProject, sendProject, updateProject, createMilestone, deleteMilestone, type Project, unlockProject } from "../api/projects";
import { formatMoney } from "../lib/money";
import { statusText, statusTone } from "../lib/status";

export default function ProjectDetail() {
  const { id } = useParams();
  const projectId = Number(id);
  const [project, setProject] = useState<Project | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [form, setForm] = useState({ title: "", clientName: "", clientEmail: "", scopeSummary: "" });
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);
  const [msForm, setMsForm] = useState({ title: "", amount: "", deliverables: "", acceptanceCriteria: "" });
  const [msError, setMsError] = useState<string | null>(null);
  const [msSaving, setMsSaving] = useState(false);

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
      await sendProject(projectId);
      setProject(await getProject(projectId));
    } catch (e) {
      setError((e as ApiError).error?.message ?? "failed to send");
    }
  }

  async function handleSave() {
    setFormError(null);
    setSaved(false);
    setSaving(true);
    try {
      setProject(await updateProject(projectId, form));
      setSaved(true);
      window.setTimeout(() => setSaved(false), 2500);
    } catch (e) {
      const apiErr = (e as ApiError).error;
      setFormError(Object.values(apiErr?.fields ?? {})[0] ?? apiErr?.message ?? "failed to save");
    } finally {
      setSaving(false);
    }
  }

  function editForm(patch: Partial<typeof form>) {
    setSaved(false);
    setForm((f) => ({ ...f, ...patch }));
  }

  async function handleAddMilestone() {
    setMsError(null);
    setMsSaving(true);
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
    } finally {
      setMsSaving(false);
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

  async function handleUnlock() {
    try {
      setProject(await unlockProject(projectId));
    } catch (e) {
      setError((e as ApiError).error?.message ?? "failed to unlock");
    }
  }

  if (error && !project) {
    return (
      <main className="ui-shell">
        <div className="ui-topbar">
          <Link to="/projects" className="ui-brand"><span className="ui-mark">S</span>ScopePay</Link>
        </div>
        <header className="page-hero">
          <div>
            <p className="eyebrow">Project</p>
            <h1 className="page-title">Project unavailable</h1>
            <p className="page-sub">{error}</p>
          </div>
        </header>
        <Link to="/projects" className="btn">Back to projects</Link>
      </main>
    );
  }

  if (!project) {
    return (
      <main className="ui-shell">
        <div className="ui-topbar">
          <Link to="/projects" className="ui-brand"><span className="ui-mark">S</span>ScopePay</Link>
        </div>
        <header className="page-hero">
          <div className="skeleton" style={{ height: 88, width: "100%" }} />
        </header>
      </main>
    );
  }

  const milestones = project.milestones ?? [];
  const total = milestones.reduce((sum, m) => sum + m.amountCents, 0);
  const currency = milestones[0]?.currency ?? "USD";
  const canSend = project.status === "draft" && milestones.length > 0;
  const editable = project.status === "draft" || project.status === "sent";
  const shareLink = project.shareToken
    ? `${window.location.origin}/share/${project.shareToken}`
    : null;

  async function handleCopy() {
    if (!shareLink) return;
    await navigator.clipboard.writeText(shareLink);
    setCopied(true);
    window.setTimeout(() => setCopied(false), 1600);
  }

  return (
    <main className="ui-shell">
      <div className="ui-topbar">
        <Link to="/projects" className="ui-brand"><span className="ui-mark">S</span>ScopePay</Link>
        <div className="ui-nav">
          <Link to="/projects" className="btn btn-ghost btn-sm">Back to projects</Link>
          <span className={`pill pill-${statusTone(project.status)}`}>{statusText(project.status)}</span>
        </div>
      </div>

      <header className="page-hero">
        <div>
          <p className="eyebrow">{project.clientName}</p>
          <h1 className="page-title">{project.title || "Untitled project"}</h1>
          <p className="page-sub">{project.clientEmail}</p>
        </div>
        <div className="page-actions">
          {project.status === "draft" && (
            <>
              {!canSend && <span className="block-meta">Add a milestone to send</span>}
              <button
                className="btn btn-primary"
                onClick={handleSend}
                disabled={!canSend}
                title={canSend ? undefined : "Add at least one milestone first"}
              >
                Send to client
              </button>
            </>
          )}
        </div>
      </header>

      {error && <p className="alert">{error}</p>}

      {shareLink && project.status !== "draft" && (
        <div className="share-link">
          <span className="share-link-label">✓ Share link</span>
          <code>{shareLink}</code>
          <button className="btn btn-sm" onClick={handleCopy}>{copied ? "Copied" : "Copy"}</button>
          <a className="btn btn-sm" href={shareLink} target="_blank" rel="noreferrer">Open</a>
        </div>
      )}

      {project.status === "revision_requested" && (
        <div className="block revision-banner">
          <div>
            <strong>Client requested changes</strong>
            <p>{project.editRequestNote || "No note provided."}</p>
          </div>
          <button className="btn btn-primary" onClick={handleUnlock}>Unlock editing</button>
        </div>
      )}

      <div className="detail-grid">
        <div className="detail-stack">
          <section className="block">
            <div className="block-head"><h2>Details</h2></div>
            <div className="form-stack">
              <label className="field">
                Title
                <input
                  value={form.title}
                  onChange={(e) => editForm({ title: e.target.value })}
                  maxLength={200}
                  placeholder="Landing page redesign"
                />
              </label>
              <div className="field-grid">
                <label className="field">
                  Client name
                  <input
                    value={form.clientName}
                    onChange={(e) => editForm({ clientName: e.target.value })}
                    required
                  />
                </label>
                <label className="field">
                  Client email
                  <input
                    value={form.clientEmail}
                    onChange={(e) => editForm({ clientEmail: e.target.value })}
                    type="email"
                    required
                  />
                </label>
              </div>
              <label className="field">
                Scope summary
                <textarea
                  rows={4}
                  value={form.scopeSummary}
                  onChange={(e) => editForm({ scopeSummary: e.target.value })}
                  placeholder="One paragraph the client sees as the scope of work."
                />
              </label>
              {formError && <p className="alert">{formError}</p>}
              <div className="form-actions">
                {editable ? (
                  <>
                    <span className="save-note" role="status" aria-live="polite">{saved ? "✓ Saved" : ""}</span>
                    <button className="btn btn-primary" onClick={handleSave} disabled={saving}>
                      {saving ? "Saving…" : "Save changes"}
                    </button>
                  </>
                ) : (
                  <span className="block-meta">Locked after client approval</span>
                )}

              </div>
            </div>
          </section>

          <section className="block">
            <div className="block-head"><h2>Original request</h2></div>
            <blockquote className="quote">{project.rawRequest}</blockquote>
          </section>
        </div>

        <section className="block">
          <div className="block-head">
            <h2>Milestones</h2>
            <span className="block-meta">
              {milestones.length} {milestones.length === 1 ? "payment" : "payments"}
            </span>
          </div>

          {milestones.length === 0 ? (
            <p className="empty">No milestones yet. Add the first payment step below.</p>
          ) : (
            <>
              <ol className="milestone-list">
                {milestones.map((m, i) => (
                  <li className="milestone-item" key={m.id}>
                    <span className="milestone-index">{i + 1}</span>
                    <div className="milestone-body">
                      <h3>{m.title}</h3>
                      {m.deliverables && <p>{m.deliverables}</p>}
                      {m.acceptanceCriteria && (
                        <span className="milestone-meta">Accepted when: {m.acceptanceCriteria}</span>
                      )}
                    </div>
                    <div className="ms-side">
                      <span className="milestone-amount">{formatMoney(m.amountCents, m.currency)}</span>
                      {editable && m.status === "pending" && (
                        <button className="btn btn-danger btn-sm" onClick={() => handleDeleteMilestone(m.id)}>
                          Delete
                        </button>
                      )}
                    </div>
                  </li>
                ))}
              </ol>
              <div className="ms-total">
                <span>Project total</span>
                <strong>{formatMoney(total, currency)}</strong>
              </div>
            </>
          )}

          {editable ? <details className="ms-add">
            <summary>+ Add milestone</summary>
            <div className="form-stack">
              <label className="field">
                Title
                <input
                  value={msForm.title}
                  onChange={(e) => setMsForm((f) => ({ ...f, title: e.target.value }))}
                  placeholder="Design handoff"
                  required
                />
              </label>
              <label className="field">
                Amount ({currency})
                <input
                  value={msForm.amount}
                  onChange={(e) => setMsForm((f) => ({ ...f, amount: e.target.value }))}
                  type="number"
                  step="0.01"
                  min="0.01"
                  placeholder="500.00"
                  required
                />
              </label>
              <label className="field">
                Deliverables
                <textarea
                  rows={3}
                  value={msForm.deliverables}
                  onChange={(e) => setMsForm((f) => ({ ...f, deliverables: e.target.value }))}
                />
              </label>
              <label className="field">
                Acceptance criteria
                <textarea
                  rows={3}
                  value={msForm.acceptanceCriteria}
                  onChange={(e) => setMsForm((f) => ({ ...f, acceptanceCriteria: e.target.value }))}
                />
              </label>
              {msError && <p className="alert">{msError}</p>}
              <div className="form-actions">
                <button className="btn btn-primary" onClick={handleAddMilestone} disabled={msSaving}>
                  {msSaving ? "Adding…" : "Add Milestones"}
                </button>
              </div>
            </div>
          </details>
            : <p className="block-meta">Milestones are locked once the client approves.</p>}
        </section>
      </div >
    </main >
  );
}
