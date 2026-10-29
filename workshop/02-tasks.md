# 2. Give the agent work

You have somewhere for the agent to work. Now give it requirements you can check.
First you'll paste a task, then let the same agent retrieve it through SBX's MCP
gateway. The task stays on your laptop throughout.

## 1. Prepare the workshop backlog

[Beans](https://github.com/hmans/beans) is a command-line task tracker backed by
Markdown files. Our backlog lives outside the shared app, under
`.local/chapters/beans/`. The agent cannot read that directory through its mount.

In HOST, run the supplied setup command once:

```bash
./chapters/00-setup/setup
```

It installs the pinned Beans CLI, installs the MCP adapter already downloaded
with the materials, and seeds four practice tasks. It marks this backlog as
disposable workshop data so the adapter can later append result notes. It leaves
your application alone. Wait for its ready message before continuing.

Read the first task in HOST:

```bash
./scripts/beans show task-101
```

The task changes the board's result count to describe the current filter. Copy
the printed requirements and acceptance criteria. Paste them into your running
agent in `sbx-manual`, followed by:

> Implement this task in the current project. Run the relevant checks inside
> this sandbox, commit the change, and refresh the app on 0.0.0.0:8080 so I can
> try it. Use `Workshop learner` and `workshop@example.invalid` as the Git author
> if none is configured. Report what changed and how you checked it.

Try the filters at <http://127.0.0.1:3102>. Open the changed file in your host
editor. To see the commit on your laptop, run in HOST:

```bash
git -C sample-app log -1 --oneline
```

The change is already in your working copy. No export or copy-back step is needed.
Read the agent's check results as well as the commit: a commit alone does not
show that the acceptance criteria passed.

## 2. Connect the task source through SBX

Copying requirements works, but each new worker needs them again.
[MCP](https://modelcontextprotocol.io/docs/getting-started/intro) lets an assistant
discover and call named tools. Here the connection is:

```text
Agent inside SBX → SBX MCP gateway → Beans adapter on your laptop → workshop backlog
```

Open `scripts/beans-mcp` in your editor. It selects this workshop's backlog and,
by default, exposes only listing and reading tasks. It does not expose arbitrary
shell commands or task deletion. The adapter runs with your host user's
permissions, so its code is part of what you trust. Keep it and the backlog
outside the writable application mount.

Register that host command with SBX. In **macOS HOST**:

```bash
sbx mcp add beans --command bash --args "$(pwd)/scripts/beans-mcp"
```

In **Windows HOST Git Bash**, use the Windows path to Bash and the script:

```bash
sbx mcp add beans --command "$(cygpath -w /bin/bash)" --args "$(pwd -W)/scripts/beans-mcp"
```

Run only the line for your platform. `add` records a command; it does not give it
to every sandbox. Read any host-command approval prompt and check the script path.
Now attach it to the sandbox you already have, in HOST:

```bash
sbx mcp load beans --sandbox sbx-manual
```

This is SBX configuration, not a per-agent MCP installation. Its gateway exposes
the server's tools to the assistant. Ask the agent:

> Use the Beans MCP tools to read the full body of task-101, including every
> acceptance criterion. Compare those requirements with the
> current code and the change we just made. Do not implement it twice or edit
> anything. Show me the tool result and explain whether the change meets it.

Compare the retrieved task with your earlier host output. You should see an
actual tool call, not an answer recalled from the pasted text. If the tools are
missing, inspect `sbx mcp ls` in HOST and the agent's `/mcp` view. Confirm you
loaded `beans` into `sbx-manual`; reconnect the agent with
`sbx run --name sbx-manual` if its tool list has not refreshed.

For a new task you can now say “Read task ID through Beans and implement it,”
instead of pasting the whole task. We used the completed task here so you can
compare the two ways of supplying exactly the same requirements.

## 3. Keep the work, replace the worker

Leave the agent with `/exit` in Claude or `/quit` in Codex. Remove this sandbox
in HOST:

```bash
sbx rm sbx-manual
```

Your app and commit remain. Keep the `beans` registration: you will attach
it to the next sandbox. If you stop the workshop here, remove the registration
with `sbx mcp rm beans`; repeat the registration step when you return.

If you need the completed warm-up, run the following in HOST after removing the
sandbox. It backs up your current app and installs the completed checkpoint:

```bash
./scripts/prepare-app.sh app-01-warmup-solution
```

Next: [package the tools and guidance](03-kits.md).
