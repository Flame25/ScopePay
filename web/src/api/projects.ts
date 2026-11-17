import { request, requestVoid } from "./client";
import { authHeaders } from "./auth";

export type ProjectStatus = "draft" | "sent" | "approved" | "in_progress" | "complete";

export type Milestone = {
  id: number;
  projectId: number;
  order: number;
  title: string;
  deliverables: string;
  acceptanceCriteria: string;
  amountCents: number;
  currency: string;
  status: string;
}

export type Project = {
  id: number;
  clientName: string;
  clientEmail: string;
  rawRequest: string;
  title: string;
  scopeSummary: string;
  status: ProjectStatus;
  milestones?: Milestone[];
  createdAt: string;
  updatedAt: string;
}

export function listProjects(): Promise<Project[]> {
  return request<Project[]>("/projects", {
    method: "GET",
    headers: authHeaders(),
  });
}

export function getProject(id: number): Promise<Project> {
  return request<Project>(`/projects/${id}`, {
    method: "GET",
    headers: authHeaders(),
  });
}

export function createProject(body: { clientName: string; clientEmail: string; rawRequest: string }): Promise<Project> {
  return request<Project>("/projects", {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify(body),
  });
}

export function updateProject(id: number, body: Partial<Pick<Project, "clientName" | "clientEmail" | "rawRequest" | "title" | "scopeSummary">>): Promise<Project> {
  return request<Project>(`/projects/${id}`, {
    method: "PATCH",
    headers: authHeaders(),
    body: JSON.stringify(body),
  });
}

export function sendProject(id: number): Promise<{ shareUrl: string }> {
  return request<{ shareUrl: string }>(`/projects/${id}/send`, {
    method: "POST",
    headers: authHeaders(),
  });
}

export function createMilestone(projectId: number, body: { title: string; deliverables?: string; acceptanceCriteria?: string; amountCents: number; currency?: string }): Promise<Milestone> {
  return request<Milestone>(`/projects/${projectId}/milestones`, {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify(body),
  });
}

export function updateMilestone(projectId: number, mid: number, body: Partial<{ title: string; deliverables: string; acceptanceCriteria: string; amountCents: number; currency: string; order: number; }>): Promise<Milestone> {
  return request<Milestone>(`/projects/${projectId}/milestones/${mid}`, {
    method: "PATCH",
    headers: authHeaders(),
    body: JSON.stringify(body),
  });
}

export function deleteMilestone(projectId: number, mid: number): Promise<void> {
  return requestVoid(`/projects/${projectId}/milestones/${mid}`, {
    method: "DELETE",
    headers: authHeaders(),
  });
}

