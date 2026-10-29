# 7. Keep decisions human

The next task is to let users reopen a resolved incident. There's a product
decision to make first: should reopening clear the current resolution note or keep
it? You'll tell the team to ask you before it starts implementing that behavior.

When the question arrives, you'll join the running sandbox through SSH and send
your answer to the coordinator. The same team should then continue using your
decision. SSH gets you into the environment; the role instructions and messages
handle the conversation.

## 1. Give it a task with a real choice

Keep using `sample-app/`, including the assignment-and-resolution feature from
[step 6](06-feature.md). If you skipped implementation, use the checkpoint at
the end of that page before continuing. Finish that environment's cleanup first.

Read the next task in HOST:

```bash
./scripts/beans show wad-103
```

Find the requirement about reopening an incident. The team needs your answer
about the current resolution note before it can decide what behavior to implement.
Change the `TASK` line to `TASK=wad-103` in `factory/chapter.env`. Keep `MODE=mcp`,
`USE_ACR=1` and `SESSION=shell`. Keep the existing prompt and append this paragraph
to `factory/PROMPT.md`:

```markdown
Before implementing reopening, ask the human whether to clear the current
resolutionNote or retain it. Explain both interpretations and wait for an answer.
The coordinator records the human's choice for the developer and QA, then resumes
implementation and review. Preserve the historical resolution entry either way.
```

The existing prompt reads the task ID supplied by `TASK`, so it now asks for
`wad-103`. The added paragraph tells the team which decision belongs to you.
The environment, tools and team are the same.
In the SANDBOX tab, at the workshop root:

```bash
./scripts/launch-factory.sh wad-ch-06
```

Then in its sandbox shell:

```bash
crew submit
crew watch
```

`crew submit` sends the reopening task to the coordinator; `crew watch` follows
the discussion. Wait for a message asking whether to clear or retain the note
before connecting through SSH to answer it. If you added the browser-access kit in
the optional browser exercise, it is also installed in this newly created environment.

## 2. Reach the existing team through SSH

Leave the SANDBOX session open. In HOST, configure SSH:

```bash
sbx setup ssh
```

On macOS or Linux, connect with:

```bash
ssh wad-ch-06.sbx
```

On Windows, use the native Windows OpenSSH client from Git Bash:

```bash
/c/Windows/System32/OpenSSH/ssh.exe wad-ch-06.sbx
```

Use this full path so you get the Windows client configured during laptop setup.
It can read the Windows paths that SBX puts in the SSH configuration.

`setup ssh` configures the host's SSH integration. The connection command opens an
interactive shell in the running sandbox, alongside the existing agent sessions.
You do not need to find an IP address or install an SSH server in the app.
See [SBX integrations](https://docs.docker.com/ai/sandboxes/integrations/).

You do not need a third terminal: your HOST tab becomes the SSH session
temporarily, while the original SANDBOX tab remains connected. After answering
the team, `exit` closes SSH and restores the HOST tab. Inside the SSH session:

```bash
crew status
```

You should see the same task and product question as in your SANDBOX tab.
The message tools in this SSH shell read the existing team's files.

## 3. Give a product answer

For this walkthrough, clear the current note while retaining history:

```bash
crew reply "Clear resolutionNote when reopening. An open incident has no current resolution, but retain the old note in history. Continue implementation and review using that decision."
crew watch
```

`crew reply` sends your answer to the coordinator, which records the decision
for this task and passes it to the developer and QA. Both should now work from
the same interpretation of reopening.

Look for the choice being used in implementation and review. When the team reports
that the app is ready, open <http://127.0.0.1:3102> and reopen a resolved incident.
Does the current note clear while history retains the old resolution?

Ctrl-C leaves the message view. Type `exit` to close SSH and restore your HOST tab.
Your original SANDBOX tab remains connected; SSH was another doorway into the
same running team. The code changes are already in `sample-app/`.

## 4. Finish this environment

Once the team is finished, leave its message view and type `exit` in the original
SANDBOX tab. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=wad-ch-06
```

Review the destroy plan. SBX 0.50.0 also removes the local MCP command registered
by this environment. If an older or manually registered `wad-ch-06-beans` entry
remains in `sbx mcp ls`, remove that entry with `sbx mcp rm wad-ch-06-beans`.

For the next project, you can answer product questions through the same team
conversation. If an agent needs more network access, use the host policy controls
from step 1.

Next: [make the factory yours](08-next.md).
