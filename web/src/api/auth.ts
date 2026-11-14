import { request } from "./client";

export type User = { id: number; email: string; name: string };
const TOKEN_KEY = "scopepay_token";

export function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(t: string | null) {
  if (t) {
    localStorage.setItem(TOKEN_KEY, t);
  } else {
    localStorage.removeItem(TOKEN_KEY);
  }
}

export function authHeaders(): HeadersInit {
  const t = getToken();
  return { "Content-Type": "application/json", ...(t ? { Authorization: `Bearer ${t}` } : {}) }
}

export async function register(email: string, name: string, password: string): Promise<User> {
  return request<User>("/auth/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, name, password }),
  });
}

export async function login(email: string, password: string): Promise<User> {
  const res = await request<{ token: string, user: User }>("/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  setToken(res.token);
  return res.user;
}

export async function me(): Promise<User> {
  return request<User>("/auth/me", { headers: authHeaders() });
}

export function logout() {
  setToken(null)
}
