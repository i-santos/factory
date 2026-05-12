# Visual Autonomous Factory Execution

## Source Intake

Operator requested a product-context use case for a highly visual, remotely accessible, Factorio-like Factory where runs continue automatically until human input is explicitly required.

## Shared Outcome

Factory product context defines visual autonomous execution as a first-class capability across use cases, roadmap, and architecture.

The resulting context must preserve the invariant that the orchestrator is deterministic application code. It may continue from structured, configured next actions, but it must not guess agentically. When human input is required, the run must pause in a visible `needs-human-action` state and present HTML controls for the user to provide the requested decision.

## Why This Is A Work Package

The intake affects multiple product-context artifacts and later implementation areas:

- use-case definition
- roadmap sequencing
- architecture invariants
- run-state model
- web GUI decision surfaces
- autonomous continuation rules
- evaluation expectations

Those concerns should be refined into focused work orders instead of being forced into one broad context edit.

## Intended Work Orders

### WO-01: Add visual autonomous execution to triangulation

Update product context so `docs/use-cases.md`, `docs/roadmap.md`, and `docs/architecture.md` describe visual autonomous run chaining, deterministic continuation, and human decision collection.

Acceptance criteria:
- `docs/use-cases.md` includes a use case for autonomous visual run chaining and human decision collection.
- `docs/roadmap.md` includes the capability in a coherent phase.
- `docs/architecture.md` explains deterministic continuation, blocked human-input state, and GUI decision surfaces.
- The context distinguishes automatic continuation from unsafe agent guessing.

### WO-02: Define run-state and next-action model

Specify the durable state shape for configured next actions, run continuation, and `needs-human-action` pauses.

Acceptance criteria:
- Run state can represent `running`, `completed`, `blocked`, `failed`, and `needs-human-action`.
- Structured outputs can carry deterministic next actions.
- Human decision requests have typed options, form fields, labels, and submit semantics.
- Continuation rules are explicit enough for the orchestrator to execute without interpreting prompt prose.

### WO-03: Design visual human-in-the-loop UX

Define the web UI behavior for showing factory progress, blocked runs, and required human decisions with low cognitive load.

Acceptance criteria:
- The GUI can show a run chain visually.
- The GUI can surface a pending human action without requiring Markdown file inspection.
- The GUI can present buttons, option lists, form fields, and free-text inputs when requested by the run state.
- The UX keeps the user oriented toward triangulation, roadmap, architecture, and eval feedback.

### WO-04: Define autonomous continuation evals

Add evaluation expectations for safe automatic continuation and human-stop boundaries.

Acceptance criteria:
- Evals cover a run that continues automatically from a deterministic next action.
- Evals cover a run that pauses for human input.
- Evals reject continuation when the next action is ambiguous or unsupported.
- Evals confirm GUI-visible state reflects runtime state.

## Known Ordering Or Dependency Constraints

1. Update triangulation context first.
2. Define durable run-state and next-action semantics before implementation.
3. Design the GUI decision UX after the state model is clear.
4. Add eval expectations before or alongside implementation work.

## Completion Signal

The package is complete when the product context and refinement outputs provide a coherent implementation path for a remotely accessible, highly visual Factory that autonomously continues safe configured run chains and visually requests human decisions when required.

## Lifecycle Status

refining

## Produced Work Orders

None yet.

## Current Status

- WO-01: candidate
- WO-02: candidate
- WO-03: candidate
- WO-04: candidate
