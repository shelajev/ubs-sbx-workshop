# Build your software factory

An agent can make a change. A useful factory lets you hand over work with its
requirements, tools and boundaries already in place, inspect the result, and
step in when a decision belongs to you. You keep the project and choose how much
to delegate.

You'll build that arrangement one capability at a time. Each exercise leaves
something you can use in your own projects: a working sandbox, a task connection,
a kit, an environment definition, and eventually a team with a way to ask you
questions.

## Start here

Complete [laptop setup](00-setup.md), then follow the pages in order. Each page
ends with the next step. Use your own model access: Claude or a ChatGPT subscription
with Codex access works for the single-agent exercises. The team continuation
uses three Claude sessions sharing one account and model.

| Step | The question you answer | What you take away |
|---|---|---|
| [1. Give an agent a workspace](01-sandbox.md) | What can it touch, and where does code run? | A shared project, isolated execution, Docker inside SBX, and a published port |
| [2. Give it work](02-tasks.md) | Can it get requirements without me copying them each time? | A host backlog and an SBX-managed MCP connection |
| [3. Package what it needs](03-kits.md) | How does the next worker get my tools and review guidance? | A kit you made and a versioned policy with a review skill |
| [4. Recreate the environment](04-environment.md) | Can I start again without remembering the setup? | One environment file describing the workspace, kits, port and task connection |
| [5. Delegate between roles](05-team.md) | Who passes implementation to review? | A coordinator, developer and QA working through Herdr |
| [6. Give the team a feature](06-feature.md) | Can a task become a reviewed change with evidence? | A commit, checks and a result note in the host backlog |
| [7. Keep decisions human](07-human.md) | How does work pause for my answer and resume? | A running team you can join and direct |
| [8. Make it yours](08-next.md) | Where can this arrangement go next? | Your project, optional providers and a path toward host orchestration |

## A useful place to pause

At the end of step 4, you have made a kit and can recreate the environment you
assembled. The application gives us real work and containers to observe; you do
not need to spend time polishing it. Keep your configuration and continue into
the team exercises when you are ready. Each page has a cleanup step and a link
to the next exercise so you can leave the workshop with a working setup.

## Read the terminal labels

**HOST** means your laptop terminal at the workshop repository root. **SANDBOX
tab** starts at that same root, then connects to the assistant. **AGENT** means
paste the quoted prompt into the assistant. **SANDBOX shell** means a shell
inside SBX, not your laptop. Edit configuration with your normal editor.

Run one workshop sandbox at a time. Keep its connected terminal open while it is
working. Removing a sandbox stops its services and deletes its private files;
the mounted `sample-app/` and its Git history stay on your laptop. Database data
inside the sandbox is temporary; the app's setup can recreate its sample data.

If port 3102 is busy, check `sbx ls` in HOST for an earlier workshop sandbox and
finish its cleanup. Do not stop an unrelated service. For an interrupted creation,
use `sbx ls` to find the sandbox, then the current page's cleanup commands before
retrying. If you pause for the day, finish that cleanup and restart the page's
sandbox when you return.

**Next: [set up your laptop](00-setup.md).**
