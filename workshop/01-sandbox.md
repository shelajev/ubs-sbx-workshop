# 1. Give an agent a workspace

Start by opening Claude Code in a sandbox. You'll ask it to run the sample app,
then give it a small coding task in the next step. First, try a few shell commands to see what it
can access: the source directory is shared with your laptop, so edits appear in
your editor, while commands and the database container run inside SBX. You'll
start manually, then package the setup as you discover what it needs.

## 1. Start your assistant in the application

From the workshop repository root:

```bash
# SANDBOX tab — before connecting, from the workshop repository
sbx run claude "$(pwd)/sample-app" --name sbx-manual --skills off --cpus 4 --memory 8g
```

`claude` chooses the agent, `"$(pwd)/sample-app"` gives the full path to the sample
application to share, and `--name` gives the environment a reusable name. `--skills off` keeps the exercise
independent of host-installed skills; the final options give it four CPUs and 8 GB.
This creates the sandbox and opens Claude Code in your app directory. Keep this
terminal open while working.

**Using Codex?** Run the same line with `codex` in place of `claude`. Follow SBX's
browser authentication prompt. Wherever this guide says to ask Claude, give
Codex the same prompt. For a Claude `!` command, ask Codex to run the command
without the `!` inside the sandbox and report its output. Use `/quit` to leave
Codex where the guide uses Claude's `/exit`.

**Sign in:** if Claude asks you to authenticate, choose your subscription account
and follow the browser login. You can also type `/login` inside Claude. Existing
SBX credentials may mean you are already authenticated. See
[Claude authentication in SBX](https://docs.docker.com/ai/sandboxes/agents/claude-code/).

Before continuing, send Claude a short message such as “Hello, what directory
are we working in?” Wait for a response to confirm that your model access works.
If it reports an authentication error, ask the instructor for help before starting
the application exercise.

The built-in Claude configuration starts with permission prompts bypassed (the
“YOLO” mode for this exercise). SBX still enforces its own access boundaries.

## 2. Try shell commands without leaving Claude

Claude's [`!` shell mode](https://code.claude.com/docs/en/interactive-mode#shell-mode-with-prefix)
runs a command directly. Enter these one at a time **in Claude's input**, not in your
host shell:

```text
!pwd
!whoami
!docker version
```

Read the directory and user names. In the Docker output, find the Server section
and its operating system. Compare them with the directory and OS on your laptop.

We want to see which files are shared with the sandbox. In your host editor, open `host-only.txt` at the workshop root (the setup helper created
it) and replace its content with:

```text
A note outside the shared app.
```

Ask Claude:

```text
Create sandbox-message.txt in this project with the text "Hello from SBX".
```

Only `sample-app/` is mounted. Predict whether Claude can read the neighbouring
host file, then enter this in **Claude's input**:

```text
!cat ../host-only.txt
```

Keep that result in view while you check the same files from your laptop.

Back in your **HOST tab**:

```bash
# HOST
cat "./sample-app/sandbox-message.txt"
cat "./host-only.txt"
```

Compare the two terminals: identify which file the sandbox could read and which
file appeared on your laptop after the agent wrote it. Use the workspace path in
your `sbx run` command to explain the difference. Remove the demonstration file
from Claude, then check its absence in your host editor:

```text
!rm sandbox-message.txt
```

## 3. Ask the agent to run the app

On Windows, first run this in your **HOST Git Bash tab**, from the workshop root:

```bash
./scripts/prepare-workspace.sh sbx-manual
```

On Windows, this gives npm dependencies their own storage inside the sandbox,
where executable links work as npm expects. Your application source and Git
history stay shared with the host. The helper reports the dependency storage it
prepared. The team launcher introduced later includes this step for you.

Give Claude this prompt:

```text
Read README.md and get this application running inside the sandbox. Work out its
dependencies and start any supporting services using the Docker engine inside SBX.
Make the app available on 0.0.0.0:8080 and leave it running. Verify it responds.
I will publish that port as localhost:3102 on my laptop using the workshop
instructions. Do not change application code yet. Explain what you started and why.
```

Watch how Claude gets the app running: which dependencies does it install, and
which containers does it start? Ask it to explain a command you don't recognize.
If it asks where to run something, keep all app commands inside this sandbox.

When it reports the app is ready, ask to see its health response. In Claude:

```text
!curl -fsS http://127.0.0.1:8080/healthz
```

This requests the application's health endpoint from inside SBX. Look for an OK
status and the database being up. If it fails, ask Claude to finish startup before
continuing. This checks local health; the server must also listen on `0.0.0.0`, as
requested above, so the published port can reach it.

Then publish the web port from your HOST tab:

```bash
# HOST
sbx ports sbx-manual --publish 3102:8080
```

`3102:8080` connects port 3102 on your host to port 8080 in this sandbox.
Open **<http://127.0.0.1:3102>**. You should see the incident list and severity/status
filters. This is our starting application:

![Incident Triage Board with eight incidents and severity and status filters](../chapters/images/incident-triage-board.png)

The browser runs on your host; the application and PostgreSQL run in SBX. The
published port connects them. If the page does not load, ask Claude to check the
server and `/healthz` before changing anything on the host.

## 4. Inspect the access policy

The assistant can run development commands freely inside SBX. Access beyond that
environment is a separate decision. In HOST, inspect the policy for this sandbox:

```bash
sbx policy ls sbx-manual --wide
```

Find the network rules, their destinations and the sources of those rules. If
the assistant encountered a network error during startup, ask it for the exact
hostname and operation. Find the matching rule in the table, or identify that
no allowance covers the destination. If SBX reports an approval request, inspect
it in HOST:

```bash
sbx policy approval ls
```

Find the row for your sandbox and the reported hostname. Read what operation
needs approval; the list may also contain requests unrelated to the app.

If a needed download was blocked, decide whether that specific destination is
appropriate. This is the command pattern; replace `HOSTNAME` with the reported
hostname before running it:

```bash
sbx policy allow network --sandbox sbx-manual HOSTNAME
```

Then ask the assistant to retry. This scopes your allowance to this sandbox;
organization policy can still restrict access. If startup already worked, there
is nothing to approve. Later you can record recurring access needs in a kit.

You have now seen three boundaries: the app directory is shared, containers run
inside SBX, and host access to the app goes through a published port. Keep
`sbx-manual` and its agent session open for the next step.

Next: [give the agent work](02-tasks.md).
