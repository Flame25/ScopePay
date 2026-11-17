import { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { type ApiError } from "../api/client";
import { useAuth } from "../auth/AuthContext";
import { listProjects } from "../api/projects";

type Project = { id: number; clientName: string; title: string; status: string };

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

  return (
    <main className="page">
      <header className="page-head">
        <h1>Projects</h1>
        <div className="head-actions">
          <span className="muted">{user?.email}</span>
          <Link to="/projects/new" className="btn-primary">New project</Link>
          <button className="btn-primary" onClick={() => { logout(); navigate("/login"); }}>Logout</button>
        </div>
      </header>

      {error && <p className="form-error">{error}</p>}
      {!items && !error && <p className="muted">Loading…</p>}

      <div className="card-container">
        {items?.map((p) => (
          <Link to={`/projects/${p.id}`} className="card card-link" key={p.id}>
            <strong>{p.title || p.clientName}</strong>
            <p className="muted">{p.clientName} · {p.status}</p>
          </Link>
        ))}
        {items?.length === 0 && <p className="muted">No projects yet.</p>}
      </div>
    </main>
  );
}
