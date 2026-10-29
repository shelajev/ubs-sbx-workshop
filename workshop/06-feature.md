# 6. Give the team a feature

The team can pass a question through development and review. Now give it a task
that changes the app. It will retrieve requirements through the connection you
already built, implement them, and leave evidence you can inspect.

## 1. Enable a result note

Read the next task in HOST:

```bash
./scripts/beans show wad-102
```

It adds incident assignment and resolution. The task's acceptance criteria are
what you and QA will use to evaluate the result.

In `factory/sbxenv.yaml`, change only the MCP server's `args` line to:

```yaml
      args: ["${{ env.fileDir }}/../scripts/beans-mcp", "--allow-notes"]
```

`--allow-notes` adds one operation: append a note to a task in the disposable
workshop backlog. The agent still cannot delete tasks or run arbitrary host
commands through this adapter. This is an explicit increase from read-only access;
review it in the next environment plan.

Replace `factory/chapter.env` with:

```text
TASK=wad-102
MODE=mcp
USE_ACR=1
SESSION=shell
```

`MODE=mcp` tells the workshop launcher to start the team and supply a task ID
instead of a copied task. It is a launcher setting, not an SBX requirement for
using MCP: your single agent already used the gateway without it.

Replace `factory/PROMPT.md` with:

```markdown
Ask the developer to retrieve the task named in ~/work/task-id through the Beans
MCP get_task tool and save it to ~/work/task.json. Read its full body, including
every acceptance criterion. Do not invent its requirements.
The developer should inspect the mounted project's documentation, install its
dependencies and start any required services inside SBX, then implement the task.
Ask QA to review the exact commit against the task and coding policy.
After review, refresh the web app on 0.0.0.0:8080 and leave it running.
Ask the developer to append a Beans result note with the change, commit and actual
check outcomes. Leave the task open. Send the human a summary and how to try it.
```

The task states the desired behavior. This prompt describes the work process:
retrieve, implement, review, check and report. The coordinator passes it between
roles; you do not have to copy the developer's answer into QA yourself.

## 2. Start, then submit

Preview the environment from HOST:

```bash
sbx env plan factory/sbxenv.yaml --env-arg name=wad-ch-05
```

Look for the project mount, Herdr and ACR kits, and the MCP command with its new
note permission. In the SANDBOX tab, at the workshop root:

```bash
./scripts/launch-factory.sh wad-ch-05
```

Read and approve the plan and the host MCP command. The launcher starts the
sessions and opens a sandbox shell. Give the selected task to the coordinator:

```bash
crew submit
```

Submission waits for the roles to be ready. Now follow the work in the same shell:

```bash
crew watch
```

Expect the developer to fetch the task, implement it and pass a specific commit
to QA. Ctrl-C leaves the message view without stopping the agents. If the view
is quiet, use `crew status` to see the stage, then `crew logs developer` (or the
role that should act next) to inspect its terminal. Check for authentication,
active work or an error before sending another request.

## 3. Inspect the result

When the team reports readiness, open <http://127.0.0.1:3102>. Try assigning and
resolving an incident. In HOST, read the task again:

```bash
./scripts/beans show wad-102
```

Look for the result note, including the commit and actual check outcomes. Inspect
the changed files in your editor. A note saying “done” is not a substitute for
trying the behavior and reading the evidence.

For an optional exercise in approving a browser download and saving the resulting
network rules as a kit, use [browser access](../chapters/05-mcp/README.md#5-let-the-human-handle-a-request-for-more-access).
Do only that section, then return here. It uses this running `wad-ch-05` team and
your existing `factory/` configuration.

## 4. Finish this environment

When the team is done, Ctrl-C leaves the message view and `exit` leaves the
sandbox shell. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=wad-ch-05
```

Review the destroy plan. SBX 0.50.0 also removes the environment's local MCP
command. If recovering an older installation, inspect `sbx mcp ls` and remove
only a remaining `wad-ch-05-beans` with `sbx mcp rm wad-ch-05-beans`.

If you want to continue without waiting for implementation, finish cleanup first,
then install the completed feature checkpoint from HOST:

```bash
./scripts/prepare-app.sh app-02-feature-solution
```

This backs up your app before replacing it and prints the backup location. It
keeps your factory configuration. The checkpoint supplies the feature needed by
the next task. At the next launch, if the helper reports that the review skill
is missing, run `install-guidance` in the sandbox shell before `crew submit`.

Next: [keep decisions human](07-human.md).
