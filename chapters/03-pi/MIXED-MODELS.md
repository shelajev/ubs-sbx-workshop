# Demonstration: different models for different responsibilities

Repeat the team's request-and-review exchange with a different combination:
Pi using Gemini to coordinate, Claude doing development and Codex reviewing.
Keep the same roles and messages so the audience can follow a familiar task
across the new assistants and providers.

This optional presenter demonstration uses **`gemini-3.8-flash`** for its Gemini
role. The attendee path uses Claude Code and Sonnet for all three roles; it does
not require any of these additional accounts or model choices.

## What changes, and what does not?

Compare `chapters/04-team/team.tsv` with `team-mixed.tsv`:

| Role | Assistant | Provider |
|---|---|---|
| Coordinator | Pi | Google |
| Developer | Claude Code | Anthropic |
| QA | Codex | OpenAI |

The fourth column selects each model. The role briefs and file-message protocol
remain the same. A different model does not require a new way to send a review request.

Open `chapters/kits/multi-provider/spec.yaml`. This is a **sandbox kit**, which
selects a starting image and declares provider credential bindings, rather than
just a mixin. Read it in these pieces:

1. `sandbox` selects the base image and entrypoint.
2. `credentials` describes each provider's request hosts and credential injection.
   API-key and OAuth routes differ: a subscription login must not be forced into
   API-key mode by setting a placeholder variable indiscriminately.
3. `permissions.network` allows the endpoints these assistants need.
4. `setup.install` installs Codex and prepares the assistant configuration appropriate
   to the selected authentication mode. Pi and Herdr remain separate mixins.

The real credentials live in SBX's host credential store, not the role table.
Each assistant can use the account route that fits it:

- **Claude Code:** use the Claude subscription login from chapter 1, or register
  an Anthropic API key with `sbx secret set anthropic`.
- **Codex:** use a ChatGPT subscription by running `sbx secret set openai --oauth`,
  or register an OpenAI API key with `sbx secret set openai`.
- **Pi on Gemini:** register a Gemini API key with `sbx secret set google`.

For an API-key-based presenter setup, run:

```bash
sbx secret set anthropic
sbx secret set openai
sbx secret set google
```

Enter the corresponding keys at the secure prompts. You do not need both a
subscription and an API key for Claude Code or Codex. Pi itself uses provider APIs,
so its normal Anthropic route needs an Anthropic API key; this mixed route instead
uses the Gemini key. The kit selects the Gemini credential service `google`;
the kit's `google_secret_service` argument can also select a legacy service named
`gemini`.

## Presenter walkthrough

Finish and remove the normal chapter sandbox first. This demonstration uses the
same mounted `sample-app/` and port 3102, so it runs on its own. In the SANDBOX tab:

```bash
./chapters/04-team/launch-mixed sbx-mixed
```

The helper uses `sbx create` with the custom sandbox kit and the ACR, Pi and Herdr
mixins, then opens the same workshop shell. Try:

```bash
crew ask "Ask developer to list the project's test commands, ask QA to check the answer, then report to me. Do not change application code."
crew watch
```

Use `crew logs coordinator` and `crew logs qa` to show which assistants answered.
Follow the replies back to the coordinator's combined answer so the audience can
see the exchange across providers. If a provider is out of quota, explain which
call failed and continue with the working team configuration.

Exit the shell when finished. In HOST:

```bash
sbx rm sbx-mixed
```

## Carry the mixed team into the MCP exercise

The custom kit also configures Claude's MCP client. To demonstrate the chapter 5
connection with this same team, use the supplied MCP reference configuration:

```bash
# SANDBOX tab — before connecting, at the workshop root
./chapters/04-team/launch-mixed sbx-mixed-mcp 05-mcp
```

The second argument chooses the chapter's completed settings and prompt. This
creates one new sandbox on the same `sample-app/`. Inside its shell, try a small
read-only request before a full feature:

```bash
crew ask "Ask developer to read task-101 through MCP and list package.json test scripts. Ask QA to check the list independently. Report the combined answer; do not change project files or write a task note."
crew watch
```

The result should include a real host task read by Claude, an independent Codex
check, and Pi's combined answer. The gateway-capable role remains Claude; the other
roles communicate with it through the same file-message helpers.

The presenter helper uses `sbx create KIT --kit MIXIN` and
attaches the host server with `--static-mcp`. These are different ways to compose
the same SBX building blocks. The reference helper uses the ACR, Pi and Herdr kits;
it does not read extra kits from your evolving `factory/` configuration.

Exit the shell, then in HOST:

```bash
sbx rm sbx-mixed-mcp
sbx mcp rm sbx-mixed-mcp-beans
```

The main attendee path continues with `factory/` and the built-in Claude environment.
