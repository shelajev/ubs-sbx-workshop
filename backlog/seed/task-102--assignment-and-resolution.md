---
# task-102
title: Add incident assignment and resolution with a required note
status: todo
type: feature
priority: critical
tags:
    - workshop
    - app
    - main-feature
created_at: 2026-09-16T09:05:00Z
updated_at: 2026-09-16T09:05:00Z
---

Responders need to take ownership of an incident and close it out with a note that
explains what happened. The board currently has no way to do either, and no history.

Implement assignment, resolution and a durable history, validated on the server.

## API contract

`GET /api/assignees`
- `200 {"assignees": [{"id": "...", "name": "..."}]}`
- The roster is fixed application data, not a database table: put it in
  `src/domain/assignees.ts`. Ids are lowercase dotted names, e.g. `rae.mendez`.

`POST /api/incidents/:id/assign`
- Request: `{"assignee": "rae.mendez"}`
- `200 {"incident": {...}}` on success; the incident's `assignee` is set.
- Status transition, stated exactly because the obvious wording is ambiguous: an `open`
  incident becomes `acknowledged` **only when a non-null known assignee is set**.
  Clearing an assignment never changes the status — an `open` incident stays `open` and
  an `acknowledged` one stays `acknowledged`. `acknowledged` means a person took
  ownership, and clearing the owner is the opposite of that.
- A `resolved` incident keeps status `resolved` whatever the assignee becomes.
- `{"assignee": null}` clears the assignment and is a valid request.
- `400 {"error":{"code":"unknown_assignee","field":"assignee"}}` if the id is not in the roster.
- `400 {"error":{"code":"invalid_payload","field":"assignee"}}` if `assignee` is missing
  or not a string or null.
- `404 {"error":{"code":"not_found","field":"id"}}` for an unknown incident id.
- `400 {"error":{"code":"invalid_id","field":"id"}}` for a malformed incident id.

`POST /api/incidents/:id/resolve`
- Request: `{"note": "text", "requestId": "optional-client-token"}`
- `200 {"incident": {...}}` on success; status becomes `resolved`, `resolutionNote`
  holds the trimmed note, and one `resolved` history entry is recorded.
- `400 {"error":{"code":"note_required","field":"note"}}` if `note` is missing, not a
  string, empty, or whitespace-only. This is enforced in the API. A browser-only guard
  does not satisfy this criterion.
- `400 {"error":{"code":"note_too_long","field":"note"}}` if the trimmed note is longer
  than 500 characters.
- `409 {"error":{"code":"already_resolved"}}` if the incident is already resolved and
  the request carries no `requestId`, or carries a `requestId` that does not match the
  one that resolved it.
- Idempotency: repeating the same request with the same `requestId` returns `200` with
  the unchanged incident and does **not** add a second history entry.
- `requestId`, when present, must be a string of 8 to 64 characters matching
  `[A-Za-z0-9._-]+`; otherwise `400 {"error":{"code":"invalid_request_id","field":"requestId"}}`.

`GET /api/incidents/:id`
- Response gains `assignee` (string or null), `resolutionNote` (string or null) and
  `history`: an array ordered oldest-first of
  `{"id": number, "kind": "assigned"|"unassigned"|"resolved", "at": ISO8601, "actor": string, "detail": {...}}`.

## UI

- The detail page shows the current assignee, an assignee picker, and a resolve form
  with a note field.
- The detail page shows the history, oldest first, with the note text visible on the
  `resolved` entry.
- After resolving, a page refresh still shows the resolved status, the note and the
  history entry: the state is in Postgres, not in the page.
- The board list shows the assignee when one is set.

## Acceptance criteria
- Every status code and error code above is exercised by a test in `tests/api/`.
- A new migration in `db/migrations/` adds the `assignee` and `resolution_note` columns
  and an `incident_events` table. The existing migration is not edited.
- History survives a reconnect: a fresh pool reads the same entries.
- Repeating a resolve with the same `requestId` leaves exactly one `resolved` entry.
- Existing filter behaviour and the `task-101` count text are unchanged.
- One browser flow assigns and resolves an incident, reloads the page, and finds the
  note and history still there.
- `npm run typecheck`, `npm run test:unit`, `npm run test:integration` and
  `npm run test:api` pass.

## Out of scope
- Authentication, real user accounts, notifications, and editing history entries.
- Reopening a resolved incident: that is `task-103`.

## Note for maintainers

The assignment transition is specified above so this task can proceed without a
product decision. `task-103` is the intentional human-decision exercise.
