# Add Autonomous Visual Run Tests

## Source Work Package

`.factory/02-yard/a-refining/20260512T135038Z-autonomous-visual-runs.dir`

## Requested Outcome

Tests cover deterministic continuation, human-stop behavior, persisted run records, and GUI/API exposure.

## Scope

- Test supported next-action continuation.
- Test unsupported next-action blocking.
- Test human-required next-action persistence.
- Test `needs-human-action` parsing.
- Test GUI/API run exposure.

## Completion Signal

`go test ./...` covers the new autonomous visual run foundation.
