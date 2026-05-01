"use client";

import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { useWorkspaceId } from "@multica/core/hooks";
import { agentTaskSnapshotOptions } from "@multica/core/agents/queries";

/**
 * Returns whether the given issue currently has an active
 * (queued / dispatched / running) execution.
 *
 * Uses the workspace-wide agent-task-snapshot query — a single API call
 * shared across the entire app. The snapshot includes every active task
 * (with issue_id) plus each agent's most recent terminal task.
 * WS task lifecycle events invalidate the snapshot in useRealtimeSync,
 * so this hook stays fresh without per-card API calls.
 */
export function useIssueActiveExecution(issueId: string): boolean {
  const wsId = useWorkspaceId();
  const { data: snapshot = [] } = useQuery(agentTaskSnapshotOptions(wsId));

  return useMemo(
    () =>
      snapshot.some(
        (t) =>
          t.issue_id === issueId &&
          (t.status === "queued" ||
            t.status === "dispatched" ||
            t.status === "running"),
      ),
    [snapshot, issueId],
  );
}
