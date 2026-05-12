# Persist And Load Run Records

## Source Work Package

`.factory/02-yard/a-refining/20260512T135038Z-autonomous-visual-runs.dir`

## Requested Outcome

Command-chain runs persist durable run session records for CLI and GUI inspection.

## Scope

- Persist run sessions under `.factory/runs/`.
- Include steps, status, next action, decision request, and stop reason.
- Add run-list helper for callers.
- Surface persistence failures as command-chain errors.

## Completion Signal

Run sessions are durable and loadable after command-chain execution.
