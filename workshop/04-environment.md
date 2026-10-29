# 4. Recreate the environment

You have been choosing a workspace, resources, kits, a port and a task connection
by hand. An environment file records those choices so your next worker starts
with the same capabilities. It is also something a teammate can read before
approving access.

## 1. Describe what you assembled

Create `factory/sbxenv.yaml` in your editor with this complete definition:

```yaml
schemaVersion: "1"
name: "${{ env.args.name }}"
agent: claude
workspace: "${{ env.args.app }}"
args:
  name:
    default: wad-repeatable
  app:
    default: ../sample-app
  host_shell:
    default: bash
sandboxOptions:
  skills: off
  cpus: 4
  memory: 8g
ports:
  - sandbox: 8080
    host: 3102
    protocol: tcp4
    hostIP: 127.0.0.1
kits:
  - source: ./hello-kit
  - source: "git+https://github.com/shelajev/acr-sbx-kit.git#ref=f446b95bccabd879912db174c190ff09377b7c7b"
mcp:
  servers:
    - name: "${{ env.args.name }}-beans"
      command: "${{ env.args.host_shell }}"
      args: ["${{ env.fileDir }}/../scripts/beans-mcp"]
```

If you chose Codex, set `agent: codex`. The `workspace` is relative to this file,
so `../sample-app` selects your existing project. The kits are the ones you just
used. The MCP command has no note-writing flag: it remains read-only.

`name` lets us reuse the definition for different workers. `app` lets us choose
another project later. `host_shell` names the host program that starts the Beans
adapter. On **Windows**, run this in HOST Git Bash to find Bash's Windows path:

```bash
cygpath -m /bin/bash
```

Set the `host_shell` default in your YAML to that output, quoted, for example
`default: 'C:/Program Files/Git/bin/bash.exe'`. Leave it as `bash` on macOS.

Keep this environment file, kits and MCP script outside `sample-app/`. They
control what your host runs and what access a new sandbox receives; the agent
should not be able to edit them through its application mount.

## 2. Preview, then create

In HOST, preview the environment:

```bash
sbx env plan factory/sbxenv.yaml
```

Find the application path, both kits, the port mapping and the host MCP command.
The plan describes proposed setup; it does not start the application. In the
SANDBOX tab, from the workshop root:

```bash
sbx env run factory/sbxenv.yaml
```

Approve the plan after reading it. SBX 0.50.0 also asks for approval of the host
MCP command. Check that it points to the workshop's `scripts/beans-mcp` and does
not include `--allow-notes`. Approve that command to use the task connection.

On Windows, prepare npm's executable-link storage for this manually created
sandbox with this HOST command before asking the agent to run the app:

```bash
./scripts/prepare-workspace.sh wad-repeatable
```

This keeps npm dependencies inside the sandbox while leaving source and Git
history shared. The team launcher introduced later performs this preparation.

Ask the agent:

> Run workshop-hello and acr --version. Read wad-101 through Beans MCP, then
> inspect the project and its AGENTS.md. Confirm whether the filter-count change
> is present; do not implement it twice. Start the app and its supporting services
> inside SBX on 0.0.0.0:8080. Report what you found and how I can check it.

Open <http://127.0.0.1:3102>. You did not need a separate port-publishing command,
MCP registration or kit-installation conversation. Your environment definition
supplied them. The application still needed starting: an environment definition
is not a saved running process or a database backup.

## 3. Prove you can start again

Exit the agent and remove the environment from HOST:

```bash
sbx env rm factory/sbxenv.yaml
```

Review the destroy plan. In SBX 0.50.0, this also removes its environment-defined
MCP registration. Your source, Git history, kit and YAML file remain.

Run `sbx env run factory/sbxenv.yaml` again in the SANDBOX tab and repeat the
prompt above. On Windows, repeat the workspace-preparation command for the new
sandbox first. The same tool, guidance and task access should be available in a
fresh worker. Exit and remove it with the same `sbx env rm` command when finished.

This is the hands-on checkpoint: you can hand someone your environment definition
and kits, give a worker a task ID, and inspect the resulting code on your laptop.
Keep `factory/` and `sample-app/` to continue later.

Next: [delegate between roles](05-team.md).
