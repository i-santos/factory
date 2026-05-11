# Model Reusable Factory Pieces

## Source Work Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir`

## Shared Outcome

Factory persists reusable circuits, machines, automations, sectors, and templates as project-local state that both CLI and GUI can read and write.

## Requested Outcome

Define the reusable factory piece model, storage layout, validation rules, and initial CLI surface.

## Scope

- Define the filesystem layout for circuits, machines, automations, sectors, templates, runtime, runs, and triangulation.
- Define minimum fields for each reusable piece.
- Define validation invariants.
- Define import/export expectations for templates.
- Keep implementation out of scope unless required for documentation links.

## Ordering Constraints

- Depends on WO-01 triangulation baseline.
- Enables WO-04 serial orchestrator and WO-06 GUI foundation.

## Completion Signal

- The repository has a reusable factory project model document.
- The document defines each reusable piece and its invariants.
- Architecture and roadmap link to the model.
