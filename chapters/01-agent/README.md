# Chapter 1: Let an agent build inside a sandbox

Start by opening Claude Code in a sandbox. You'll ask it to run the sample app,
then give it a small coding task. First, try a few shell commands to see what it
can access: the source directory is shared with your laptop, so edits appear in
your editor, while commands and the database container run inside SBX. You'll
start manually, then make the setup repeatable in chapter 2.

## 1. Start Claude in your application

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

You are in the application directory, running as the sandbox user. Docker's server
is inside SBX; you do not need a host Docker engine.

We want to see which files are shared with the sandbox. In your host editor, open `host-only.txt` at the workshop root (the setup helper created
it) and replace its content with:

```text
A note outside the shared app.
```

Ask Claude:

> Create sandbox-message.txt in this project with the text "Hello from SBX".

Only `sample-app/` is mounted. Predict whether Claude can read the neighbouring
host file, then enter this in **Claude's input**:

```text
!cat ../host-only.txt
```

The agent wrote into the shared application directory. Reading the neighbouring
file should fail: you created that file outside the directory you mounted.

Back in your **HOST tab**:

```bash
# HOST
cat "./sample-app/sandbox-message.txt"
cat "./host-only.txt"
```

You can read both on the host. Edits inside the mounted app are bidirectional;
unrelated host directories have not been shared. The same applies to deletions
inside the mounted directory. Remove the demonstration file from Claude:

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
prepared. Later chapter launchers include this step for you.

Give Claude this prompt:

> Read README.md and get this application running inside the sandbox. Work out its
> dependencies and start any supporting services using the Docker engine inside SBX.
> Make the app available on 0.0.0.0:8080 and leave it running. Verify it responds.
> Do not change application code yet. Explain what you started and why.

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

![Incident Triage Board with eight incidents and severity and status filters](../images/incident-triage-board.png)

The browser runs on your host; the application and PostgreSQL run in SBX. The
published port connects them. If the page does not load, ask Claude to check the
server and `/healthz` before changing anything on the host.

## 4. Give the agent its first change

Our first demo task is a small change to the sample application: make its result
count follow the active filter. This is the warm-up exercise referred to in later
chapters. Try a severity or status filter, then give Claude this prompt:

> Read WORKSHOP-TASK.md and implement the active-filter result-count task. Run the
> relevant checks inside this sandbox, commit your change, and refresh the running
> app so I can try it. Use `Workshop learner` and `workshop@example.invalid` as the
> Git author if none is configured. Report what changed and how you checked it.

Try the filters again in your browser. Ask Claude to show the commit and explain
what it checked. You can inspect it yourself without leaving Claude:

```text
!git status --short
!git log -1 --oneline
```

The change and its Git commit are **already on your host** in `sample-app/`.
Open the changed file in your editor. Every later chapter will use this same working copy.

Type `/exit` in Claude to return to the host shell in your SANDBOX tab.
In HOST, remove the sandbox:

```bash
# HOST
sbx rm sbx-manual
```

The sandbox is gone, but you can still open the changed code and its Git commit
in `sample-app/`. In the next chapter, you'll give a new sandbox that same project
and a written task.

## Skip to the completed chapter

To skip to the completed exercise, exit and remove `sbx-manual` if it exists. Then run
this from the workshop repository:

```bash
./scripts/prepare-app.sh app-01-warmup-solution
```

This puts the supplied completed exercise in `sample-app/`. Your previous copy is
saved under `.local/saved-app.*`; the command prints its location. Continue with
chapter 2 using the same `sample-app/` path.

Next: [environment files and the host launcher](../02-launcher/README.md).
