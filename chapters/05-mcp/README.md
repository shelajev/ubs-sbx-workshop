# Chapter 5: Give the team controlled access to host tools

Each sandbox has started with a copy of its task. We want the team to read the
current requirements from Beans and leave a result note when it finishes, so it
needs access to the backlog on your host. That backlog sits outside the mounted
application directory. In this chapter, you'll connect it through MCP.

[MCP](https://modelcontextprotocol.io/docs/getting-started/intro) lets an assistant
discover and call a server's named tools. [SBX's MCP gateway](https://docs.docker.com/ai/sandboxes/mcp-gateway/)
connects the assistant inside the sandbox to our small host server:

```text
Claude developer → SBX MCP gateway → host Beans MCP server → workshop backlog
```

## 1. Understand the tool boundary

The workshop includes a prebuilt MCP adapter for Beans. It exposes listing tasks,
reading a task and, when enabled, appending a note. It does not expose arbitrary
shell commands or task deletion. We can make the integration useful without giving
the agent our entire host filesystem or connecting a production issue tracker.

Chapter 0 downloaded the adapter for your computer. Install it into the
workshop's tools directory from HOST:

```bash
./scripts/install-mcp.sh --from-local
```

`--from-local` uses that earlier download. The installer checks its checksum,
places the executable in `.local/chapters/bin/` and prints its version. It is
prebuilt for your operating system, so you do not need to install Go or compile it.

The gateway will start this program on your host and exchange MCP messages through
its standard input and output (**stdio**). This gives the sandbox a connection to
Beans without needing a public web service or another account login.

Open `scripts/beans-mcp` in your editor. This short entry point selects the same
workshop backlog as `scripts/beans`. By default it is read-only; `--allow-notes`
enables the result-note tool. The underlying adapter also requires the disposable
backlog marker created in chapter 2.

## 2. Declare the connection alongside the sandbox

Add `host_shell` to the existing `args` mapping in `factory/sbxenv.yaml`, then add
the `mcp` section:

```yaml
args:
  # Keep the existing app and name arguments.
  host_shell:
    default: bash
mcp:
  servers:
    - name: "${{ env.args.name }}-beans"
      command: "${{ env.args.host_shell }}"
      args: ["${{ env.fileDir }}/../scripts/beans-mcp", --allow-notes]
```

Read this as a host command: run Bash, give it `scripts/beans-mcp`, and enable
notes with `--allow-notes`. `${{ env.fileDir }}` is the directory containing our
environment file, so the script path leads from `factory/` back to `scripts/`.
On Windows, the launcher supplies the path to Git Bash as `host_shell`.

The server name combines the sandbox name with `-beans`. For this chapter, the
registration will be called `sbx-ch-05-beans`; you'll see that name when you
inspect the connection and remove it at the end.

These fields describe both registration and attachment. You can also register and
load servers with `sbx mcp add` and `sbx mcp load`; recording them here makes the
connection part of every new factory environment.

The MCP adapter runs on your host with your user's permissions. Before starting
it, look again at the tools it exposes: we rely on the adapter's code to restrict
requests to those operations on the selected backlog. The gateway connects Claude
to the program. Each request must still be handled within the limits the adapter
defines.

## 3. Tell the team when to use its new tools

Edit `factory/chapter.env` to contain:

```text
TASK=task-102
MODE=mcp
USE_ACR=1
SESSION=shell
```

`TASK=task-102` selects the next feature: assigning incidents and recording their
resolution. `MODE=mcp` starts the team and supplies only that task ID, so the agent
must fetch the current requirements through the gateway. `SESSION=shell` gives us
a place to explore the tools before submitting work.

Replace `factory/PROMPT.md` with:

```markdown
Ask the developer to retrieve the task named in ~/work/task-id through the Beans
MCP get_task tool and save it to ~/work/task.json. Do not invent its requirements.
The developer should inspect the mounted project's documentation, install its
dependencies and start any required services inside SBX, then implement the task.
Ask QA to review the exact commit against the task and coding policy.
After review, refresh the web app on 0.0.0.0:8080 and leave it running.
Ask the developer to append a Beans result note with the change, commit and actual
check outcomes. Leave the task open. Send the human a summary and how to try it.
```

With the server in the environment file, Claude can call its tools. The prompt
gives the team a reason to use them: fetch this task, do the work and write back
the result. Claude is the developer connected to the gateway, so the coordinator
asks it for the requirements through the messages you tried in chapter 4. This
uses the all-Claude team from that chapter; no new model account is needed.

In HOST, preview the connection:

```bash
sbx env plan factory/sbxenv.yaml --env-arg name=sbx-ch-05
```

Look for the host program, note permission and application mount. In SANDBOX,
starting at the workshop root:

```bash
./scripts/launch-factory.sh sbx-ch-05
```

SBX shows an approval plan before creation. Review the workspace, kits and host
MCP command, then approve the declared access. In SBX 0.50.0, environment-defined
MCP commands also require host-command approval. Approve only this workshop's
Beans wrapper after checking its path and arguments. Keep `scripts/`, `factory/`
and the host backlog outside the writable `sample-app/` mount. Do not use
`--skip-host-commands` for this exercise; it disables the environment's MCP commands.

## 4. Explore the actual MCP tools

In the SANDBOX shell, open a personal Claude conversation:

```bash
claude
```

This is your exploration session, separate from the team's developer. Type `/mcp`
to inspect the gateway connection, then ask:

> Use the Beans MCP tools to list the workshop tasks and read task-102. Explain the
> requested feature. Do not implement it or write a note yet.

Watch the tool calls. You should see the assignment-and-resolution requirements
stored in the host backlog. Ask what information is required to resolve an incident.
Compare Claude's answer with the original task in HOST:

```bash
./scripts/beans show task-102
```

The requirements should match. You have read a host task from inside SBX without
mounting the backlog directory.

Type `/exit` to return to the SANDBOX shell. Now give the job to the team:

```bash
crew submit
crew watch
```

Submission waits for the roles to be ready. Expect the developer to fetch the task,
implement it, and pass it through QA. If a role needs attention, `crew logs developer`
(or `qa` or `coordinator`) shows its terminal. Ctrl-C leaves the message view;
it does not stop the agents. `crew status` shows the stage and replies addressed to you.
If you see no progress, use the [quiet-team walkthrough](../04-team/README.md#4-see-a-file-message-and-its-wakeup)
to inspect the role that should act next before sending another request.

## 5. Let the human handle a request for more access

After the team finishes its feature, try a normal development need: downloading a
browser for future UI checks. In the SANDBOX shell:

```bash
crew ask "Ask the developer to run npx playwright install --with-deps chromium inside this sandbox, then npm run test:browser. Report passed, failed and skipped test counts. If a browser or system-package download is blocked, tell me the exact destination and why it is needed. Do not change host policy."
crew watch
```

Read the reported failure. A common blocked destination is `cdn.playwright.dev`.
Current SBX may report “Approval required” and suggest `sbx policy approval ls`.
Inspect that request in HOST; the sandbox-scoped allow below is another way to
record your decision for this workshop environment.
If that is the destination your agent reports and you decide to allow it, use HOST:

```bash
sbx policy allow network --sandbox sbx-ch-05 cdn.playwright.dev
```

This grants network access only for this sandbox and host. Return to the SANDBOX
shell (Ctrl-C leaves `crew watch`) and ask it to retry:

```bash
crew ask "The reported download host is now allowed. Ask the developer to retry installation and npm run test:browser, then report passed, failed and skipped test counts."
crew watch
```

`--with-deps` installs both Chromium and its Linux system libraries inside SBX.
Downloading the browser alone can leave required libraries missing. Ask the
developer to explain the check results: did Chromium start, and which tests
actually ran? A skipped test has not exercised the browser. If startup fails,
use its error message to identify what the sandbox still needs.

Redirects and system package repositories may introduce another destination;
Ubuntu mirrors such as `archive.ubuntu.com`, `security.ubuntu.com` or the ARM
mirror `ports.ubuntu.com` are common.
They are intentionally not pre-authorized in this exercise: when a request is
blocked, inspect it, ask the agent why it is needed, and from your host terminal,
allow only the specific hostname you accept. If installation and browser checks
already pass, continue to the kit definition below.

For a requirement you want every future worker to have, create
`factory/browser-access/spec.yaml` in your editor:

```yaml
schemaVersion: "2"
kind: mixin
name: workshop-browser-access
permissions:
  network:
    allow: [cdn.playwright.dev]
```

Use the destinations you actually needed. In HOST:

```bash
sbx kit validate factory/browser-access
```

Add this entry to `factory/sbxenv.yaml`'s existing `kits` list:

```yaml
  - source: ./browser-access
```

The kit records the decision for the next sandbox creation; it does not modify the
current one. On the next chapter's plan, look for this extra kit. Central AI
governance can apply additional organization rules; the presenter will demonstrate
that separately. [Network policy and tool governance](https://docs.docker.com/ai/sandboxes/security/)
are distinct controls.

## 6. See the factory's result

Open **<http://127.0.0.1:3102>**. Try assigning and resolving an incident. In HOST:

```bash
./scripts/beans show task-102
git -C sample-app log -1 --oneline
```

The task should contain the team's result note. The commit and changed code are
already in `sample-app/`. Compare the actual app behavior with the requirements;
a note saying “done” is not a substitute for trying the result.

When the team is finished, press Ctrl-C to leave `crew watch`, then type `exit`
to leave the SANDBOX shell. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=sbx-ch-05
```

Review the destroy plan. In SBX 0.50.0, removing this local environment also
revokes and removes its environment-defined MCP command. If you are recovering
an older installation or a manually registered server, inspect `sbx mcp ls` and
remove only a remaining `sbx-ch-05-beans` entry with `sbx mcp rm sbx-ch-05-beans`.
Keep the application as it is; the next task builds on your
assignment-and-resolution feature.

Next: [answer a human question through SSH](../06-human/README.md).
