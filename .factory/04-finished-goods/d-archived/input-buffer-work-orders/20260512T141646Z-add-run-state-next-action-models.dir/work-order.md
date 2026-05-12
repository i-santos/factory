# Add Run-State, Next-Action, And Decision-Request Models

## Source Work Package

`.factory/02-yard/a-refining/20260512T135038Z-autonomous-visual-runs.dir`

## Requested Outcome

Typed Go models represent durable run status, structured next actions, decision requests, and supported decision controls.

## Scope

- Add run status constants.
- Add `NextAction`.
- Add `DecisionRequest`, `DecisionControl`, and options.
- Add durable `RunSession`.
- Add circuit and command result support for next actions and decision requests.

## Completion Signal

JSON model shape aligns with the product model and supports `needs-human-action`.
