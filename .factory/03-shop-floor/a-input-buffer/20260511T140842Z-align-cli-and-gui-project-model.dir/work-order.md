# Align CLI And GUI Around The Same Project Model

## Source Work Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir`

## Shared Outcome

Factory becomes a gamified, Factorio-like application where users visualize and operate factories made from reusable circuits, machines, sectors, and automations. The product must support both CLI and GUI workflows, and each factory is a filesystem folder/project with its own initialized workspace, triangulation, reusable pieces, runtime programs, and exportable templates.

## Requested Outcome

The CLI and first GUI slice operate on the same local factory project model, with shared validation and consistent inspection behavior.

## Scope

- Ensure CLI visualization, GUI rendering, and orchestration load the same reusable pieces from the workspace.
- Add validation for the reusable pieces needed by both CLI and GUI.
- Expose validation through the CLI.
- Make GUI/API failures surface model validation problems clearly.
- Add tests proving CLI and GUI paths use the same workspace model.
- Document the shared model and validation command.

## Excludes

- Full GUI editing.
- Cloud sync.
- Advanced collaboration.
- Replacing the zero-build GUI foundation with a SPA stack.

## Ordering Constraints

- Depends on the first CLI and GUI slices.
- Must not introduce a separate GUI-only schema.
- Must keep one canonical active workspace root.

## Completion Signal

- A factory initialized in the CLI can be validated and opened in the GUI.
- The same reusable circuit, machine, automation, and sector records appear through CLI graph JSON and GUI graph API.
- Invalid machine/circuit or automation/machine references are reported by validation.
- Tests cover the shared validation and graph behavior.
