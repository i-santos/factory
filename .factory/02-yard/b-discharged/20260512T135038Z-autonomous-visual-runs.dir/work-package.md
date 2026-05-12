# Autonomous Visual Runs

## Source Intake

Operator requested an implementation work package for the first usable slice of the gamified orchestrator app: durable run state, deterministic autonomous continuation, and GUI-visible human-in-the-loop decision requests.

## Shared Outcome

Factory can record and expose visual run state for autonomous command/circuit execution. A supported deterministic next action can continue automatically, while unsupported, missing, ambiguous, or human-required next actions stop visibly.

This package should produce a usable foundation, not the full remote game UI. The goal is the smallest implementation slice that makes the newly shipped product context executable and visible.

## Why This Is A Work Package

The intake spans multiple coordinated implementation areas:

- typed Go model additions
- command/circuit result semantics
- durable run persistence
- deterministic continuation validation
- GUI/API state exposure
- tests for autonomous and human-stop boundaries

Those areas are related and should ship together as one package, but they should be refined into focused work orders.

## Intended Work Orders

### WO-01: Add run-state, next-action, and decision-request models

Add typed Go models for durable run status, structured next actions, decision requests, and supported decision controls.

Acceptance criteria:
- Run status includes `idle`, `running`, `succeeded`, `blocked`, `failed`, and `needs-human-action`.
- A typed next-action model includes kind, input, and requires-human semantics.
- A typed decision-request model supports bounded options, text, textarea, checkbox, form-like controls, and submit actions.
- JSON serialization matches `docs/project-model.md`.

### WO-02: Implement deterministic continuation rules

Extend command/circuit run handling so structured result data can request a supported next action and unsafe results stop visibly.

Acceptance criteria:
- Supported deterministic next actions can continue automatically.
- Missing, unsupported, invalid, ambiguous, or human-required next actions do not auto-run.
- Human-required output records `needs-human-action`.
- Existing command-chain behavior remains compatible.

### WO-03: Persist and load run records

Persist run records under the Factory workspace so CLI and GUI can inspect current and historical runs.

Acceptance criteria:
- Run records are written under `.factory/runs/` or the configured workspace run directory.
- Records include status, steps, next action, and decision request when present.
- Load/list helpers expose persisted runs to callers.
- Persistence failures surface as errors instead of being hidden.

### WO-04: Expose run state through GUI/API

Add the first GUI/API surface for run monitor data and human decision panels.

Acceptance criteria:
- An API endpoint exposes run records as JSON.
- The GUI can render run statuses including `needs-human-action`.
- The GUI can render a minimal decision panel from decision-request data.
- The UI avoids requiring users to inspect Markdown files for the current run decision.

### WO-05: Add autonomous visual run tests

Add tests for continuation and human-stop boundaries.

Acceptance criteria:
- Tests cover successful automatic continuation from a supported next action.
- Tests cover stop-on-missing or unsupported next action.
- Tests cover `needs-human-action` with a decision request.
- Tests cover GUI/API exposure of run state data.

## Known Ordering Or Dependency Constraints

1. Add models first.
2. Implement continuation rules against those models.
3. Add persistence so GUI/API reads durable state.
4. Expose GUI/API after persisted state exists.
5. Add or extend tests throughout, with final package verification by `go test ./...`.

## Completion Signal

The package is complete when the app has a tested foundation for durable autonomous visual runs:

- typed run-state and decision models exist
- deterministic next actions can continue safely
- unsafe next actions stop visibly
- human-required output becomes `needs-human-action`
- GUI/API can expose run state and decision-request data

## Lifecycle Status

release-ready

## Produced Work Orders

- `.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T141646Z-add-run-state-next-action-models.dir/work-order.md`
- `.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T141647Z-implement-deterministic-continuation-rules.dir/work-order.md`
- `.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T141648Z-persist-load-run-records.dir/work-order.md`
- `.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T141649Z-expose-run-state-gui-api.dir/work-order.md`
- `.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T141650Z-add-autonomous-visual-run-tests.dir/work-order.md`

## Release-Ready Evidence

- `.factory/04-finished-goods/b-release-ready/20260512T141646Z-run-state-next-action-models.md`
- `.factory/04-finished-goods/b-release-ready/20260512T141647Z-deterministic-continuation-rules.md`
- `.factory/04-finished-goods/b-release-ready/20260512T141648Z-run-record-persistence.md`
- `.factory/04-finished-goods/b-release-ready/20260512T141649Z-run-state-gui-api.md`
- `.factory/04-finished-goods/b-release-ready/20260512T141650Z-autonomous-visual-run-tests.md`
- `.factory/04-finished-goods/b-release-ready/20260512T141651Z-autonomous-visual-runs-package.md`

## Current Status

- WO-01: integrated
- WO-02: integrated
- WO-03: integrated
- WO-04: integrated
- WO-05: integrated
