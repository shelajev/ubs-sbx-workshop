# Chapter 2: Make the sandbox repeatable

Your warm-up change is in `sample-app/`. Let's give a new agent that project and
a written task without retyping the sandbox settings. You'll put those settings
in an environment file, store the task in Beans, and use a short host script to
bring them together. The new sandbox will mount the same application directory.

## 1. Describe the environment

An [SBX environment file](https://docs.docker.com/ai/sandboxes/configuration/environment-files/)
records the settings you previously typed into `sbx run`. In HOST, create the
configuration directory:

```bash
mkdir -p factory
```

Create `factory/sbxenv.yaml` in your editor:

```yaml
schemaVersion: "1"
name: sbx-env-first
agent: claude
workspace: ../sample-app
sandboxOptions:
  skills: off
  cpus: 4
  memory: 8g
ports:
  - sandbox: 8080
    host: 3102
    protocol: tcp4
    hostIP: 127.0.0.1
```

`agent` selects the built-in Claude environment. `workspace` is relative to this
file: `../sample-app` shares our application, including `.git`. The resource and
skills settings are the choices from chapter 1. `ports` publishes the app's port
8080 at localhost:3102 every time this environment is created. It does not start
the app; the agent does that when needed.

Preview it in HOST:

```bash
sbx env plan factory/sbxenv.yaml
```

Find your sample-app path and the 3102 → 8080 mapping. A plan shows the proposed
setup without creating it. In your SANDBOX tab, still at the repository root:

```bash
sbx env run factory/sbxenv.yaml
```

Read and approve the plan. This creates the environment and opens Claude.
Ask the assistant to inspect the project and confirm your filter-count change is
present. Claude is reading the files and Git history from your host, including
the change you made in the previous chapter.

Type `/exit` to leave Claude. In HOST:

```bash
sbx env rm factory/sbxenv.yaml
```

Approve removal. The source stays in `sample-app/`. Removing this practice
sandbox lets us apply new configuration on the next creation.

## 2. Put work in a task tracker

[Beans](https://github.com/hmans/beans) is a command-line issue tracker. Its tasks
are Markdown files: a person can read them in an editor, and an agent can read
requirements, update notes and use the CLI without a browser integration.

In an ordinary project, `beans init` creates the task store, `beans create` adds
work, and `beans show` reads it. Our backlog will stay **on the host**, outside the
mounted app. Later, MCP will expose selected operations to the agents.

In HOST, install the pinned CLI and initialize the supplied workshop tasks:

```bash
./scripts/install-beans.sh
./scripts/backlog-init.sh --disposable
```

The first helper downloads Beans. The second creates our four demo tasks under
`.local/chapters/beans/`. `--disposable` marks this as practice data that our later
MCP server may append notes to. Chapters 1–7 walk through `task-101` to `task-103`;
`task-104` is an optional take-home or presenter exercise for extending the incident
table with a service column and filter.

Try the CLI:

```bash
./scripts/beans list
./scripts/beans show task-101
./scripts/beans create "Try a small improvement" --type task --status todo
```

`list` shows the backlog; `show` prints the warm-up requirements; `create` adds a
new task and prints its ID. Use `./scripts/beans show YOUR-ID` with that ID to
read it. Open `.local/chapters/beans/.beans/` in your editor to see the Markdown
behind the commands.

`scripts/beans` is a small wrapper around the real Beans CLI. It selects this
workshop's executable and backlog so each command need not repeat their paths.
It does not run an agent. We will keep using the supplied tasks for the walkthrough.

## 3. Connect a task to a sandbox

We want a reusable environment with a different name for each chapter. At the top
of `factory/sbxenv.yaml`, replace the `name` and `workspace` values and add `args`.
Keep `sandboxOptions` and `ports` as they are:

```yaml
schemaVersion: "1"
name: "${{ env.args.name }}"
agent: claude
workspace: "${{ env.args.app }}"
args:
  name:
    required: true
  app:
    default: ../sample-app
```

The name is an input now. The app path defaults to the same working directory;
chapter 7 will show how to select another project.

Create `factory/chapter.env`:

```text
TASK=task-101
MODE=manual
USE_ACR=0
SESSION=claude
```

`sbxenv.yaml` configures SBX. This smaller file configures the workshop's host
launcher, which connects a task to that environment:

| Setting | What it tells the launcher |
|---|---|
| `TASK=task-101` | Read the filter-count task from our host backlog. |
| `MODE=manual` | Let us work with the assistant ourselves. |
| `USE_ACR=0` | Leave coding-policy installation for the next chapter. |
| `SESSION=claude` | Open Claude after preparing the sandbox. |

Create `factory/PROMPT.md`:

```markdown
Read ~/work/task.json and inspect this project. Confirm whether the warm-up
filter-count change is already present; do not implement it twice. Explain
how the current code meets the task and what you would check.
```

The agent will read this prompt alongside the task. The task describes the
filter-count change; the prompt asks the agent to check the work you already did.
Your environment file supplies the workspace where it can inspect that code.

## 4. Start the repeatable sandbox

Preview the configuration in HOST:

```bash
sbx env plan factory/sbxenv.yaml --env-arg name=sbx-ch-02
```

Then in the SANDBOX tab, at the workshop root:

```bash
./scripts/launch-factory.sh sbx-ch-02
```

The launcher creates `sbx-ch-02` and opens Claude in `sample-app/`. It gives Claude
the selected task as `~/work/task.json` and your prompt as `~/work/PROMPT.md`.
The launcher does not implement or verify the task for you. Tell Claude:

> Read ~/work/PROMPT.md and follow those instructions. Explain what you found.

You should recognize your own warm-up change. In your host editor, the same files
and Git history are still in `sample-app/`.

Type `/exit` to end the session. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=sbx-ch-02
```

Next: [give each worker shared guidance through a kit](../02.5-acr/README.md).
