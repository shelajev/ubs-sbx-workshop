# Workshop introduction

[Open the Google Slides deck](https://docs.google.com/presentation/d/1OLvxvAt-lByz606_1s37lgpuPQJ-SLFfcR18zrUBh9Q/edit)

A ten-minute introduction followed by eighty minutes of hands-on work. These are
the speaker notes for the ten slides in the Google Slides deck.

## 1. Safety, speed, and governance in the agentic era

About 40 seconds.
A personal software factory gives you a way to hand off more of the work between an idea and a change you can inspect. You choose the task and the access it needs. The agent gets an environment where it can work, and you keep the reusable setup for your next project.
Today we have 90 minutes: ten minutes of introduction, then hands-on work. We will build the foundation with one agent. The guide continues into teams and human decisions when you are ready.
Source: https://github.com/shelajev/ubs-sbx-workshop/blob/main/workshop/README.md

## 2. What can your agent touch?

About 60 seconds, including a show of hands.
Who has run a coding assistant with permission prompts disabled? What could it reach?
Agents need to install dependencies, edit files and start services. Docker Sandboxes gives them an execution environment with explicit shared directories and network controls. An intentionally mounted directory remains shared, so changes there affect your project.
In the first exercise, run the commands and compare what you can see from the host and sandbox. Inspect the full policy table and find the relevant rows. We will observe the boundary ourselves.
Source: https://github.com/shelajev/ubs-sbx-workshop/blob/main/workshop/01-sandbox.md

## 3. A task becomes a reviewed change

About 75 seconds. Trace Task, Factory, Result, then Human.
A task goes in. The system runs the project, changes it and checks the result. You get a change you can inspect, together with the checks and any unresolved questions. You still decide whether to accept it.
The human connection stays in the picture. You supply the goal, choose access and answer questions the task did not settle. As the factory grows, separate roles can pass work between themselves without you carrying every intermediate message.
We will build toward that picture in stages. Today the useful foundation is one working agent, a task connection and an environment you can recreate. The team is a continuation in the guide.
The payoff is having your own setup to use again: different tasks and projects can use the same tools, guidance and boundaries.
Source: https://github.com/shelajev/ubs-sbx-workshop/blob/main/README.md

## 4. A project for our agents to work on

About 40 seconds.
This is our sample project: an incident list with a browser interface, an API and PostgreSQL. It gives the agent dependencies to install and a database container to run inside SBX. A published port lets us reach it from the host browser.
Our first task is a small change to the filtered incident count. We will paste the task, then connect the same backlog through MCP so the agent can retrieve requirements itself.
The app gives us something concrete to observe. Later exercises add assignment and resolution, then a reopening task with a human product decision.
Sources: https://github.com/shelajev/ubs-sbx-workshop/blob/main/workshop/01-sandbox.md and workshop/02-tasks.md

## 5. Host and sandbox

About 75 seconds.
Your project files and the Beans task backlog stay on your laptop. Start SBX manually from a host terminal and mount the sample project. The agent installs dependencies and runs its application processes and Docker containers inside the sandbox.
Edits to the mounted project are visible on the host. Replacing a sandbox does not discard those source edits. Later, sbxenv records the environment so we can recreate it with the same kits, workspace, port and MCP connection.
For the first task, paste its requirements. Then configure MCP through SBX on the host. The agent can ask the gateway to retrieve a task from Beans without mounting the backlog directory into the sandbox. The adapter runs on the host and is part of the access you grant.
We will compare files across the boundary and open the app through a published port. Keep host and sandbox terminals clearly separated as you follow the guide.
Sources: https://github.com/shelajev/ubs-sbx-workshop/blob/main/workshop/01-sandbox.md, workshop/02-tasks.md and workshop/04-environment.md

## 6. Today and the continuation

About 90 seconds.
The introduction takes ten minutes, leaving eighty for hands-on work. The target is a repeatable single-agent environment. Follow the guide at your pace; it remains available after the session.
First, run SBX manually. Observe the shared-file boundary, inspect network policy, run Docker inside SBX and reach the app from the host.
Second, paste a Beans task and observe the change in your host working directory. Connect Beans through SBX-managed MCP, then ask the agent to read that same task by ID.
Third, build a small kit yourself. Change it and recreate the sandbox to see when setup takes effect. Add ACR so the environment can install versioned coding guidance and a review skill.
Fourth, capture the workspace, kits, port and MCP connection in sbxenv and recreate the environment. These files and kits are the main takeaways from today.
After today, continue into Herdr, separate Claude roles, a feature and a human decision. Other model providers and host orchestration are later options.
Source: https://github.com/shelajev/ubs-sbx-workshop/blob/main/workshop/README.md

## 7. Reusable tools and guidance

About 85 seconds.
A kit lets you describe tools and setup that a fresh sandbox should receive. We will make a tiny command, install it through a kit, then edit the kit and recreate the sandbox. That makes the lifecycle visible before adding more tools.
The ACR kit installs the package manager. You then ask the agent to use ACR to install a pinned coding-policy package and make its rules and review skill available. The kit and the guidance package have separate jobs.
Finally, sbxenv records how the environment fits together: which project to mount, which kits to apply, which port to publish and which MCP connection to load. Recreate the environment and check that those choices come back.
You leave with configuration you can adapt for your own tools and project. In the continuation, a Herdr kit adds the sessions used by the team.
Beans, ACR and Herdr are my tool choices for this example. Docker Sandboxes supplies the environment, kits, network controls and MCP gateway.
Sources: https://github.com/shelajev/ubs-sbx-workshop/blob/main/workshop/03-kits.md and workshop/04-environment.md

## 8. Later: a team of Claude sessions

About 45 seconds. Preview of the continuation.
Once the environment is repeatable, give separate sessions different responsibilities. The coordinator routes work and asks for decisions, the developer implements and checks the change, and QA reviews the exact commit against the task.
The supplied team uses three Claude Code sessions, all with Sonnet and the same Claude subscription. Herdr manages their sessions and message delivery. You can inspect the messages and follow a task through review.
For today's single-agent work, Claude or Codex access is suitable. The documented team continuation uses Claude. Trying another provider or host orchestration comes later.
Source: https://github.com/shelajev/ubs-sbx-workshop/blob/main/workshop/05-team.md

## 9. Access to tools, networks and people

About 60 seconds.
MCP lets the agent use selected host operations. We start with read-only Beans tools to list and retrieve tasks. Later, the feature exercise enables appending a result note. SBX manages the connection, while the adapter runs on the host.
Network access is a separate decision. Inspect the full policy table and find the relevant row. When a request is blocked, identify the destination and decide whether the task needs that access.
A product ambiguity needs a human answer rather than a network rule. In the continuation, join the existing team through SSH, read its question and send a decision. The team records the answer and continues.
Sources: https://github.com/shelajev/ubs-sbx-workshop/blob/main/workshop/02-tasks.md, workshop/06-feature.md and workshop/07-human.md

## 10. Let’s build it

About 30 seconds. Transition to hands-on work.
Open the repository and click Start the workshop. Begin with setup if needed, then Run an agent manually. You need your own model-provider access; a Claude subscription or a ChatGPT subscription with Codex access works for the single-agent exercises.
Keep the guide and a host terminal open. We will work toward reusable kits and an environment we can recreate. Continue building the factory afterward with the linked team and human-decision exercises.
Workshop: https://github.com/shelajev/ubs-sbx-workshop
Guide: https://github.com/shelajev/ubs-sbx-workshop/blob/main/workshop/README.md
Setup: https://github.com/shelajev/ubs-sbx-workshop/blob/main/workshop/00-setup.md
