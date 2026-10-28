# Workshop introduction

[Open the Google Slides deck](https://docs.google.com/presentation/d/1OLvxvAt-lByz606_1s37lgpuPQJ-SLFfcR18zrUBh9Q/edit)

An approximately eight-minute introduction before chapter00. The Google Slides deck includes these speaker notes.

## 1. Build a software factory

Opening, about 30 seconds.
The session title is “Docker's Agentic Platform: Sandboxes, MCP, and the Infrastructure of Autonomous Development.” Today we will assemble our own small software factory using Docker Sandboxes and supporting tools. Start with one coding agent, give it an environment where it can act, then grow it into a team you can give work.
We are using platform building blocks. The factory is the system we build together, not the name of a packaged Docker product.
Source: workshop README and the agreed session brief.

## 2. What can your agent touch?

About 40 seconds.
Quick show of hands: who has run a coding assistant with permission prompts disabled? What could it reach when you did that?
That is our starting point. Agents need to install dependencies, change files, run tests and start services. SBX gives us an execution boundary, with explicit decisions about shared files and network access. We will inspect that boundary rather than assume that “YOLO” makes it disappear.
An intentionally mounted directory is still shared. Isolation is not a promise that every resource you grant access to is protected from changes.
Source: chapters/01-agent/README.md.

## 3. A task becomes a reviewed change

About 45 seconds.
Our host starts the sandbox and supplies the project and task. Inside, a coordinator routes work, a developer implements it, and QA reviews the exact commit. The developer records the outcome back in our host backlog through MCP.
The agents inspect the repository and work out its dependencies and startup commands. Containers they need run inside SBX. The host launcher does not set up Node or PostgreSQL for them.
When the work is done we inspect the changes in the same mounted working copy. There is no extra host judge that accepts changes or closes tasks automatically.
Source: chapters/support/launch, chapters/support/bin/prepare, chapters/support/roles, chapters/05-mcp/PROMPT.md, chapters/07-factory/README.md.

## 4. A small app with real dependencies

About 35 seconds.
This is our shared sample application: a browser UI for tracking incidents, with an API and PostgreSQL. It gives agents real development work: install packages, run a database container, change behavior and test it. A published sandbox port lets us try the result in our browser.
We are building a coding factory. The agents implement features in this sample app; they are not an incident-response team operating production services. Later we replace this repository and task with our own.
Screenshot: chapters/images/incident-triage-board.png. Source: chapters/00-setup/README.md and chapters/01-agent/README.md.

## 5. Host and sandbox

About 55 seconds.
Keep this boundary in mind throughout. Beans, our file-based task tracker, stays on the host. A small script creates the sandbox from an SBX environment file and mounts the application working directory. The agents work inside that environment and run any containers there.
We do not need Docker Desktop or a host Docker engine. Every chapter uses the same mounted working copy, including its Git history. We replace sandbox environments as their capabilities grow; edits remain in the project on the host.
MCP supplies selected host operations; it does not mount the whole backlog directory. SSH lets us join the existing environment when a human decision is needed.
Source: chapters/01-agent, chapters/02-launcher, chapters/05-mcp, chapters/06-human.

## 6. Each step solves the next problem

About 60 seconds.
First, run an agent manually and inspect its access. Next, make that environment repeatable with sbxenv and a small launcher. Add a simple kit, then ACR to distribute coding guidance.
Now give three Claude Code sessions separate responsibilities using the same account and Sonnet model. Add Herdr, role instructions, file messages and explicit wakeups so those conversations can coordinate.
The team needs a live task and somewhere to report the result: connect a small host Beans server through SBX's MCP gateway. Inspect its limited tool surface and make a deliberate network-access decision. Join through SSH when the team needs a human answer.
Finally, keep this environment and bring another project and task. Every chapter leaves something working that the next one extends.
Source: chapters/README.md.

## 7. The environment grows with the job

About 55 seconds.
We try a tiny hello-world kit first so the mechanism is visible. Then ACR supplies the skill and coding policy. The policy package and the kit are different things: the kit installs the tool that retrieves and materializes the guidance.
We add Herdr as a mixin for session management. Pi is an optional extension for trying another assistant afterward. A template supplies the base image; kits layer the tools and settings onto it. We discuss templates without asking everyone to build an image.
These are technologies you can reuse with a different set of tools. Beans, ACR, Herdr, Pi and this workflow are my choices for the workshop, not third-party endorsements by Docker. Docker Agent is another team framework you could run inside SBX.
Sources: chapters/02.5-acr/README.md, chapters/03-pi/README.md, chapters/04-team/README.md, README.md. Docker Agent: https://docs.docker.com/ai/docker-agent/

## 8. One team, separate responsibilities

About 55 seconds.
SBX is also a place to try new assistants and models without installing their tool stacks on your laptop. Keep the project, task and guidance the same, then compare what each assistant actually does. Tools and processes live in the sandbox; edits to the mounted project persist.
Our default team uses Claude Code for coordinator, developer and QA, with Sonnet in all three roles. One Claude subscription supplies model access; the conversations have different responsibilities. File messages carry the durable handoff; Herdr manages sessions and delivers the wakeup. We do not equate a status indicator with completed work.
Everyone inspects the roles without configuring another provider. The mixed-provider variation is an optional presenter demonstration or take-home exercise. It is not needed for the MCP or human-decision chapters.
Source: chapters/04-team/team.tsv, chapters/04-team/team-mixed.tsv, chapters/03-pi/MIXED-MODELS.md.

## 9. Useful work needs explicit access

About 55 seconds.
Our MCP server has a small job: list tasks, read a task, append a result note to the disposable backlog. It has no separate OAuth setup. SBX connects the host stdio process to the sandbox. The server itself still runs with host-user permissions, so the narrow adapter is part of what we trust.
A legitimate download can be blocked. We read the denial and decide whether to allow that destination for this sandbox. A product ambiguity is a different issue: use SSH and the handoff protocol to supply a human answer without replacing the team.
At the end I can demonstrate organization governance, cloud execution and live mounts on my account. You do not need organization enrollment or cloud access for the attendee path.
Source: chapters/05-mcp/README.md, chapters/06-human/README.md, chapters/08-presenter/README.md.

## 10. What you need today

About 45 seconds.
Complete installation, downloads and the first successful Claude response beforehand when possible; we leave time for troubleshooting. Bring up the instructions and keep a host terminal ready. The hands-on path targets macOS on Apple silicon or Windows x64 with Git Bash and SBX v0.50.0. Application commands and containers run inside the sandbox.
Docker login and Claude login are separate. The main path uses your Claude subscription for every role, with no additional provider or API key. Check the account requirements in chapter00; we will verify a response before starting the team. We are not supplying model accounts.
Commands are labelled HOST, SANDBOX or the assistant name. Keep the launch terminal open and use another terminal to inspect the running environment. Stop the previous sandbox before reusing port3102. If a coding job takes longer, use the supplied sample-app checkpoint and continue learning the infrastructure.
Source: chapters/00-setup/README.md and chapters/README.md.

## 11. Let's build it

Transition to hands-on work, about 20 seconds.
Open the repository and start with chapter00. We will get SBX installed and model access working, then let one agent inspect and run the sample application. Keep the question from the start in view: what can it touch, and how do we choose that deliberately?
Your take-home artifact is the growing factory configuration. Chapter07 shows how to give it another repository and a task of your own.
Workshop: https://github.com/shelajev/wad-sbx-workshop
Setup: https://github.com/shelajev/wad-sbx-workshop/blob/main/chapters/00-setup/README.md
