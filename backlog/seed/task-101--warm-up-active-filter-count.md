---
# task-101
title: 'Warm-up: active-filter result count'
status: todo
type: task
priority: high
tags:
    - workshop
    - warmup
    - app
created_at: 2026-09-16T09:00:00Z
updated_at: 2026-09-16T09:00:00Z
---

The board header shows `Showing 8 incident(s)`. The `(s)` is wrong for a single
result, and the text does not say that a filter is narrowing the list.

Change `resultCountText` in `src/domain/counts.ts` so the header reads correctly.

## Acceptance criteria
- With no filter active and N incidents, the text is exactly `Showing N incidents`.
- With no filter active and exactly 1 incident, the text is exactly `Showing 1 incident`.
- With a filter active, the text is exactly `Showing M of N incidents matching FILTER`,
  where M is the number of matching incidents, N is the unfiltered total, and FILTER is
  `describeFilter(filter)` from `src/domain/incident.ts`.
- With a filter active and exactly 1 match, the text is exactly
  `Showing 1 of N incidents matching FILTER` — `incidents` stays plural because it
  describes the total.
- `GET /api/incidents` returns the same string in its `countText` field.
- `tests/unit/counts.test.ts` covers zero, one and many results, both with and without an
  active filter.
- `npm run typecheck`, `npm run test:unit` and `npm run test:integration` pass.

## Out of scope
- Any change to filtering behaviour, the API shape, or the page layout.
