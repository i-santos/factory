# Gamified Factory Triangulation

This triangulation is the current product baseline. It supersedes the earlier command-centric model; compatibility with that model is intentionally out of scope.

Factory is a gamified local automation environment. A user creates one or more factories as folders on their computer, visualizes each factory as a playable production system, and runs reusable prompt-powered circuits through real program orchestration.

## Product Thesis

Factory should feel like a game about building production lines, not a plain workflow admin tool. The user assembles circuits into machines, organizes machines into sectors, and dispatches automations that run machines in sequence.

The orchestrator is deterministic application code. It loads factory configuration, executes machines, invokes circuits through `codex exec`, waits for results, records state, and moves to the next configured step. Circuits are the agent runtime boundary.

## Model

- **Factory:** one filesystem folder/project.
- **Circuit:** reusable agent abstraction with one prompt-program.
- **Circuit program:** pseudocode prompt executed by `codex exec`.
- **Runtime kernel:** contract loaded by every circuit before execution.
- **Machine:** ordered set of circuits.
- **Sector:** organization-only grouping for machines and reusable pieces.
- **Automation:** ordered workflow of machines.
- **Template:** exported reusable factory configuration.

## Triangulation Artifacts

- [Roadmap](roadmap.md)
- [Use Cases](use-cases.md)
- [Architecture](architecture.md)
- [Runtime Kernel](runtime-kernel.md)

## Invariants

- Circuits are the only agent-executed runtime unit.
- Machines and automations are orchestrated by the application, not improvised by an agent.
- The initial orchestrator runs serially.
- Parallel work happens inside a circuit when a circuit explicitly spawns subagents.
- Sectors do not define behavior.
- Every circuit loads a `runtime-kernel.md`.
- Every circuit program must follow the runtime-kernel pseudocode contract.
- New factories start from an init prompt and produce triangulation before deep implementation.
