# Optional chapter 3: Give the factory a choice of assistant

The main workshop continues from shared guidance to
[chapter 4's all-Claude team](../04-team/README.md). This extension is for your own
experiments or a presenter demonstration. It is not a prerequisite for teams,
MCP or human intervention, and it requires separate API access for model requests.

You have a project and shared coding guidance ready for an agent. Let's try a
different assistant with that same setup. SBX gives you somewhere to install its
tools and try its models while controlling what it can access. The assistant's
installation stays inside the sandbox; it can still work on your real project.

We will install [Pi](https://pi.dev), an open-source coding assistant, and
use it to explore the same application and coding policy. You can run the model
conversation with Anthropic API access, or follow a partner's or the presenter's
conversation if you have only a Claude subscription.

Keeping the same project, task and guidance gives you a familiar change to discuss
with Pi. You can compare how it finds the relevant code and explains its review.

This introduces three separate choices: the assistant program, its model provider,
and the particular model. The main team uses Claude Code with one model for every role.

## 1. Understand what the Pi kit supplies

Open `chapters/kits/pi/spec.yaml`. Read these parts before using the whole file:

| Part | Purpose |
|---|---|
| `args.pi_version` | Select a repeatable version of the assistant. |
| `permissions.network.allow` | Permit the download hosts and model API endpoints it needs. |
| `setup.install` | Download the right build, verify its checksum and install it as the sandbox user. |

The longer architecture and checksum code makes that installation step
portable. You do not need to rewrite it. The important choices are **what to
install**, **where it may connect**, and **which user runs setup**.

A [template](https://docs.docker.com/ai/sandboxes/customize/templates/) supplies the
starting filesystem. A [kit](https://docs.docker.com/ai/sandboxes/customize/kits/)
adds capabilities and access requirements. We can combine kits on the existing
Claude template without building our own image.

## 2. Compose it with the guidance kit

Add this entry under the existing `kits` list in `factory/sbxenv.yaml`:

```yaml
  - source: ../chapters/kits/pi
```

Keep the ACR entry. Relative kit paths are resolved from the environment file,
so this points from `factory/` to the supplied Pi kit.

Set `SESSION=shell` in `factory/chapter.env`. We want a shell where we can start Pi
ourselves; `MODE=manual` and `USE_ACR=1` stay as they are.

In HOST:

```bash
sbx env plan factory/sbxenv.yaml --env-arg name=sbx-ch-03
```

Check that both kits and the same application mount appear. In the SANDBOX tab:

```bash
./scripts/launch-factory.sh sbx-ch-03
```

The launcher opens an interactive shell inside SBX, already in the mounted project.
It keeps this sandbox session open while you work. Try:

```bash
# SANDBOX shell
pi --version
```

You should see the version selected by the kit. Now we will use the actual assistant.

## 3. Connect an assistant to a provider

Claude subscription login and Pi's API access are separate. Choose your route
before starting a model conversation:

- **Have Anthropic API access?** Follow this section, then try the prompts in
  sections 4 and 5 yourself.
- **Have only a Claude subscription?** You have already installed Pi and checked
  its version. Follow sections 4 and 5 with a partner or the presenter: look for
  the code and policy rule Pi uses to support its answer. You do not need to set a
  secret or send Pi a model request on your machine. Then go to
  [finish this chapter](#6-finish-this-chapter). In chapter 4, use the all-Claude
  role table to build and run your own team.

For the API route, register your Anthropic key from HOST:

```bash
sbx secret set anthropic
```

Enter the real key at the secure prompt, not in a command or file. Skip this if the
service is already configured. SBX keeps the real credential on the host and injects
it into permitted model requests. Allowing a network host and providing a credential
solve different problems: the request needs both permission and authentication.

In the SANDBOX shell, start Pi with the Anthropic provider and the model used in
this example:

```bash
pi --provider anthropic --model claude-sonnet-5
```

Our kit's Pi wrapper supplies a placeholder API value when this route needs one;
the SBX proxy supplies the real credential. Open the kit's wrapper section if you
want to see that mechanism. No real key is written into the project.

Pi knows how to make a model request, while SBX controls where it can connect and
supplies the credential for that request.

Accept the workspace trust prompt if shown. In Pi, type `/model` to see available
models and choose one your account can use for this optional experiment.

## 4. Use the same guidance with a different assistant

Ask Pi:

> Look around this project and explain what our sample application does. Where is
> the filter-result count calculated? Do not change anything yet.

Then:

> Read AGENTS.md and use the installed review-change skill. Review the warm-up
> described in ~/work/task.json. Explain one rule you checked against actual code.
> Do not modify the application.

Follow up in your own words. Does its answer help you understand the change? You
have changed the assistant while keeping the application, task and guidance.

## 5. Use this environment to try another model

Keep Pi open and type `/model`. Look for another **Anthropic model your API
account can use**, so you can try changing models with the provider connection
already configured. Select it if you have access. If only one model is available
to you, keep it and continue.

To try a different provider, you need its credential binding and its API endpoint
must be allowed by SBX. Selecting a model name does not configure either one. The
supplied Pi kit already allows Google's `generativelanguage.googleapis.com` endpoint,
but the normal attendee environment does not bind a Google credential. You can see
that additional credential configuration in the
[mixed-provider demonstration](MIXED-MODELS.md), which runs Pi on Google with
**`gemini-3.8-flash`**, alongside Claude and Codex.

Before repeating the request, start a fresh conversation with `/new` so the first
answer doesn't influence the next one. Then compare the reviews: does each answer
find the relevant code, apply the shared rule and explain its conclusion? Judge
those details as well as how easy the answer is to read.

This is a reusable way to evaluate tools: put an assistant's installation and
network requirements in a kit, create its sandbox, give it a bounded task, and
inspect what it actually does. When you remove the sandbox, its installed tools
and running processes go away. The mounted project and its edits remain, so use
a separate checkout when an experiment should not affect work you want to keep.

## 6. Finish this chapter

If you opened Pi, type `/quit` to return to the sandbox shell. The
[Pi command reference](https://github.com/earendil-works/pi/blob/v0.85.1/packages/coding-agent/README.md#commands)
explains its other interactive controls. If you followed someone else's model
conversation, your SANDBOX tab is already at its shell prompt.

Type `exit` in the SANDBOX shell to return to the host. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=sbx-ch-03
```

Next: [give the assistants roles and communication](../04-team/README.md).
