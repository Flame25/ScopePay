import type { ProjectStatus } from "../api/projects";

export const statusLabel: Record<ProjectStatus, string> = {
  draft: "Draft",
  sent: "Awaiting approval",
  approved: "Approved",
  in_progress: "In progress",
  complete: "Complete",
  revision_requested: "Changes requested",
};

/** Pill modifier: known statuses map to a tone, anything unknown falls back to draft. */
export function statusTone(status: string): ProjectStatus {
  return status in statusLabel ? (status as ProjectStatus) : "draft";
}

export function statusText(status: string): string {
  return statusLabel[statusTone(status)];
}
