"use client";

import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@multica/core/api";
import { issueKeys } from "@multica/core/issues/queries";

/**
 * Subscribes to the per-issue task list and returns whether the issue
 * currently has an active (queued / dispatched / running) execution.
 *
 * The query cache is kept fresh by the global WS sync
 * (`use-realtime-sync.ts`) which invalidates `["issues", "tasks"]`
 * on every task lifecycle event — no local WS subscription needed.
 */
export function useIssueActiveExecution(issueId: string): boolean {
  const { data: tasks = [] } = useQuery({
    queryKey: issueKeys.tasks(issueId),
    queryFn: () => api.listTasksByIssue(issueId),
    staleTime: 30_000,
    refetchOnWindowFocus: true,
  });

  return useMemo(
    () =>
      tasks.some(
        (t) =>
          t.status === "queued" ||
          t.status === "dispatched" ||
          t.status === "running",
      ),
    [tasks],
  );
}
