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

The circuit owns:
- loading `runtime-kernel.md`
- executing one prompt-program
- using tools when requested by `Tool.call(...)`
- reasoning through `Pseudo.*` operations
- returning structured output
- recording trace evidence

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

The exact layout can evolve, but the model must remain stable:

- circuits are reusable prompt runtimes
- machines are ordered circuit sets
- automations are ordered machine workflows
- sectors are organization-only
- templates are exportable factory definitions

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

The map should show execution state directly: idle, running, blocked, failed, and completed.

## CLI Architecture

The CLI exposes the same model without requiring the GUI.

Initial command families:
- `factory init`
- `factory circuit ...`
- `factory machine ...`
- `factory automation ...`
- `factory sector ...`
- `factory run machine <name>`
- `factory run automation <name>`
- `factory template export`
- `factory template import`
- `factory validate`

## Architectural Invariants

- A circuit is the only agent-executed unit.
- A machine is not an agent; it is an ordered circuit definition.
- An automation is not an agent; it is an ordered machine workflow.
- A sector never constrains execution.
- The orchestrator must not interpret prompt-specific meaning.
- Prompt-programs must follow the runtime kernel.
- Factory state must be reusable and portable across local projects.
