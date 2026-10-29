# Presenter extensions: explore the boundaries you can change

These notes give the presenter a choice of demonstrations to explore with the
audience: sandbox inspection, live mounts, cloud execution and organization
governance. Choose the ones you have prepared and have access to. For each,
start with a question about the factory the attendees built and show what changes.
Attendees can follow along on screen; completing the workshop doesn't depend on
setting these up on their laptops.

The seeded backlog also contains `wad-104`, an optional take-home or presenter
exercise that adds a service column and filter to the incident table. It is not
part of the chapter 1–7 walkthrough.

## Prepare the main workshop

For the 90-minute session, use the [attendee path](../../workshop/README.md) and
its [delivery notes](../../workshop/PRESENTER.md): 10 minutes of introduction,
then hands-on work toward a kit and repeatable single-agent environment.
The team and human-decision exercises are continuation material. The complete
factory rehearsal below is useful for preparing those demonstrations.

Use the same SBX release as attendees: **v0.50.0**. Recheck the release and
installation links before each delivery. Keep the existing v2 kits for this lab;
moving to v3 kits is a separate exercise.

The two-hour attendee path keeps the team and uses Claude Code with the same
Sonnet model for coordinator, developer and QA. Skip chapter 3; additional
providers and chapter 7's own-project exercise are extensions.

| Minutes | Activity |
|---|---|
| 15 | Setup recovery and an authenticated Claude response |
| 20 | One agent, the running app, isolation and ports |
| 15 | Repeatable environment and the task backlog |
| 15 | A small kit and shared guidance |
| 15 | Three Claude roles and a small handoff |
| 20 | Scoped MCP task access and the tested change |
| 10 | Human decision and inspection of the result |
| 10 | Catch-up, questions and recovery buffer |

Before delivery, rehearse from a fresh download in a separate workshop clone:

1. Complete chapter 0 and get a successful Claude response in chapter 1. Verify
   the database, app and browser access through the published port.
2. Launch the chapter 4 reference after the guidance interlude. Confirm all three
   roles acknowledge and answer the read-only handoff using subscription access.
3. Run chapter 5 on the target release. Review its host-command approval, verify
   `/mcp`, read a real task and append a result note to the disposable backlog.
4. Install Chromium inside SBX. Inspect any network approval before allowing the
   specific host for the workshop sandbox. Run browser tests and record passed,
   failed and skipped counts. ARM machines may use `ports.ubuntu.com` for packages.
5. Join chapter 6 through SSH and answer the product question. Check that the
   current resolution note clears while history retains the old resolution.
6. Exercise cleanup and a catch-up checkpoint. Keep a prepared feature checkpoint
   available if a coding task takes too long during delivery.

The success criterion is a running change, actual verification evidence, an
inspectable commit and a backlog result note. A role's “done” message is not proof.
Avoid unrelated dependency upgrades on the eve of delivery.

Validation recorded on October 29: SBX 0.49.2 environment plans and local kit
validation passed. In a disposable SBX, fresh materials, checkpoint recovery,
helper tests, Beans MCP protocol tests and the starter app's 19 non-browser tests
passed; the app's health endpoint also responded through a published host port.
The 0.50.0 CLI was checked inside SBX, including removal of `--auto-approve`.
The revised all-Claude configuration was also exercised in a fresh Claude sandbox
with the pinned ACR and Herdr kits: subscription authentication returned a real
Sonnet response, all three roles acknowledged, and coordinator → developer → QA
→ human completed a read-only handoff. Catch-up and launcher regression checks
passed, including paths with spaces and preserved custom role choices.
This does not establish a full 0.50.0 runtime or Windows rehearsal. Chromium
download required approval for `cdn.playwright.dev`; browser checks must be
completed after that access decision. Rehearse the MCP gateway and SSH steps above
on the target host before calling the whole workshop verified.

## Look at the environment as a whole

Prepare one demo sandbox and keep it connected before presenting. The examples
below use `wad-demo`; substitute your prepared sandbox name throughout. In HOST:

```bash
sbx inspect wad-demo
```

Ask the audience to find the kits, published port and workspace/mount configuration.
Connect each item to the chapter where it was introduced. This is useful when
someone asks “what did we actually create?” or “why can this agent reach that?”
Review the output before projecting it because it can identify host paths and
account configuration. It describes the environment, not the success of a task.

## Grant access to a directory while the sandbox exists

Every workshop chapter shares the application directory. Now suppose a worker
needs an additional
reference file: can we grant that access without rebuilding its environment?

Open a shell in the selected sandbox and leave it open in terminal A:

```bash
# HOST — terminal A
sbx exec -it wad-demo bash
```

This also starts it if it was stopped. In host terminal B, prepare a scratch folder
and a message. Run these commands from the workshop repository root.

```bash
# HOST — terminal B
mkdir -p "./.local/mount-demo"
```

Create `.local/mount-demo/message.txt` in your editor with this content:

```text
A message from the host
```

Share it read-only from HOST:

```bash
sbx mount wad-demo "$(pwd)/.local/mount-demo:/home/agent/mailbox:ro"
```

`sbx mount` shares that directory at the
specified path inside the sandbox. The trailing `ro` makes it read-only. Ask the
audience to predict which of these operations will work, then try them in terminal A:

```bash
# SANDBOX shell — terminal A
cat /home/agent/mailbox/message.txt
echo changed > /home/agent/mailbox/message.txt
```

Reading should succeed and writing should fail. Ask the audience which part of
the mount command explains the difference before pointing back to `ro`.

Back on the host, revoke the mount:

```bash
# HOST — terminal B
sbx umount wad-demo "$(pwd)/.local/mount-demo:/home/agent/mailbox"
```

The spelling is `umount`; the arguments identify the same host directory and target
path. Try `cat /home/agent/mailbox/message.txt` again in the sandbox. The host file
still exists, but it is no longer visible through that mount.

A possible extension is to mount the same scratch directory into two sandboxes so
one can leave files for another. Discuss the idea rather than implementing a second
messaging system now. Sharing storage does not by itself define who reads a message,
when an agent wakes up, or how two writers coordinate. Those are the concerns we
already addressed inside the factory with files and Herdr.

Exit the inspection shell and stop the demo sandbox when finished.

## Run somewhere other than your laptop

Use a prepared cloud sandbox on an account with access. The question for the audience
is what changes when the compute is remote and what remains familiar.

```bash
# HOST
sbx --cloud --help
sbx --cloud ls
```

`--cloud` selects the cloud backend. The help output shows which commands that
backend currently supports; the list shows your remote environments. Select the
prepared demo by its name and open its agent:

```bash
# HOST — replace CLOUD_DEMO_NAME with that actual name
sbx --cloud run claude --name CLOUD_DEMO_NAME
```

Talk to the agent and ask it to inspect the prepared app/environment. Exit when you
are ready to show its service. If your prepared service listens on port 8080:

```bash
sbx --cloud ports CLOUD_DEMO_NAME --publish 8080
```

Open the URL returned by the cloud service. Unlike a local `3102:8080` mapping,
cloud publishing names a sandbox port and returns a remote URL.

Ask what would happen to our **workshop task-access server running on the host** in this arrangement.
It does not automatically move with the sandbox or become remotely reachable.
Moving the whole factory requires designing that connection; this demo only shows
the remote execution environment.

Remove the published demo port afterward with
`sbx --cloud ports CLOUD_DEMO_NAME --unpublish 8080`, and stop the cloud environment
using the account's supported lifecycle commands. Attendees do not need cloud access
to finish the workshop.

## Add an organization's rules

On an enrolled presenter account, show a prepared rule for a named tool. Ask the
agent to make a harmless allowed call and a harmless denied one. Read the response
together: what was blocked, who owns that rule, and where could an operator change it?

Relate this to the earlier local network decision and the adapter's limited tool
surface. Compare who configures each control and which operations it governs.
