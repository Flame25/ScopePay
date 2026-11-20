import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import type { ApiError } from "../api/client";
import { approveShare, getShare, requestChanges, type ShareProject } from "../api/share";
import { formatMoney } from "../lib/money";
import { statusText, statusTone } from "../lib/status";
import "../styles/share.css";

export default function Share() {
  const { token } = useParams();
  const [data, setData] = useState<ShareProject | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const [note, setNote] = useState("");
  const [showRequest, setShowRequest] = useState(false);

  useEffect(() => {
    if (!token) return;
    getShare(token)
      .then(setData)
      .catch((e) => setError((e as ApiError).error?.message ?? "proposal not found"));
  }, [token]);

  async function handleApprove() {
    if (!token) return;
    setBusy(true);
    setError(null);
    try {
      setData(await approveShare(token));
    } catch (e) {
      setError((e as ApiError).error?.message ?? "failed to approve");
    } finally {
      setBusy(false);
    }
  }

  async function handleRequestChanges() {
    if (!token) return;
    setBusy(true);
    setError(null);
    try {
      setData(await requestChanges(token, note));
      setShowRequest(false);
      setNote("");
    } catch (e) {
      setError((e as ApiError).error?.message ?? "failed to request changes");
    } finally {
      setBusy(false);
    }
  }

  if (error && !data) {
    return (
      <main className="ui-shell share-page">
        <div className="ui-topbar">
          <span className="ui-brand"><span className="ui-mark">S</span>ScopePay</span>
        </div>
        <div className="share-hero">
          <h1 className="share-title">Proposal unavailable</h1>
          <p className="share-sub">{error}</p>
        </div>
      </main>
    );
  }

  if (!data) {
    return (
      <main className="ui-shell share-page">
        <div className="ui-topbar">
          <span className="ui-brand"><span className="ui-mark">S</span>ScopePay</span>
        </div>
        <div className="share-hero"><div className="skeleton share-loading" /></div>
      </main>
    );
  }

  const total = data.milestones.reduce((sum, m) => sum + m.amountCents, 0);
  const currency = data.milestones[0]?.currency ?? "USD";

  return (
    <main className="ui-shell share-page">
      <div className="ui-topbar">
        <span className="ui-brand"><span className="ui-mark">S</span>ScopePay</span>
        <div className="ui-nav">
          <span className="block-meta">{data.viewerRole === "client" ? "Client view" : "Freelancer view"}</span>
          <span className={`pill pill-${statusTone(data.status)}`}>{statusText(data.status)}</span>
        </div>
      </div>

      <header className="share-hero">
        <p className="eyebrow">Project proposal</p>
        <h1 className="share-title">{data.title || "Untitled project"}</h1>
        <p className="share-sub">Prepared for {data.clientName}</p>
      </header>

      {error && <p className="alert">{error}</p>}

      {data.scopeSummary && (
        <section className="block">
          <div className="block-head"><h2>Scope</h2></div>
          <p className="share-scope">{data.scopeSummary}</p>
        </section>
      )}

      <section className="block share-block">
        <div className="block-head">
          <h2>Milestones</h2>
          <span className="block-meta">
            {data.milestones.length} {data.milestones.length === 1 ? "payment" : "payments"}
          </span>
        </div>

        {data.milestones.length === 0 ? (
          <p className="empty">No milestones defined yet.</p>
        ) : (
          <>
            <ol className="milestone-list">
              {data.milestones.map((m, i) => (
                <li className="milestone-item" key={m.id}>
                  <span className="milestone-index">{i + 1}</span>
                  <div className="milestone-body">
                    <h3>{m.title}</h3>
                    {m.deliverables && <p>{m.deliverables}</p>}
                    {m.acceptanceCriteria && (
                      <span className="milestone-meta">Accepted when: {m.acceptanceCriteria}</span>
                    )}
                  </div>
                  <span className="milestone-amount">{formatMoney(m.amountCents, m.currency)}</span>
                </li>
              ))}
            </ol>

            <div className="share-total">
              <span>Project total</span>
              <strong>{formatMoney(total, currency)}</strong>
            </div>
          </>
        )}
      </section>

      <div className="share-actions">
        {data.status === "sent" && data.viewerRole === "client" && (
          <>
            <div className="share-actions-copy">
              <strong>Ready to proceed?</strong>
              <p>Approving confirms the scope and milestone amounts above.</p>
            </div>

            <button className="share-cta" onClick={handleApprove} disabled={busy}>
              {busy ? "Approving…" : "Approve proposal"}
            </button>

            {!showRequest && (
              <button className="btn btn-sm" onClick={() => setShowRequest(true)}>Request changes instead</button>
            )}
          </>
        )}

        {data.status === "sent" && data.viewerRole === "freelancer" && (
          <>
            <div className="share-actions-copy">
              <strong>Sent to {data.clientName}</strong>
              <p>Waiting for the client to review and approve. You can still request changes.</p>
            </div>
            <button className="btn btn-sm" onClick={() => setShowRequest(true)}>Request changes</button>
          </>
        )}

        {showRequest && (
          <div className="share-request">
            <label className="field">
              What should change?
              <textarea
                rows={3}
                value={note}
                onChange={(e) => setNote(e.target.value)}
                placeholder="e.g. split milestone 2, and add a final QA pass"
                maxLength={2000}
              />
            </label>
            <div className="form-actions">
              <button className="btn btn-ghost btn-sm" onClick={() => setShowRequest(false)}>Cancel</button>
              <button className="btn btn-sm" onClick={handleRequestChanges} disabled={busy}>
                {busy ? "Sending…" : "Send request"}
              </button>
            </div>
          </div>
        )}

        {data.status === "approved" && (
          <>
            <div className="share-approved">
              <span className="share-check">✓</span>
              Proposal approved — milestones can now be invoiced.
            </div>

            {!showRequest && (
              <button className="btn btn-sm" onClick={() => setShowRequest(true)}>Request changes</button>
            )}
          </>
        )}

        {data.status === "revision_requested" && (
          <div className="share-actions-copy">
            <strong>Changes requested</strong>
            <p>Waiting for the freelancer to reopen the proposal and revise it.</p>
            {data.editRequestNote && <p className="block-meta">Your note: {data.editRequestNote}</p>}
          </div>
        )}

        {data.status === "draft" && (
          <div className="share-actions-copy">
            <strong>Not sent yet</strong>
            <p>The freelancer still needs to send this proposal.</p>
          </div>
        )}

        {data.status === "in_progress" && (
          <div className="share-actions-copy">
            <strong>Work in progress</strong>
            <p>Milestone payments are collected as work is completed.</p>
          </div>
        )}

        {data.status === "complete" && (
          <div className="share-approved">
            <span className="share-check">✓</span>
            Project complete.
          </div>
        )}
      </div>

      <p className="share-footnote">
        Payments are made directly to the freelancer's PayPal account.
        ScopePay does not hold or escrow funds.
      </p>
    </main>
  );
}
