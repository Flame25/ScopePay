import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { request, type ApiError } from "../api/client";
import { authHeaders } from "../api/auth";
import { useAuth } from "../auth/AuthContext";

type Project = { id: number; clientName: string; title: string; status: string };

export default function Projects() {
  const [items, setItems] = useState<Project[] | null>(null);
  const [error, setError] = useState<string | null>(null);
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    request<Project[]>("/projects", { headers: authHeaders() })
      .then(setItems)
      .catch((e) => setError((e as ApiError).error?.message ?? "failed to load"));
  }, []);

  return (
    <main className="page">
      <header style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <h1> Projects </h1>
        <span className="muted">{user?.email} </span>
        <button className="btn-primary" onClick={() => { logout(); navigate("/login"); }}>
          Logout
        </button>
      </header>
    </main >
  );
}
