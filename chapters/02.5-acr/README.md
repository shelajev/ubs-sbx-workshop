# Interlude: Give every worker the same tools and guidance

The next agent can see your code and task, but it also needs to know how your team
expects code to be written and reviewed. You'll give it a versioned coding policy
and a review skill, then use them to review the warm-up change.

We'll start with a small SBX kit so you can see how installation works. Then you'll
use a kit to install ACR, the tool that puts the policy and skill into your project.

## 1. Try a kit small enough to understand completely

A [mixin kit](https://docs.docker.com/ai/sandboxes/customize/kits/) adds capabilities
to an existing sandbox environment. In your editor, create
`factory/hello-kit/spec.yaml` with:

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
user. It creates a tiny command and makes it executable. `$HOME` is the sandbox
home. This shell code belongs to the kit; you are not running it on your laptop.

Validate the definition in HOST:

```bash
sbx kit validate factory/hello-kit
```

A valid result means SBX accepts its structure. Now try its actual effect in your
SANDBOX tab, starting at the workshop root:

```bash
sbx run claude "$(pwd)/sample-app" --name sbx-kit-first --skills off --cpus 4 --memory 8g --kit ./factory/hello-kit
```

`--kit` adds our installation step to the built-in Claude environment. In Claude:

```text
!workshop-hello
```

Expect `Hello from a kit`. The tool was installed during sandbox creation as part
of the kit setup. Exit with `/exit`. In HOST:

```bash
sbx rm sbx-kit-first
```

## 2. Use a kit to supply a real tool

The [ACR kit](https://github.com/shelajev/acr-sbx-kit) installs a package manager
for agent guidance. Before adding this Git-hosted kit, check the allowed publishers
in HOST:

```bash
sbx settings get kit.allowedSources
```

The kit needs `github.com/shelajev/` in that list. On a personal installation whose
current list is exactly `["docker.io/"]`, approve the workshop publisher with:

```bash
sbx settings set kit.allowedSources '["docker.io/","github.com/shelajev/"]'
```

This permits kits from that GitHub publisher. Preserve any other existing entries
when adding it; on a managed installation, ask your administrator to approve it.

Our [workshop policy package](https://github.com/shelajev/coding-policy/tree/workshop)
contains four coding rules and a `review-change` skill. Open that package and read
a rule: how would it help someone review the change you just made?

With ACR installed by the kit, you can choose which policy package to add to a
project. Updating that shared guidance doesn't require building a new sandbox image.

Add this section to `factory/sbxenv.yaml`:

```yaml
kits:
  - source: "git+https://github.com/shelajev/acr-sbx-kit.git#ref=f446b95bccabd879912db174c190ff09377b7c7b"
```

Unlike our local hello kit, this source is a Git repository pinned to a commit.
Every attendee receives the same kit revision. Keep `USE_ACR=0` for now: we want
to install the guidance together before automating it.

In HOST:

```bash
sbx env plan factory/sbxenv.yaml --env-arg name=sbx-ch-02-5
```

Find the ACR kit and the unchanged application mount. In your SANDBOX tab:

```bash
./scripts/launch-factory.sh sbx-ch-02-5
```

SBX may also ask you to approve the kit's credential bindings. In the linked kit's
`spec.yaml`, find `credentials`: it declares an **optional** service named
`acr-github`, for private repositories, higher API limits and publishing. Its
`inject` entries name the GitHub API, archive and upload hosts that may receive
that credential.

Read those service and destination names in the approval prompt. A binding
describes where a configured credential may be used; it does not mean every
download needs one. Our policy package is public, so you do not need to create a
GitHub token or run a secret-setup command for this exercise.

## 3. Install and inspect shared guidance

Ask Claude:

> Run acr --version. Then run the workshop's install-guidance helper. Show me
> agents.yaml, AGENTS.md and the installed review-change skill. Explain which file
> declares the package and which files tell an assistant how to work.

The helper installs a missing dependency, generates guidance, and checks that it
is current. Its main operations are:

```text
acr install github:shelajev/coding-policy@b85031eb0c8963b28b63eaa12efcbd34c850d32d --if-missing --non-interactive
acr realize --agent claude-code --agent codex
acr check --agent claude-code --agent codex
```

`install` adds the pinned policy without replacing an existing dependency choice.
`realize` writes the agent-facing guidance. Claude and Codex have generated formats;
Pi can read the shared `AGENTS.md` and skill too. `check` confirms that generated
instructions match the installed package. You should now have a review skill for
both Claude and Codex.

For a project without `agents.yaml`, the helper selects those two assistants and
sets `--freshness none` to keep the exercise on its chosen policy revision. An
existing project's package choices and configuration are kept. You can inspect
`chapters/support/bin/install-guidance` on the host to see the complete helper.

Open `sample-app/AGENTS.md` in your host editor. These are real files in the same
mounted project. Ask Claude:

> Read AGENTS.md and use the installed review-change skill to review the warm-up
> against ~/work/task.json. Explain one rule you checked and the code that supports
> your conclusion. Report a problem only if you find one. Do not edit application code.

Look at the rule and the code Claude cites. Can you follow how it reached its
conclusion? Ask it to explain if the connection is unclear; this is your chance
to see how the installed guidance affects a review.

## 4. Make guidance part of future preparation

In HOST, set `USE_ACR=1` in `factory/chapter.env`. Future launches install the policy
for a project without existing guidance; this project's installed guidance remains
in its working copy. The helper preserves a project's existing `AGENTS.md` instead
of silently replacing it.

Exit Claude. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=sbx-ch-02-5
```

We now have a repeatable environment and shared coding guidance. Next we will give
three Claude sessions different responsibilities using the same account and model.

Next: [build an all-Claude team](../04-team/README.md).

Optional extension: [try Pi and API-based models](../03-pi/README.md).
