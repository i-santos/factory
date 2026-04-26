# CLI Factory Architecture

Purpose: define the next factory model where the CLI creates and runs configurable factories, while agents execute command prompt programs as the runtime.

## Core Decision

A factory is a configured domain, not a hard-coded workflow.

The factory CLI owns:
- creating factory workspaces
- registering commands
- registering event bindings
- registering sectors
- invoking agent runtime sessions
- waiting for structured command results
- deciding the next command from durable event and result contracts

The agent owns:
- executing one command prompt program
- using repository tools to do real work
- producing artifacts and handoffs
- returning a structured result
- recording trace evidence

The CLI must not know what `load-intake`, `drain-work-package`, or any other user-defined command means. It only knows how to resolve command definitions, invoke an agent with the resolved command prompt, validate the result envelope, emit events, and follow configured event bindings.

## Factory Shape

A factory app is the CLI in the repository root. A factory workspace is the `.factory/` directory that stores workflow definitions, handoffs, artifacts, events, logs, and other project-local state.

A factory workspace contains:
- `control-room`: config, command registry, event bindings, schemas, and runtime logs
- `sectors`: domain modules that provide findings and callable actions
- `workspace state`: intake, work packages, active work, shipped artifacts, and archives

Recommended project-local shape:

```text
.factory/
├── 00-control-room/
│   ├── a-config/
│   ├── b-commands/
│   ├── c-schemas/
│   ├── d-events/
│   ├── e-state/
│   └── f-logs/
├── 01-dock/
├── 02-yard/
├── 03-shop-floor/
├── 04-finished-goods/
└── 05-sectors/
    └── <sector-name>/
        ├── findings/
        └── actions/
```

## Command Model

A command is the only executable factory unit.

Steps, audits, reviewers, refiners, drains, loaders, and sector actions are commands at different abstraction levels. They may still be stored in folders named `steps/`, `audits/`, or `actions/` for human orientation, but their runtime boundary is the same:

```json
{
  "name": "command-name",
  "prompt": "project://.factory/00-control-room/b-commands/command-name.md",
  "aliases": [],
  "description": "Human-facing summary.",
  "inputSchema": "project://.factory/00-control-room/c-schemas/command-name.input.schema.json",
  "resultSchema": "factory://schemas/command-result.v1.json",
  "emits": ["command.completed", "command.blocked"],
  "consumes": ["operator.requested"]
}
```

Built-in commands are app-provided defaults. Project commands are user-created commands. The command catalog is the merged registry, not a fixed list compiled into the CLI.

## Event Model

Events connect commands without turning the factory into a rigid state machine.

An event is a durable fact that something happened:

```json
{
  "eventId": "evt_...",
  "eventType": "command.completed",
  "factoryId": "default",
  "command": "load-intake",
  "status": "succeeded",
  "artifactRefs": ["project://.factory/03-shop-floor/d-work-packages/pkg.dir"],
  "data": {
    "package_path": ".factory/03-shop-floor/d-work-packages/pkg.dir"
  }
}
```

An event binding maps an event to the next command:

```json
{
  "on": "command.completed",
  "where": {
    "command": "load-intake",
    "result.status": "work-package-created"
  },
  "run": "drain-work-package",
  "input": {
    "package_path": "$event.data.package_path"
  },
  "mode": "auto"
}
```

The CLI may loop only through this contract:
1. run one command
2. validate its result envelope
3. append a command event
4. find matching event bindings
5. invoke the next configured command
6. stop on no binding, blocked result, failed validation, operator approval requirement, or iteration guard

## Command Result Envelope

Every command must return a deterministic envelope:

```json
{
  "status": "succeeded|blocked|failed|waiting",
  "resultType": "free-form-type-owned-by-command",
  "summary": "Short operator-facing summary.",
  "artifacts": [],
  "events": [],
  "next": {
    "recommendation": "optional-command-name",
    "input": {}
  },
  "missingEvidence": [],
  "errors": []
}
```

The CLI may use only envelope fields, explicit event bindings, and configured schemas for orchestration decisions. It must not parse arbitrary prose from agent output to decide the next autonomous step.

## Sector Model

A sector is a domain module that composes commands.

Sectors own:
- findings: durable project knowledge for retrieval and grounding
- actions: sector-scoped commands
- optional adapters: links to external tools or agent capabilities

Example:

```text
.factory/05-sectors/product/
├── findings/
│   └── experience-principles.md
└── actions/
    └── refine-experience.md
```

A command may depend on a sector action by calling it as another command. The sector does not get a separate execution primitive.

## Boundary Rules

- Commands are executable units.
- Events are orchestration facts.
- Event bindings are orchestration rules.
- Sectors are domain modules.
- Findings are retrieval context, not executable behavior.
- The CLI runs command chains but does not interpret command-specific meaning.
- Agents execute command prompts but do not call the factory CLI to execute nested factory work.
- A command may recommend a next command, but automatic continuation requires a matching event binding or explicit operator request.
