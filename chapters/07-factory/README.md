# Chapter 7: Give your factory another project

Now try the factory on a project you care about. Choose a small feature or bug
fix in a repository you know well enough to judge the result. A library or CLI
works too. You'll keep the environment, roles and communication you've built,
then give the team a different project and task.

## 1. Choose the project

Pick a repository you can clone and a change small enough to try during the
remaining workshop time. A familiar public project is a useful starting point:
you can judge the result without first learning a new codebase.

In HOST, clone it. Replace `https://github.com/OWNER/REPOSITORY.git` with its actual
clone URL; keep `projects/my-project` as the local directory for these instructions:

```bash
mkdir -p projects
git clone --config core.autocrlf=false --config core.eol=lf https://github.com/OWNER/REPOSITORY.git projects/my-project
```

`mkdir` creates a place for your chosen project. The Git options keep text files
with Linux-compatible line endings for the tools inside SBX, including on Windows.
This working copy will be mounted directly, including its Git history. The
application from the earlier chapters remains in `sample-app/`. Work on only this
new project in the next sandbox; there is no separate output copy to retrieve.

## 2. Write a task for someone seeing the project for the first time

Create `factory/TASK.md` in your editor, replacing the bracketed sections:

```markdown
# [Short title]

## Problem
[What happens now? Give a concrete example.]

## Wanted behavior
[What should happen after the change?]

## How to check it
[One or two observable examples or focused tests.]

## Scope
[What should the team leave alone?]
```

Add it to the host backlog:

```bash
./scripts/beans create "My project's first factory task" --type task --status todo --body-file factory/TASK.md
```

Beans prints the new task's ID. In `factory/chapter.env`, replace `wad-103` on
the `TASK=` line with that ID. Keep `MODE=mcp`, `USE_ACR=1` and `SESSION=shell`; the gateway will read this task just
as it read the sample tasks. The new task lives in your workshop backlog on the
host; creating it doesn't open an issue in the project's upstream repository.

## 3. Adapt the instructions, not the whole factory

Replace `factory/PROMPT.md` with:

```markdown
Read your role brief. Use crew send for messages; it stores each handoff and wakes its recipient.
Ask the developer to retrieve the task named in ~/work/task-id through Beans MCP
and save it to ~/work/task.json before implementing anything.

This is a different project, mounted at ~/work/app. First ask the developer to read its README,
contributor guide, agent instructions, dependency manifests and build configuration.
Determine the required runtimes, dependencies and supporting services. Install and
start what is needed inside SBX, using containers for services where appropriate.
Discover and run a relevant baseline check; report any existing failure separately
from failures caused by the change. Then implement only the requested change. QA reviews the
exact commit and runs the relevant checks. Do not replace existing project guidance.

Project setup: work out and follow the repository's development setup inside SBX;
the human does not need to supply installation commands.
Use the sandbox's Docker engine if tests need containers. Determine the stack from the repository rather than assuming Node or PostgreSQL.
If a dependency, toolchain or network request is unavailable, report the exact need
to the human. Do not claim a skipped check passed.

Preview: for a library or CLI, demonstrate the change with a focused test or CLI
example; no web server is required. For a web app, use its documented start command
and listen on 0.0.0.0:8080 to use the existing host port 3102 mapping.

After review, ask the developer to append a Beans result note with the commit,
summary, actual check outcomes and QA findings. Leave the task open. Ask the human
about ambiguous requirements rather than choosing a product answer.
```

The team will read your task alongside this prompt, then inspect the repository
to work out how to build and test it. Add anything it can't learn from the code,
such as a reproduction or a design constraint. It can discover the installation
commands from the project itself.

`AGENTS.md`, if this project has one, contains instructions for coding assistants.
The setup helper keeps it. Read the project's guidance before adding another policy so you understand the conventions the agents should follow.
As you watch them set up the project, look for instructions you'd want to reuse as
a skill or tools you'd want to install through a kit next time.

## 4. Let the team work

In SANDBOX, starting at the workshop repository root:

```bash
./scripts/launch-factory.sh wad-my-project projects/my-project
```

The first argument names the sandbox. The second changes its mounted workspace.
The launcher still uses your `factory/sbxenv.yaml`, kits, team and MCP connection.
Inside its shell:

```bash
crew submit
crew watch
```

Follow the team as it discovers the project. Ctrl-C returns to the shell, where
you can ask a follow-up:

```bash
crew ask "Explain which checks demonstrate the requested behavior and which commit QA reviewed."
crew watch
```

Use `crew reply "your answer"` for a product question. Use HOST for a scoped network
policy change, as in chapter 5. `crew logs developer` shows the assistant's current
terminal output if a tool or login needs attention; it does not attach an interactive
session to that assistant.

## 5. Inspect the result in your working copy

Try the changed behavior: a web app uses <http://127.0.0.1:3102>; a library or CLI
should have a focused demonstration inside SBX. On HOST:

```bash
git -C projects/my-project status --short
git -C projects/my-project log -1 --oneline
./scripts/beans show YOUR-TASK-ID
```

Replace `YOUR-TASK-ID` with your task's ID. Open the changed files in your editor
and compare the behavior with your requirements. The existing history and new
commits are in this working copy. Decide what to keep through your normal review
process; the workshop does not publish changes upstream.

When finished, press Ctrl-C to leave `crew watch`, then type `exit` to leave
the SANDBOX shell. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=wad-my-project --env-arg "app=$(pwd)/projects/my-project"
```

SBX 0.50.0 also removes this environment's local MCP command. If an older or
manually registered `wad-my-project-beans` entry remains in `sbx mcp ls`, remove it
with `sbx mcp rm wad-my-project-beans`. The app argument identifies the same
workspace used on creation. You now have an
environment and team you can reuse: change the task and select a project for the
next job.

Next: [presenter demonstrations of further SBX capabilities](../08-presenter/README.md).
