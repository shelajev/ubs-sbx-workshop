# Build a software factory with Docker Sandboxes

Give an agent a place to work, a task it can retrieve, and tools and guidance it
can reuse. Then grow that arrangement into a team that implements, reviews and
asks you when a decision needs a human. Your project stays on your laptop; the
agents and their application processes run inside Docker Sandboxes.

**[Start the workshop →](workshop/README.md)**

The live session is 90 minutes: a 10-minute introduction followed by hands-on
work toward a repeatable single-agent environment. The guide continues into
teams, human decisions and applying the factory to your own project.

## What you will build

| Step | What you gain |
|---|---|
| [Run an agent manually](workshop/01-sandbox.md) | Inspect file isolation, run Docker inside SBX and reach the app through a published port. |
| [Give it work](workshop/02-tasks.md) | Paste a Beans task, then let the agent retrieve requirements through SBX-managed MCP. |
| [Package tools and guidance](workshop/03-kits.md) | Build a kit and use ACR to install a versioned policy and review skill. |
| [Recreate the environment](workshop/04-environment.md) | Record the workspace, kits, port and MCP connection in one environment file. |
| [Delegate between roles](workshop/05-team.md) | Give separate Claude sessions responsibility for coordination, implementation and review. |
| [Complete a feature](workshop/06-feature.md) | Follow a task through implementation and QA to a result note on your host. |
| [Keep decisions human](workshop/07-human.md) | Join a running team, answer its question and let it continue. |
| [Make it yours](workshop/08-next.md) | Bring another project, compare providers or explore host orchestration. |

The sample app tracks service incidents with a browser UI, an API and PostgreSQL.
It gives us real containers, code changes and checks to observe. Your takeaways
are the environment definition, reusable kits, task connection and working
practices you can adapt to your projects.

## Get started

Bring a Mac with Apple silicon or a Windows x64 machine with Git Bash, a browser,
and access to a model provider. A Claude subscription or a ChatGPT subscription
with Codex access works for the single-agent exercises. The team continuation
uses three Claude sessions sharing one account and Sonnet model.

[Set up your laptop](workshop/00-setup.md) for installation, authentication and
sample materials. The instructions target SBX v0.50.0. Docker containers run
inside SBX. Experienced Linux users can try the exercises on a self-supported
basis; the documented attendee setup covers macOS and Windows.

## About the tool choices

Docker Sandboxes provides environments, kits, network controls and the MCP
gateway. Beans, ACR, Herdr, Pi, Claude Code, Codex and the coding policy are the
author's choices for this example, not third-party endorsements by Docker.
[Docker Agent](https://docs.docker.com/ai/docker-agent/) is another option for
team coordination; the SBX boundaries apply whichever framework you choose.

## Repository map

| Directory | What is inside |
|---|---|
| `workshop/` | The sequential attendee guide, from setup to a factory |
| `chapters/` | Reference configurations, focused exercises and presenter extensions |
| `chapters/examples` and `chapters/kits` | Small examples and reusable kits |
| `chapters/support` | Launcher, sandbox helpers and role instructions |
| `scripts` | Materials and setup helpers, each introduced where it is needed |
| `backlog/seed` | Tasks used in the exercises |

## License

Workshop source is licensed under [Apache 2.0](LICENSE). The [Beans MCP adapter](mcp/beans/LICENSE)
and third-party components retain their own licenses and notices.
