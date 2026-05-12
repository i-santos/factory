# Gamified Factory Roadmap

## Phase 0: Product Triangulation

Goal: establish the source of truth before more implementation.

Deliverables:
- product glossary
- non-negotiable use cases
- architecture baseline for CLI plus GUI
- feasible implementation sequence
- acceptance criteria for the first playable factory

Exit criteria:
- roadmap, use cases, and architecture agree on circuits > machines, sectors as organization, and automations as machine workflows.

## Phase 1: Runtime Contract

Goal: define how circuits execute prompt-programs.

Deliverables:
- `runtime-kernel.md`
- pseudocode syntax and semantics
- allowed primitives: `<code>`, `<comment>`, `Runtime.input`, `Assert`, `Tool.call(...)`, `Trace.event(...)`, and `Pseudo.*`
- validation rules for circuit programs
- example circuit programs

Exit criteria:
- a circuit prompt can be parsed, checked against the kernel contract, and executed by an agent with predictable boundaries.

## Phase 2: Reusable Factory Pieces

Goal: persist user-defined circuits, machines, sectors, automations, and templates.

Deliverables:
- filesystem layout
- JSON schemas or equivalent typed model
- CLI create/list/show/update/validate commands
- template export/import
- sample factory definition

Exit criteria:
- users can build reusable pieces without editing source code.

## Phase 3: Serial Orchestrator

Goal: run a machine or automation end-to-end.

Deliverables:
- automation loader
- machine loader
- circuit runner using `codex exec`
- serial execution engine
- input/output handoff between circuits
- execution logs and status records

Exit criteria:
- a sample automation runs multiple machines and circuits in sequence with durable state.

## Phase 4: New Factory Initialization

Goal: create a new factory folder from an init prompt.

Deliverables:
- new factory command
- workspace bootstrap
- init prompt execution
- generated triangulation
- starter circuits, machines, sectors, and automations

Exit criteria:
- an empty folder becomes a usable factory project with documented triangulation.

## Phase 5: Gamified GUI Foundation

Goal: make the factory visible and operable as a game-like production system.

Deliverables:
- factory map/canvas
- visual machines and circuits
- sector grouping
- automation dispatch controls
- live execution states
- success, blocked, failed, and running feedback

Exit criteria:
- a user can inspect and run a sample automation from the GUI.

## Phase 6: CLI And GUI Alignment

Goal: ensure both surfaces operate on the same local factory project.

Deliverables:
- shared project model
- shared validators
- CLI commands for GUI-authored pieces
- GUI editing for CLI-authored pieces

Exit criteria:
- a factory created in the CLI can be opened in the GUI, modified, run, exported, and reused.

## Phase 7: Autonomous Visual Operations

Goal: let users run factories from a remote web UI with low cognitive load, automatic continuation for safe next actions, and visual human-in-the-loop decision collection.

Deliverables:
- durable run-state model with `needs-human-action`
- deterministic next-action envelope for safe automatic continuation
- web run monitor that shows active, completed, blocked, failed, and waiting-for-human states
- decision request model for buttons, option lists, forms, text inputs, and structured panels
- GUI decision submission flow that writes user input back into the run state
- evals for automatic continuation, ambiguous next actions, and human-stop boundaries

Exit criteria:
- a user can start a run from the web UI, watch it continue through configured next actions, and provide requested human input without manually locating and reading Factory Markdown files.
