import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { type ApiError } from "../api/client";
import { useAuth } from "../auth/AuthContext";
import { listProjects, type Project } from "../api/projects";
import { formatMoney } from "../lib/money";
import { statusText, statusTone } from "../lib/status";

export default function Projects() {
  const [items, setItems] = useState<Project[] | null>(null);
  const [error, setError] = useState<string | null>(null);
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    listProjects()
      .then(setItems)
      .catch((e) => setError((e as ApiError).error?.message ?? "failed to load"));
  }, []);

  function totalCents(p: Project) {
    return (p.milestones ?? []).reduce((sum, m) => sum + m.amountCents, 0);
  }

  const count = items?.length ?? 0;

  return (
    <main className="ui-shell">
      <div className="ui-topbar">
        <Link to="/projects" className="ui-brand"><span className="ui-mark">S</span>ScopePay</Link>
        <div className="ui-nav">
          <span>{user?.email}</span>
          <button className="btn btn-ghost btn-sm" onClick={() => { logout(); navigate("/login"); }}>
            Log out
          </button>
        </div>
      </div>

      <header className="page-hero">
        <div>
          <p className="eyebrow">Workspace</p>
          <h1 className="page-title">Projects</h1>
          <p className="page-sub">
            {items ? `${count} ${count === 1 ? "proposal" : "proposals"} with milestone payments` : "Loading your proposals…"}
          </p>
        </div>
        <Link to="/projects/new" className="btn btn-primary">New project</Link>
      </header>

      {error && <p className="alert">{error}</p>}

      {!items && !error && (
        <div className="project-grid">
          {[0, 1, 2, 3].map((i) => <div className="skeleton skeleton-card" key={i} />)}
        </div>
      )}

      {items && items.length > 0 && (
        <div className="project-grid">
          {items.map((p) => {
            const milestones = p.milestones ?? [];
            return (
              <Link to={`/projects/${p.id}`} className="project-card" key={p.id}>
                <div className="project-card-top">
                  <h2>{p.title || p.clientName}</h2>
                  <span className={`pill pill-${statusTone(p.status)}`}>{statusText(p.status)}</span>
                </div>
                <p className="project-card-client">{p.clientName}</p>
                <div className="project-card-meta">
                  <span>
                    {milestones.length} {milestones.length === 1 ? "milestone" : "milestones"}
                    {" · "}
                    {formatMoney(totalCents(p), milestones[0]?.currency ?? "USD")}
                  </span>
                  <span className="project-card-arrow">→</span>
                </div>
              </Link>
            );
          })}
        </div>
      )}

      {items?.length === 0 && (
        <div className="empty-state">
          <h2>No projects yet</h2>
          <p>Start with the client's request — you can turn it into scoped milestones next.</p>
          <Link to="/projects/new" className="btn btn-primary">Create your first project</Link>
        </div>
      )}
    </main>
  );
}
