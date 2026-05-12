# Gamified Factory Evals

This document defines evaluation expectations for the current triangulation. Evals should protect the product shape before deeper implementation work.

## Autonomous Continuation

The factory should continue automatically only when continuation is deterministic.

Eval expectations:
- a run with a supported next action and valid input continues to the configured next step
- a run with a missing next action stops instead of guessing
- a run with an unsupported next-action kind stops instead of guessing
- a run with invalid next-action input stops with visible error state
- a run with `requiresHuman: true` enters `needs-human-action`

## Human-In-The-Loop

Human decisions should be visible, structured, and low-cognitive-load.

Eval expectations:
- a `needs-human-action` run includes a decision request
- a bounded decision renders as bounded controls, not free text
- submitted input is stored as structured run state
- a submitted decision resumes only through a configured next action
- the user can understand the required action without locating or reading Markdown files

## GUI State Fidelity

The GUI should reflect runtime state directly.

Eval expectations:
- running state appears as active visual work
- succeeded state appears as completed visual work
- blocked state appears as stopped visual work with a reason
- failed state appears as failed visual work with a reason
- `needs-human-action` appears as an actionable visual decision
