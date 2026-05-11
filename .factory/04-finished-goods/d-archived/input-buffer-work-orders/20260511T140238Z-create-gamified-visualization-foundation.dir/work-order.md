# Create The Gamified Visualization Foundation

## Source Work Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir`

## Shared Outcome

Factory becomes a gamified, Factorio-like application where users visualize and operate factories made from reusable circuits, machines, sectors, and automations. The product must support both CLI and GUI workflows, and each factory is a filesystem folder/project with its own initialized workspace, triangulation, reusable pieces, runtime programs, and exportable templates.

## Requested Outcome

A user can open a local, task-focused visual Factory map from the CLI and inspect the current workspace as a playable production system rather than a plain text graph.

## Product Decision

Use a zero-build local web GUI served by the Go CLI for the first foundation slice.

This keeps the GUI foundation aligned with the current Go CLI, avoids committing to a JavaScript application stack too early, and still gives a real browser surface that can later evolve into a richer frontend.

## Scope

- Add a CLI command that serves a local browser UI for the current factory workspace.
- Reuse the existing `BuildFactoryGraph` model as the GUI data source.
- Render a first playable factory map with visible lanes, machines, sectors, actions, and relationships.
- Show runtime-oriented states and labels using existing graph data where available.
- Provide controls or clear command affordances for dispatching machine and automation execution from the same project model.
- Keep the orchestrator deterministic and outside the frontend.
- Add focused tests for the GUI handler/rendering path.
- Document the GUI command and first-slice boundaries.

## Excludes

- Full game mechanics.
- Cloud or collaboration features.
- A committed React/Vite/SPA toolchain.
- GUI authoring of every factory concept.
- Advanced live run streaming.

## Ordering Constraints

- Depends on the reusable project model and serial orchestrator.
- Should build on the current visualization graph rather than introducing a separate GUI-only model.
- Enables WO-07 CLI/GUI alignment work.

## Completion Signal

- A user can run a CLI command and open a local web UI for the current factory workspace.
- The UI renders a visual factory map from real workspace data.
- The UI includes visible machines/circuits/sectors and relationships.
- The UI exposes execution affordances for machine/automation workflows without making the orchestrator agentic.
- Tests verify the GUI handler or renderer returns expected graph-backed UI content.
