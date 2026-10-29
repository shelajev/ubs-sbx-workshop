---
# task-103
title: Reopen a resolved incident
status: todo
type: feature
priority: high
tags:
    - workshop
    - app
    - needs-decision
created_at: 2026-09-16T09:10:00Z
updated_at: 2026-09-16T09:10:00Z
---

An incident was resolved too early and has to go back on the board.

`POST /api/incidents/:id/reopen`
- Request: `{"reason": "text", "requestId": "optional-client-token"}`
- `200 {"incident": {...}}`; status becomes `open`, and a `reopened` history entry is
  recorded with the reason in `detail`.
- `400 {"error":{"code":"reason_required","field":"reason"}}` for a missing, empty or
  whitespace-only reason.
- `409 {"error":{"code":"not_resolved"}}` if the incident is not currently resolved.
- The same `requestId` rules as `task-102` apply.

## Unresolved requirement — do not guess

The brief does not say what happens to `resolutionNote` when an incident is reopened,
and the two defensible readings produce different API responses and different tests:

- **Position A (keep):** `resolutionNote` stays as it is. It is a record of what the
  last responder concluded, and the `reopened` history entry already shows the state
  changed. Reopening loses no information.
- **Position B (clear):** `resolutionNote` becomes `null`, because an open incident has
  no current resolution and leaving the old note there misleads the next responder. The
  note is still readable in the `resolved` history entry, so nothing is lost.

A coordinator must not pick one of these. Report `needs-human` with both positions and
the exact field in question, and wait for the recorded decision.

## Browser workflow

On a resolved incident's detail page, provide a reopen form requiring a reason.
Submit it through the reopen API and display any validation error. After success,
show the open status, current resolution note according to the recorded decision,
and history including the reopen reason. Keep the earlier resolution note visible
in its historical entry. The reopen form is only available for resolved incidents.

## Acceptance criteria
- All of the status codes above are covered by `tests/api/` tests.
- The behaviour of `resolutionNote` on reopen matches the decision recorded in
  `$FACTORY_DIR/decisions/<task>-a<attempt>.json` in the run workspace, and the decision file
  is referenced in the implementation report.
- A test asserts the decided behaviour explicitly, naming the decision in the test title.
- `task-102` behaviour and its tests are unchanged.
- A real-browser test resolves an incident with a note, reopens it with a reason,
  reloads the page, and verifies the decided current-note behaviour and retained history.
- `npm run typecheck`, `npm run test:unit`, `npm run test:integration` and
  `npm run test:api` pass. `npm run test:browser` passes with no skipped tests.
