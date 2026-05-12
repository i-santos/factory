# Gamified Factory Architecture

This architecture supersedes the earlier command-centric model. Compatibility with that model is intentionally out of scope for the current pivot.

## System Shape

Factory has two product surfaces:

- **CLI:** creates, validates, runs, imports, exports, and inspects local factory projects.
- **GUI:** visualizes and operates factories as a gamified production system.

Both surfaces read and write the same factory folder. One factory equals one filesystem project.

## Runtime Responsibilities

The orchestrator is real application code. It is not an agent.

The orchestrator owns:
- loading factories
- loading automations
- loading machines
- loading circuits
- invoking circuits with `codex exec`
- passing configured input and output between circuits
- waiting for each circuit and machine to finish
- recording execution state
- stopping on failure or blocked output
- continuing from configured next actions when the result is deterministic
- surfacing human decision requests as durable run state

The circuit owns:
- loading `runtime-kernel.md`
- executing one prompt-program
- using tools when requested by `Tool.call(...)`
- reasoning through `Pseudo.*` operations
- returning structured output
- recording trace evidence

The orchestrator must not infer intent from prose in a circuit result. Automatic continuation is allowed only when a result includes a structured next-action envelope that matches configured machine, automation, or factory-command policy. If the next action is missing, ambiguous, unsupported, or marked as human-required, the orchestrator records a blocked or `needs-human-action` state instead of guessing.

## Execution Model

Automation execution is serial at the orchestrator level:

```text
automation
  machine 1
    circuit 1
    circuit 2
  machine 2
    circuit 1
```

The orchestrator waits for each circuit to finish before invoking the next circuit. The orchestrator waits for each machine to finish before invoking the next machine.

Parallel execution is not an orchestrator concern in the first version. If parallel work is needed, a circuit may spawn subagents as part of its own prompt-program.

## Autonomous Continuation

Autonomous continuation is deterministic application behavior, not agentic improvisation.

A completed run step may request continuation with a structured next action such as:

```json
{
  "kind": "process-work-package",
  "input": {
    "packagePath": ".factory/02-yard/a-refining/example.dir"
  },
  "requiresHuman": false
}
```

The orchestrator may execute that next action only when:

- the next-action kind is supported by the factory configuration
- the input validates against the action schema
- the current run policy allows that transition
- the result does not require human input
- the transition does not target a protected base branch directly

If any condition fails, the orchestrator records a visible non-continuing state. The default failure mode is to stop and ask for human input, not to continue with a best guess.

## Human-In-The-Loop State

Human input is a first-class runtime state.

Run status values include:

- `idle`
- `running`
- `succeeded`
- `blocked`
- `failed`
- `needs-human-action`

A `needs-human-action` run stores a decision request with:

- a human-readable title
- explanation of why input is required
- available options when the choice is bounded
- form fields when structured input is needed
- free-text fields only when bounded controls are insufficient
- submit target that resumes or updates the run

The GUI renders this state as visual work on the factory map and in the run monitor. The user should not need to inspect Markdown records, logs, or internal folders to understand what decision is needed.

## Runtime Kernel

Every circuit loads a runtime kernel before executing its program. The kernel defines the syntax and execution contract for prompt-programs. The detailed contract is documented in [Circuit Runtime Kernel Contract](runtime-kernel.md).

Initial primitives:

- `<code>`: executable pseudocode block
- `<comment>`: guidance and non-executable context
- `Runtime.input`: input supplied by the orchestrator
- `Assert`: deterministic guard that blocks execution on failure
- `Tool.call(tool_ref, input)`: external tool/script invocation
- `Trace.event(event)`: durable runtime trace
- `Pseudo.*`: pseudocode operation inferred by the circuit agent within the kernel boundary

The contract must be deterministic enough for reliable orchestration, while preserving agent reasoning for high-level pseudocode operations.

## Project Storage

Recommended project shape:

```text
factory-project/
├── .factory/
│   ├── circuits/
│   ├── machines/
│   ├── automations/
│   ├── sectors/
│   ├── templates/
│   ├── runtime/
│   │   └── runtime-kernel.md
│   ├── runs/
│   └── triangulation/
```

The exact layout can evolve, but the model must remain stable. The reusable project model is specified in [Reusable Factory Project Model](project-model.md):

- circuits are reusable prompt runtimes
- machines are ordered circuit sets
- automations are ordered machine workflows
- commands are operator-facing prompt-program entries
- event bindings connect command results to configured follow-up commands
- sectors are organization-only
- templates are exportable factory definitions

The builder and GUI must keep commands distinct from machines. Registered commands use `factory run <command>`. Reusable machines use `factory run machine <name>`.

## New Factory Flow

Creating a factory starts from a folder and an init prompt.

Flow:

1. create or select factory folder
2. initialize workspace
3. create or select init circuit
4. run the init prompt through the circuit runtime
5. produce initial triangulation
6. persist starter circuits, machines, sectors, and automations
7. make the factory available in CLI and GUI

## GUI Architecture

The GUI must be built around the factory map, not around forms first.

Primary views:
- factory map
- circuit editor
- machine editor
- automation editor
- sector organizer
- run monitor
- template import/export

The map should show execution state directly: idle, running, blocked, failed, completed, and needs-human-action.

Human decision collection is part of the GUI surface. The GUI should render requested decisions with HTML controls such as buttons, option lists, forms, text inputs, or structured decision panels. These controls submit structured input back into the run state so the orchestrator can resume only through configured transitions.

## CLI Architecture

The CLI exposes the same model without requiring the GUI.

Initial command families:
- `factory init`
- `factory circuit ...`
- `factory machine ...`
- `factory automation ...`
- `factory sector ...`
- `factory run <command>`
- `factory run machine <name>`
- `factory run automation <name>`
- `factory template export`
- `factory template import`
- `factory validate`

## Architectural Invariants

- A circuit is the only agent-executed unit.
- A machine is not an agent; it is an ordered circuit definition.
- An automation is not an agent; it is an ordered machine workflow.
- A command is not a reusable machine; it is an operator-facing registry entry.
- A sector never constrains execution.
- The orchestrator must not interpret prompt-specific meaning.
- The orchestrator may continue automatically only from configured structured next actions.
- Human-required, ambiguous, unsupported, failed, or blocked outputs stop the run visibly.
- Prompt-programs must follow the runtime kernel.
- Factory state must be reusable and portable across local projects.
