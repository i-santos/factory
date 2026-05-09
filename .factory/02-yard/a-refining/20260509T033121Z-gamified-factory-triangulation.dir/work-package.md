# Gamified Factory Triangulation

## Source Intake

Operator request received on 2026-05-09:

- Produce the project triangulation before further implementation.
- Treat this request as the new source of truth and allow it to override the current command-centric model.
- Start from the supplied use cases.
- Build a feasible roadmap.
- Define the architecture around CLI plus GUI.
- Keep compatibility concerns out of scope for this pivot.

## Shared Outcome

Factory becomes a gamified, Factorio-like application where users visualize and operate factories made from reusable circuits, machines, sectors, and automations. The product must support both CLI and GUI workflows, and each factory is a filesystem folder/project with its own initialized workspace, triangulation, reusable pieces, runtime programs, and exportable templates.

The first concrete product artifact is a triangulation package that aligns:

- product roadmap
- use cases and interaction model
- architecture and runtime contracts

## Product Model

- **Circuit:** reusable agent abstraction that runs one prompt-program through `codex exec`.
- **Circuit program:** prompt written in well-defined pseudocode using a loaded `runtime-kernel.md`.
- **Runtime kernel:** simple structured execution contract loaded by every circuit. It defines tags and primitives such as `<code>`, `<comment>`, `Runtime.input`, `Tool.call(...)`, `Factory.call(...)`, `Trace.event(...)`, `Assert`, and `Pseudo.*(...)`.
- **Machine:** ordered set of one or more circuits. A machine runs end-to-end; circuit outputs usually feed later circuits, but shared input and circuit-local branching are allowed.
- **Sector:** organizational grouping only. Sectors do not constrain execution behavior.
- **Automation:** configured workflow linking machines. The orchestrator dispatches an automation, runs each machine in order, and continues until the automation is complete or blocked.
- **Factory:** one filesystem folder/project containing workspace state, circuits, machines, sectors, automations, templates, and triangulation.

## Non-Negotiable Use Cases

1. Users must visualize a gamified factory. It should feel like a real game, similar in spirit to Factorio.
2. Users must be able to create new circuits and machines.
3. Circuits have programs. Machines are sets of circuits.
4. Circuits are the agent abstraction and run prompt-programs written in pseudocode.
5. Circuits must be treated as runtimes. Every circuit loads a `runtime-kernel.md`.
6. The runtime kernel defines a simple, pragmatic, structured pseudocode contract.
7. Programs must follow the runtime-kernel contract.
8. The contract must execute deterministically enough for orchestration, while still allowing agent reasoning for `Pseudo.*`.
9. The syntax includes `<code>`, `<comment>`, `Pseudo.*`, `Tool.call(tool_ref, input)`, and other kernel-defined primitives.
10. The factory orchestrates machines.
11. Machines run circuits in serial by default. Circuit output may feed the next circuit, while parallel work happens inside a circuit through subagents.
12. The orchestrator is a real program, not an agent.
13. The orchestrator initially runs machines serially and waits for each machine to finish.
14. Circuits are executed as `codex exec "<prompt in pseudocode>"`.
15. Automations link machines and allow an operation to continue across multiple machines.
16. Sectors organize machines but do not define behavior or execution boundaries.
17. Users can create several factories.
18. Creating a new factory initializes a workspace and creates triangulation from an init prompt.
19. Users can reuse circuits, machines, automations, sectors, and whole factory templates.

## Why This Is A Work Package

This is larger than one delivery unit. It changes the product model, terminology, runtime contract, storage model, CLI behavior, GUI behavior, orchestration engine, initialization flow, visualization strategy, and documentation. It must first produce triangulation artifacts, then split implementation into coherent work orders.

## Intended Work Orders

### WO-01: Produce Product Triangulation Artifacts

- **Outcome:** A committed triangulation set exists and supersedes the previous command-centric direction.
- **Scope:** Write product roadmap, use cases, architecture, terminology, and initial acceptance criteria.
- **Includes:**
  - roadmap from current state to playable gamified factory
  - use-case catalog for circuits, machines, sectors, automations, templates, and factory creation
  - architecture for CLI plus GUI
  - glossary and invariant rules
  - migration note that compatibility is intentionally out of scope
- **Excludes:** Implementation changes beyond docs/artifacts.
- **Completion Signal:** The repository has explicit roadmap, use-case, and architecture documents that agree on circuits > machines > sectors/automations.

### WO-02: Define The Circuit Runtime Kernel Contract

- **Outcome:** A minimal deterministic prompt-program contract exists for circuits.
- **Scope:** Design `runtime-kernel.md`, allowed syntax, primitive semantics, trace requirements, and failure rules.
- **Includes:**
  - `<code>` and `<comment>` semantics
  - `Runtime.input`, `Tool.call`, `Trace.event`, `Assert`, and `Pseudo.*`
  - deterministic execution expectations
  - pseudocode examples
  - validation strategy for prompt-programs
- **Excludes:** Full orchestrator implementation.
- **Completion Signal:** Example circuit programs can be checked against the kernel contract.

### WO-03: Model Reusable Factory Pieces

- **Outcome:** The app can persist circuits, machines, automations, sectors, and templates as reusable project state.
- **Scope:** Define filesystem layout, schemas, CRUD commands, and validation.
- **Includes:**
  - circuit registry and prompt-program storage
  - machine definitions with ordered circuit lists
  - automation definitions with ordered machine lists
  - sectors as organization-only metadata
  - factory template export/import format
- **Excludes:** GUI authoring.
- **Completion Signal:** CLI can create, list, validate, export, and import reusable pieces.

### WO-04: Build The Serial Orchestrator

- **Outcome:** A real program can execute one machine or automation end-to-end.
- **Scope:** Implement orchestration over machines and circuits using `codex exec`.
- **Includes:**
  - load automation
  - load each machine
  - load each circuit program
  - run circuits serially
  - pass outputs between circuits when configured
  - wait for each circuit and machine completion
  - record execution state, logs, and failures
- **Excludes:** Orchestrator-managed parallel execution; parallelism remains inside circuit/subagent behavior.
- **Completion Signal:** A sample automation runs multiple machines and circuits serially with durable results.

### WO-05: Initialize New Factories From An Init Prompt

- **Outcome:** A user can create a new factory folder/project and get initial triangulation.
- **Scope:** Design and implement factory creation flow.
- **Includes:**
  - choose/create folder
  - initialize workspace
  - run init prompt through the circuit runtime
  - produce initial triangulation
  - persist starter circuits, machines, sectors, and automations
- **Excludes:** Rich GUI onboarding.
- **Completion Signal:** A new empty folder becomes a usable factory project with triangulation artifacts.

### WO-06: Create The Gamified Visualization Foundation

- **Outcome:** The GUI shows factories as a playable visual system, not a plain admin dashboard.
- **Scope:** Build initial visual metaphor and interaction model.
- **Includes:**
  - factory map/canvas
  - circuits and machines as visible game objects
  - sector organization
  - automation execution status
  - controls/buttons to dispatch machine or automation execution
  - visual runtime feedback for running, blocked, succeeded, and failed states
- **Excludes:** Full game mechanics beyond useful visualization.
- **Completion Signal:** A user can inspect and run a sample automation from the GUI.

### WO-07: Align CLI And GUI Around The Same Project Model

- **Outcome:** CLI and GUI operate on identical factory folders and reusable pieces.
- **Scope:** Ensure both surfaces read/write the same schemas and execution state.
- **Includes:**
  - shared storage contracts
  - shared validation
  - CLI commands for all GUI-authored concepts
  - GUI controls for CLI-authored concepts
- **Excludes:** Advanced collaboration/cloud sync.
- **Completion Signal:** A factory created in CLI can be opened in GUI and vice versa.

## Known Ordering Or Dependency Constraints

1. WO-01 must happen first because implementation needs the triangulation baseline.
2. WO-02 should happen before WO-04 because the orchestrator needs a stable circuit contract.
3. WO-03 should happen before WO-04 and WO-06 because reusable pieces need schemas before execution and visualization.
4. WO-05 depends on WO-01, WO-02, and the first usable version of WO-03.
5. WO-06 can begin with mock/sample state after WO-01, but production wiring depends on WO-03 and WO-04.
6. WO-07 follows the first CLI and GUI slices and keeps both surfaces aligned.

## Current Status

| Work Order | Status | Readiness |
| --- | --- | --- |
| WO-01 | integrated | done |
| WO-02 | integrated | done |
| WO-03 | integrated | done |
| WO-04 | candidate | ready |
| WO-05 | candidate | ready |
| WO-06 | candidate | ready for design baseline |
| WO-07 | candidate | waiting-for first CLI and GUI slices |

## Package Processing State

Status: `blocked-or-waiting`

The package has produced and integrated the triangulation baseline. The remaining work orders now cross runtime contract design, reusable project model, orchestrator implementation, new-factory initialization, GUI foundation, and CLI/GUI alignment. Continue by selecting the next explicit refinement target instead of forcing all implementation slices in one unattended run.

Recommended next candidates:

- WO-02: Define The Circuit Runtime Kernel Contract
- WO-03: Model Reusable Factory Pieces
- WO-06: Create The Gamified Visualization Foundation

## Produced Work Orders

- WO-01: `.factory/03-shop-floor/a-input-buffer/20260509T033940Z-produce-product-triangulation-artifacts.dir/work-order.md`
- WO-02: `.factory/03-shop-floor/a-input-buffer/20260509T035600Z-define-circuit-runtime-kernel-contract.dir/work-order.md`
- WO-03: `.factory/03-shop-floor/a-input-buffer/20260509T035806Z-model-reusable-factory-pieces.dir/work-order.md`

## Integrated Work Orders

- WO-01: Product triangulation artifacts integrated from `factory/order/20260509T033940Z-produce-product-triangulation-artifacts` into `factory/package/20260509T033121Z-triangulation-strategy`.
  - Merge commit: `3c2f71031369558847765530fc587b6ab10d7278`
  - Release-ready evidence: `.factory/04-finished-goods/b-release-ready/20260509T034200Z-gamified-factory-triangulation-docs.md`
- WO-02: Circuit runtime kernel contract integrated from `factory/order/20260509T035600Z-define-circuit-runtime-kernel-contract` into `factory/package/20260509T033121Z-triangulation-strategy`.
  - Merge commit: `35ea801be173515d2406e60ef6b01db3d99dae48`
  - Release-ready evidence: `.factory/04-finished-goods/b-release-ready/20260509T035900Z-circuit-runtime-kernel-contract.md`
- WO-03: Reusable factory project model integrated from `factory/order/20260509T035806Z-model-reusable-factory-pieces` into `factory/package/20260509T033121Z-triangulation-strategy`.
  - Merge commit: `b43306f492b18bf4be56d9e48ec9e483ae58c934`
  - Release-ready evidence: `.factory/04-finished-goods/b-release-ready/20260509T040100Z-reusable-factory-project-model.md`

## Completion Signal

This package is complete when:

- The project has a committed triangulation baseline.
- The product model is circuits > machines, with sectors as organization and automations as machine workflows.
- The architecture clearly separates the real-program orchestrator from agent/circuit runtime execution.
- The initial roadmap is feasible and split into implementation work orders.
- CLI and GUI requirements are both explicit.
- The previous command-centric architecture has been intentionally superseded.

## Next Step

```text
<factory refine-package>
{ "package_path": ".factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir" }
</factory>
```
