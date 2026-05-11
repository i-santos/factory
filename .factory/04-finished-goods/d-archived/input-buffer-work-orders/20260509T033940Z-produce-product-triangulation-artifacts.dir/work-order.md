# Produce Product Triangulation Artifacts

## Source Work Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir`

## Shared Outcome

Factory becomes a gamified, Factorio-like application where users visualize and operate factories made from reusable circuits, machines, sectors, and automations. The product supports CLI and GUI workflows, and each factory is a filesystem folder/project with reusable runtime pieces and exportable templates.

## Requested Outcome

Create the initial triangulation baseline that supersedes the previous command-centric model and aligns roadmap, use cases, architecture, terminology, and acceptance criteria.

## Scope

- Add product roadmap documentation.
- Add use-case catalog documentation.
- Add gamified factory architecture documentation.
- Update the existing architecture entrypoint to point at the new source of truth.
- Keep implementation changes out of scope.

## Ordering Constraints

- Must be completed before runtime-kernel, reusable-piece, orchestrator, new-factory, GUI, and CLI/GUI-alignment implementation work.

## Completion Signal

- The repository has explicit roadmap, use-case, and architecture documents that agree on circuits > machines > sectors/automations.
- The architecture separates the real-program orchestrator from circuit runtime execution.
- The docs state that compatibility with the earlier command-centric model is intentionally out of scope.
