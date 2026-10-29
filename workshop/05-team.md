# 5. Delegate between roles

So far, you've talked to Claude yourself. For a team to work, someone
needs to pass the developer's result to QA and bring questions back to you.
We'll give those responsibilities to separate agent sessions managed by
[Herdr](https://github.com/herdrdev/herdr). They exchange assignments and replies
through files, with a notification to tell the recipient it has a message.

You'll follow a small request through the team before asking it to change code.
All three roles use Claude Code and the same Sonnet model through your Claude
subscription. The roles share that account's usage limits.

Start here after removing `wad-repeatable` at the end of step 4. This continuation
uses Claude Code for all roles. If you used Codex for the single-agent exercises,
you can follow the role design and observe the team demonstration; running this
part yourself requires Claude access. Do not change provider credentials merely
to follow along.

## 1. Add session management

Open `chapters/kits/herdr/spec.yaml`. It pins a tool version,
allows its download hosts and installs the executable. It supplies a session
manager; the role instructions decide what those sessions should do.

Add this entry to the existing `kits` list in `factory/sbxenv.yaml`:

```yaml
  - source: ../chapters/kits/herdr
```

Keep the hello and ACR entries. Herdr adds session management to the environment
you already assembled. Set `agent: claude` in the environment file for this team.

[Docker Agent](https://docs.docker.com/ai/docker-agent/) is another way to define
agent teams, models, instructions and tools. It could also run inside SBX. This
workshop uses Herdr with coding assistants, following the author's own setup;
the environment and access-control concepts apply to either approach.

## 2. Give the assistants different responsibilities

Create `factory/team.tsv` in your editor. Each line has four fields, separated by spaces or tabs:

```text
coordinator	claude	anthropic	sonnet
developer	claude	anthropic	sonnet
qa	claude	anthropic	sonnet
```

Each row means **role, assistant, provider, model**. Keep these values for the
workshop. [`sonnet` is Claude Code's model alias](https://code.claude.com/docs/en/model-config#model-aliases);
you do not need to find a version-specific model ID. The coordinator routes work and questions; the
developer changes code; QA checks the proposed change against the requirements.

These are three separate conversations using the same Claude account and model.
The `anthropic` field identifies the provider; it does not require an API key or
switch Claude Code out of subscription authentication. The sessions share your
account's usage limits. If your account restricts Sonnet, select a working model
in a personal Claude session and use that same model in all three rows.

Different assistants and providers are an optional presenter extension. The
[mixed-provider demonstration](../chapters/03-pi/MIXED-MODELS.md) describes that extra setup.

Read the three short role briefs in `chapters/support/roles/`. Notice that only the
developer writes application code, and QA reviews a specific commit. A model name
alone does not establish those responsibilities. Keep Claude as the developer for
now; it already has access to the SBX MCP gateway.

## 3. Prepare the team workspace

Until now, SBX could open one assistant directly. The team needs role briefs,
message helpers and a startup instruction as well. We supply a small launcher
that copies those files into the sandbox and opens a shell there. It uses your
existing environment file. The `crew` commands it installs run **inside SBX**.

Create `factory/chapter.env` in your editor:

```text
TASK=wad-101
MODE=manual
USE_ACR=1
SESSION=shell
```

`TASK` selects the warm-up for this introduction. `MODE=manual` means you will
start the team yourself. `USE_ACR=1` prepares guidance for a project that needs
it, preserving this project's installed guidance. `SESSION=shell` opens a sandbox
shell so you can inspect and start the roles.

Create `factory/PROMPT.md`:

```markdown
Ask developer to read wad-101 through Beans MCP and inspect the existing change.
Ask QA to review the same commit against those requirements and the coding policy.
Do not change code. Send the human your combined assessment and evidence.
```

This is the request you can later submit; creating a team does not submit it.
Keep the `host_shell` default you verified in step 4, including the Windows Bash
path on Windows. This manual team launch still uses the MCP connection declared
in your environment file.
In the SANDBOX tab, at the workshop root:

```bash
./scripts/launch-factory.sh wad-ch-04
```

The launcher creates `wad-ch-04`, copies the workshop helpers and role briefs to
`~/work`, and opens a shell. Review the SBX plan and host MCP command approval.
It also prepares Windows dependency storage, so you no longer run that helper
separately. If guidance setup reports a missing review skill, run
`install-guidance` in this SANDBOX shell; it installs and checks the policy you
used in step 3.

## 4. Start the sessions and understand what that means

Before starting the team, check the shared account in a personal conversation
from the SANDBOX shell:

```bash
claude --model sonnet
```

If asked, use `/login` with your Claude subscription. Ask for a short response to
confirm access, then type `/exit`. Do not set an Anthropic API secret for this
subscription route. Back in the same SANDBOX shell:

```bash
start-team
```

This supplied helper reads `team.tsv` and performs three Herdr operations:

| Operation | Effect |
|---|---|
| `herdr server` | Starts the session service inside this sandbox. |
| `herdr agent start` | Launches the assistant selected for a role. |
| `herdr agent prompt` | Gives that assistant its role brief. |

Open `chapters/support/bin/start-team` on the host if you want the complete working
example. Workspace IDs, environment variables and terminal creation are plumbing
around these three operations.

Inside SBX, inspect the sessions once:

```bash
herdr agent list
```

`list` should show coordinator, developer and QA. To inspect one assistant's
terminal, run `crew logs coordinator` in the same sandbox shell. This is useful
for understanding its introduction or diagnosing a login problem.
The workshop's `crew` helper lets you send messages to the team and read its
replies. Use `crew help` to see its commands.

## 5. See a file message and its wakeup

We will ask a small question without changing the app. In the SANDBOX shell:

```bash
crew ask "Ask developer to list the test scripts in package.json, ask QA to check that list, and send me the combined answer. Do not change the application."
```

`ask` waits for role acknowledgments, stores your message, and notifies the
coordinator. Now follow the replies in the same SANDBOX shell:

```bash
crew watch
```

`watch` displays readable messages as the team passes work between
roles. Expect a developer answer, a QA check and a combined reply to you. Press
Ctrl-C to return to the shell; the agents continue running.

If no reply appears, the assistant might still be working, waiting for login or
having trouble reaching its model. Press Ctrl-C to leave the message view, then
inspect the coordinator in the SANDBOX shell:

```bash
crew status
crew logs coordinator
```

`status` shows the task stage and messages addressed to you. `logs` shows the
assistant's current terminal output; it does not open an interactive conversation
with that assistant. If the coordinator has handed work to the developer, inspect
`crew logs developer` next; use `crew logs qa` for the review. Look for an active
request, an authentication prompt or an API error. If you need help, show the
instructor that output and your role table. Once the role can continue, use
`crew watch` to follow replies. A quiet message view alone does not tell us whether
an assistant has stopped.

The helper saves the message and notifies the coordinator using these two
operations. This block is a **reference** to help you read the helper; you have
already sent the request with `crew ask`:

```text
handoff send --to coordinator --from human --kind question --body "Your request"
crew-notify coordinator
```

Saving the request in a file lets the recipient read it later, but the assistant
might keep working without noticing that file. The notification gives it a turn
to read the message. We can then follow the handoff through the stored messages,
even when an assistant's busy/idle status is unreliable.

To see how a conversation is stored, ask the coordinator through the SANDBOX shell:

```bash
crew ask "Show one message file from ~/work/factory/messages/ and explain its sender, recipient and body. Do not change any files."
crew watch
```

Look for the same request and roles you just followed in the message view.
Press Ctrl-C when you have read the answer.

Agents use `crew send ROLE "message"` to save and deliver their replies together, just as
you use `crew ask`. We'll keep using these helpers for the rest of the workshop.

## 6. Finish this environment

You have started the sessions yourself and followed their messages. In the next
step, the launcher will start them automatically; submitting work remains a
separate action so you can inspect the team before giving it a task.

Type `exit` in the SANDBOX shell. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=wad-ch-04
```

Removing this environment also removes its `wad-ch-04-beans` MCP registration
on SBX 0.50.0. The source and configuration remain. We have a team that can discuss
work; next we give it a feature and enable one extra backlog operation: a result note.

Next: [give the team a feature](06-feature.md).
