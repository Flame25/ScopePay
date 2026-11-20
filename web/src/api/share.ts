import { authHeaders } from "./auth";
import { request } from "./client";
import type { Milestone, ProjectStatus } from "./projects";


export type ShareViewerRole = "client" | "freelancer";

export type ShareProject = {
  title: string;
  clientName: string;
  scopeSummary: string;
  status: ProjectStatus;
  editRequestNote: string;
  viewerRole: ShareViewerRole;
  milestones: Milestone[];
};

export function getShare(token: string): Promise<ShareProject> {
  return request<ShareProject>(`/share/${token}`, { headers: authHeaders() });
}

export function approveShare(token: string): Promise<ShareProject> {
  return request<ShareProject>(`/share/${token}/approve`, { method: "POST", headers: authHeaders() });
}

export function requestChanges(token: string, note: string): Promise<ShareProject> {
  return request<ShareProject>(`/share/${token}/request-changes`, {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({ note }),
  });
}
