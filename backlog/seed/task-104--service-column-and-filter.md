---
# task-104
title: Show the affected service on the board and filter by it
status: todo
type: task
priority: normal
tags:
    - workshop
    - app
    - second-run
created_at: 2026-09-16T09:15:00Z
updated_at: 2026-09-16T09:15:00Z
---

The board shows the service name in the metadata line but there is no way to see only
one service's incidents. This is the take-home task: a second, genuinely different
piece of work through the same factory.

## Acceptance criteria
- `GET /api/incidents?service=checkout-api` returns only that service's incidents.
  An unknown service value returns an empty list with `matching: 0`, not an error, and
  `total` still reports the unfiltered count.
- A malformed service value (anything outside `[a-z0-9-]{1,40}`) is ignored exactly the
  way an unknown severity is ignored today.
- `GET /api/services` returns `{"services":[{"id":"checkout-api","count":1}, ...]}`
  ordered by descending count then id, counting all incidents regardless of filter.
- The board renders a service filter control, and combining it with severity and status
  filters works.
- The `task-101` count text includes the service in `describeFilter` output when a
  service filter is active.
- Tests cover the API filter, the services endpoint, the combined filter, and the
  count text.
- `npm run typecheck`, `npm run test:unit` and `npm run test:integration` pass.
