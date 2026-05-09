# Build Serial Machine Orchestrator

## Source Work Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir`

## Shared Outcome

Factory can run reusable machines and automations as real program orchestration over circuit prompt-programs.

## Requested Outcome

Implement the first serial orchestrator slice for machines and automations.

## Scope

- Load circuits, machines, and automations from project-local factory state.
- Run machine circuits serially.
- Run automation machines serially.
- Execute circuits through an injectable runner so tests do not require Codex.
- Add CLI entrypoints for machine and automation runs.
- Keep orchestrator-managed parallelism out of scope.

## Ordering Constraints

- Depends on WO-02 runtime kernel contract.
- Depends on WO-03 reusable project model.

## Completion Signal

- Tests prove a machine runs circuits in order.
- Tests prove an automation runs machines in order.
- CLI exposes serial machine and automation execution.
