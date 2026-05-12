# Implement Deterministic Continuation Rules

## Source Work Package

`.factory/02-yard/a-refining/20260512T135038Z-autonomous-visual-runs.dir`

## Requested Outcome

Command-chain execution can continue from supported structured next actions and visibly stop for unsafe continuations.

## Scope

- Accept `needs-human-action` result status.
- Continue automatically from supported `command` or `factory-command` next actions.
- Stop for unsupported next-action kinds.
- Stop for missing next-action commands.
- Stop for human-required next actions.
- Preserve existing event-binding continuation behavior.

## Completion Signal

Supported deterministic next actions can continue, while unsafe next actions do not auto-run.
