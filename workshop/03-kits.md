# 3. Package the tools and guidance

A fresh worker can see your code, but it also needs your tools and expectations.
You'll make a small kit, then use a kit to install ACR and give the agent a
versioned coding policy and review skill.

## 1. Make your first kit

In HOST, create a directory for your configuration:

```bash
mkdir -p factory/hello-kit
```

Create `factory/hello-kit/spec.yaml` in your editor:

```yaml
schemaVersion: "2"
kind: mixin
name: workshop-hello
setup:
  install:
    - user: agent
      command: |
        mkdir -p "$HOME/.local/bin"
        printf '#!/bin/sh\necho "Hello from a kit"\n' > "$HOME/.local/bin/workshop-hello"
        chmod +x "$HOME/.local/bin/workshop-hello"
```

The `setup.install` block is code SBX runs **inside** the sandbox, as the `agent`
user. It creates a command and makes it executable. `$HOME` is the sandbox home;
this code belongs in the YAML file, not your host terminal.

Validate the definition in HOST:

```bash
sbx kit validate factory/hello-kit
```

A valid result means SBX accepts its structure. Now try its effect in your
SANDBOX tab, starting at the workshop root:

```bash
sbx run claude "$(pwd)/sample-app" --name sbx-kit-first --skills off --cpus 4 --memory 8g --kit ./factory/hello-kit
```

Use `codex` instead of `claude` if that is your chosen assistant. `--kit` adds
our installation step to the built-in agent environment. Ask the agent:

> Run ~/.local/bin/workshop-hello inside this sandbox and show its output.

Compare the output with the command in your kit file. Change the greeting in your
**host kit file** to something meaningful to you. The installed command does not change immediately:
kit setup runs when SBX creates the environment. Leave the agent and remove the
sandbox in HOST:

```bash
sbx rm sbx-kit-first
```

Run the same `sbx run` command again and ask for the greeting. Your new text
should appear. You have packaged a tool that each new worker receives without
an installation prompt. Exit the agent and remove `sbx-kit-first` again before
continuing.

A real project kit could install a pinned database migration CLI, a test runner,
or a small diagnostic command. Keep project-specific dependencies declared in
the project; use kits for capabilities you want to reuse across workers.

## 2. Add a tool for shared guidance

The [ACR kit](https://github.com/shelajev/acr-sbx-kit) installs a package manager
for agent guidance. Before using this Git-hosted kit, check allowed publishers
in HOST:

```bash
sbx settings get kit.allowedSources
```

The kit needs `github.com/shelajev/` in that list. On a personal installation whose
current list is exactly `["docker.io/"]`, add the workshop publisher with:

```bash
sbx settings set kit.allowedSources '["docker.io/","github.com/shelajev/"]'
```

Preserve any other entries you saw. On a managed installation, ask your
administrator to approve the publisher. This allows kits from that publisher;
it is separate from the sandbox's runtime network policy.

Create your next sandbox in the SANDBOX tab:

```bash
sbx run claude "$(pwd)/sample-app" --name sbx-guidance --skills off --cpus 4 --memory 8g --kit ./factory/hello-kit --kit 'git+https://github.com/shelajev/acr-sbx-kit.git#ref=f446b95bccabd879912db174c190ff09377b7c7b'
```

Again, choose `codex` in place of `claude` if needed. You are combining your own
kit with a Git-hosted kit pinned to a commit. Every creation uses that revision.

SBX may ask you to approve the kit's credential bindings. The kit declares an
**optional** service named `acr-github` for private repositories, higher API limits
and publishing. Read the service and destination names in the approval prompt.
Our package is public; you do not need a GitHub token for this exercise.

Attach the task connection in HOST so the new agent can review against the
original requirements:

```bash
sbx mcp load beans --sandbox sbx-guidance
```

On **Windows**, prepare dependency storage for this fresh sandbox before the
review skill runs application checks. In HOST Git Bash:

```bash
./scripts/prepare-workspace.sh sbx-guidance
```

Read the storage path the helper reports. Your source stays shared; npm will
install dependencies into this sandbox's storage, as in step 1.

## 3. Give the agent a review skill

Open the [workshop policy package](https://github.com/shelajev/coding-policy/tree/workshop)
and read one rule. It contains four coding rules and a `review-change` skill.
ACR installs that guidance into your project; the kit supplied the ACR executable.

Paste this prompt into the agent. These commands are instructions for it to run
inside SBX, not a list to run on your laptop:

> In the project directory, run acr --version. Install the public workshop policy
> with `env -u GH_TOKEN -u GITHUB_TOKEN acr install github:shelajev/coding-policy@b85031eb0c8963b28b63eaa12efcbd34c850d32d --if-missing --non-interactive`.
> If agents.yaml does not exist yet, add `--agent claude-code --agent codex --freshness none`
> to that install command. Preserve existing package choices and guidance; ask me
> before replacing any existing instructions. Then run
> `env -u GH_TOKEN -u GITHUB_TOKEN acr realize --agent claude-code --agent codex`
> and `env -u GH_TOKEN -u GITHUB_TOKEN acr check --agent claude-code --agent codex`.
> Show me agents.yaml, AGENTS.md and the installed review-change skill. Explain
> which file declares the package and which files tell an assistant how to work.

`install` selects the pinned package. `realize` writes agent-facing guidance.
`check` verifies that it matches the package. Clearing those two token variables
for these commands prevents a proxy placeholder being mistaken for a GitHub token
when downloading a public package. The sample project may already have
`agents.yaml`; that is why the prompt preserves existing configuration.

Open `sample-app/AGENTS.md` in your host editor. The policy and skill are files in
your shared project, so the next worker can use them too. Now ask the agent:

> Read task-101 through Beans MCP. Read AGENTS.md and use the installed review-change
> skill to review the warm-up against those requirements. Explain one rule you
> checked and the code supporting your conclusion. Report a problem only if you
> find one. Do not edit application code.

Look at the rule and the code it cites. Can you follow the reasoning? Installing
instructions makes them available; checking the agent's review shows whether it
actually used them.

## 4. Keep the reusable parts

Exit the agent. In HOST, remove the sandbox:

```bash
sbx rm sbx-guidance
```

The kit definition stays in `factory/`; the generated guidance stays in
`sample-app/`. Remove the manually registered task connection now, before moving
its definition into an environment file:

```bash
sbx mcp rm beans
```

You have built one kit and composed it with another. Next you will record those
choices together with the workspace, port and MCP connection.

Next: [recreate the environment](04-environment.md).
