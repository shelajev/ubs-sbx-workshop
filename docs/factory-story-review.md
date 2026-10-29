# Factory story: content map and review record

Maintainer record for the attendee guide in `workshop/`. Source baseline:
`bb9a046` on main. The attendee entry point is `README.md` → `workshop/README.md`.

## Map of substantial text

| Attendee block | Source blocks reused | Substantial new or adapted text to review |
|---|---|---|
| Root README and workshop map | Root README, chapter guide | Capability-led opening, 90-minute checkpoint, sequential navigation and takeaways |
| 00 setup | Chapter 0 almost whole | Claude/Codex choice, step links and team continuation prerequisite |
| 01 sandbox | Chapter 1 §§1–3 | Codex command adaptation, early policy inspection, retaining the sandbox for task work |
| 02 tasks | Chapter 2 Beans explanation; chapter 1 coding task; chapter 5 tool boundary | One explained preparation command, paste/commit/retrieve comparison, manual SBX MCP registration and loading, cleanup |
| 03 kits | ACR interlude hello kit, publisher approval, ACR explanation and policy review | Edit-and-recreate exercise, manual two-kit creation, direct agent prompt replacing launcher-dependent guidance installation, review via MCP |
| 04 environment | Chapter 2 environment-file block; chapter 5 MCP declaration | Complete composed environment, Windows host shell setup, fresh-worker verification, session checkpoint |
| 05 team | Chapter 4 almost whole | Launcher introduced here with each input file explained, preserving existing kits/MCP, role prerequisites and cleanup |
| 06 feature | Chapter 5 §§3–4 and 6 | Read-only-to-notes permission transition, existing MCP connection, explicit catch-up and return from optional browser exercise |
| 07 human | Chapter 6 almost whole | Feature prerequisite, sequential links and terminology |
| 08 next | Chapter 7 / mixed-model links | Explain host dispatch as a future extension, not a tool already supplied |
| Presenter notes | Existing delivery guidance | Observable checkpoints, kit priority, how to resume and continue |

## Review brief

Read as a normal attendee from a fresh clone. Check that commands have a purpose,
location and expected result; that every file, helper, task and sandbox exists
before use; that cleanup preserves required next-step state; and that the narrative
contains no revision history or outside planning discussion. Prioritize lost readers,
broken commands, hidden prerequisites and unclear value over cosmetic rewriting.
Review the substantial blocks above and their transitions together.

## Validation and Gemini review

Completed on October 29, 2026:

- Local Markdown links and heading anchors across the root and attendee guide passed.
- The kit and environment YAML were extracted from the attendee code blocks.
  The real installed SBX 0.49.2 client accepted the hello kit and both Claude and
  Codex environment plans, including kits, workspace, port and read-only MCP.
- A fresh Claude sandbox composed the authored hello kit with the pinned ACR kit.
  `workshop-hello` printed its greeting; ACR 0.2.1 installed the pinned policy,
  realized 11 changes for Claude/Codex, and reported the realization current.
- Manual `sbx mcp add` and `load` connected the real Beans adapter to that running
  sandbox. Claude retrieved task-101 through the gateway, reported its title and
  acceptance criterion, and confirmed no note-writing tool was exposed.
  To honor the operator's host execution policy, this rehearsal ran the adapter
  and disposable backlog in a second SBX via `sbx exec -i`; it did not execute
  the host adapter command printed for attendees on this laptop.
- In SBX, the material download, explained setup helper, repeated setup, task
  read, shell syntax, material-mode tests, dependency-storage tests, guidance
  tests and team-default/launcher regression tests passed.
- The 0.50.0 Linux CLI was downloaded and its `env run` and `mcp load` interfaces
  inspected inside SBX. Environment planning there required Docker login, so
  actual plans used the already authenticated host client. No credentials were
  copied into the validation sandbox.

These checks do not establish a complete 0.50.0 daemon rehearsal, native Windows
MCP process launch, a Codex model session, or the entire team coding exercise.
The runtime and team helpers are reused from main; only their success-message
pointers changed. The content map identifies the changed attendee instructions.

### Gemini review and disposition

Gemini CLI 0.65.0 reviewed the complete assembled guide inside a mountless SBX,
using the content map above. A second pass reviewed the final attendee text after
fixes and returned **“No remaining blockers/major issues.”**

- Removed the unused `MODE=team` change before `MODE=mcp`.
- Fixed the materials helper's completion link to the attendee entry page.
- Made the Windows `host_shell` dependency explicit at the team transition.
- Moved the detailed session timing out of the attendee map; kept necessary
  installation prerequisites and troubleshooting guidance for attendees.
- Quoted the optional note flag in the YAML example after the final minor review.
- Kept `agent: claude` at the team transition: Codex users may have changed it.
- Kept `brew trust docker/tap`: the suggestion to replace it with `brew tap`
  conflicts with the current [official Docker installation instructions](https://docs.docker.com/ai/sandboxes/install/#install-on-macos).

Raw first and final reviews are retained locally under
`.local/reviews/factory-story-gemini.txt` and `factory-story-gemini-final.txt`.
No Gemini process ran on the host, and no host credentials were mounted.

## Attendee rehearsal follow-up

October 29: completed the manual sandbox, shared-file boundary, Docker and published
HTTP port exercises; pasted task-101; retrieved the same task through SBX-managed
MCP; changed and recreated a kit; installed and checked ACR guidance; and recreated
the environment. Fresh Claude sandboxes reused the existing configured model
access without another login.

The three Claude roles completed a read-only handoff and implemented task-102,
reviewed it and appended a Beans note. SSH joined the running team, delivered a
human decision and observed its recorded acknowledgement. The latter rehearsed
only the handoff, not implementation of task-103.

Rehearsal fixes:

- Ask attendees to inspect files and find rows in the full policy table, without
  supplying the observations in advance.
- Tell the agent about the already published host port 3102.
- Prepare dependency storage on Windows before the fresh ACR sandbox runs checks.
- Explain Claude's background-task confirmation when exiting.
- Make browser-download recovery reachable during feature work.
- Tell the coordinator to advance the existing implementation/review/correction
  stages. This prompt change was reviewed, but not rerun as another feature task.

Application checks run by the agents passed, but the feature's browser flow was
explicitly unverified after a blocked browser download. With operator approval,
Playwright and native dependencies were then installed and Chromium launched
successfully inside SBX. No further application suite was run. No computer use
or desktop UI inspection was performed.

The actual host client remained SBX 0.49.2. This does not establish a full 0.50.0,
native Windows or Codex model rehearsal. Host helper execution was not approved:
setup, Beans and its adapter ran in a separate SBX, and the launcher's transfer and
prepare operations were reproduced through the SBX control plane. The exact host
helper/adapter invocation therefore remains unverified on this laptop. The Windows
dependency-storage branch was exercised in Linux SBX, not on Windows.

Gemini reviewed the final fixes inside SBX and returned **“Approved”**, with no
blockers. Raw review: `.local/reviews/factory-walkthrough-final-gemini.txt`.
Rehearsal transcripts and team decision/message evidence remain in the disposable
local walkthrough checkout; they are not attendee materials.
