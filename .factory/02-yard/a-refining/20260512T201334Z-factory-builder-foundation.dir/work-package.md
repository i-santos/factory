# Factory Builder Foundation

## Source Intake

The operator wants to restart the Factory product foundation around a visual interface and deterministic engine for creating Factory primitives: circuits, machines, automations, commands, and event bindings.

The old `evolve-factory-cli` seed command should be removed from, archived away from, or de-emphasized in the default visible product flow. The first product experience should be a Factory Builder where users create concrete building blocks instead of starting from a meta command for evolving the CLI.

## Shared Outcome

Factory has a usable builder foundation: the model, engine, API, and GUI clearly separate commands, circuits, machines, automations, and event bindings; users can create and inspect those primitives without editing JSON or Markdown by hand; and persisted definitions are validated before execution.

## Why This Is A Work Package

This is larger than one work order because it changes several coordinated surfaces:

- core builder model and persistence contracts
- CRUD APIs or service functions for Factory primitives
- validation rules across primitive references
- GUI behavior and visual creation surfaces
- cleanup of the existing `evolve-factory-cli` seed command visibility
- tests for model, persistence, validation, and GUI/API behavior

Those pieces should ship together as one cohesive package, but they can be refined into focused implementation work orders.

## Intended Work Orders

### WO-01: Clarify Factory primitive model and default seed state

Define the first-class model boundaries for commands, circuits, machines, automations, and event bindings. Remove, archive, or hide `evolve-factory-cli` from the default visible product flow.

Acceptance criteria:
- The model distinguishes registered commands from reusable machines.
- `evolve-factory-cli` is no longer presented as the starting point for the product.
- Existing workspace defaults are updated without breaking current validation.
- Documentation or project context explains the primitive boundaries.

### WO-02: Add builder persistence and validation services

Add create, list, read, and update support for circuits, machines, automations, commands, and event bindings using deterministic structured definitions under the Factory workspace.

Acceptance criteria:
- A circuit definition can be created and persisted.
- A machine can be created with one or more circuit references.
- An automation can be created with one or more machine references.
- A command and event binding can be created through the builder layer.
- Invalid references, duplicate IDs, and malformed inputs are rejected with clear validation errors.

### WO-03: Expose Factory Builder API endpoints

Expose the builder service through local GUI/API endpoints so the UI can create and inspect primitives without direct filesystem editing.

Acceptance criteria:
- API endpoints support list/read/create/update for the core primitives.
- API responses use stable JSON shapes.
- Validation errors return clear, user-facing error responses.
- API tests cover successful creation and invalid reference failures.

### WO-04: Build the first visual Factory Builder UI

Add a visual builder surface to the GUI for creating and inspecting circuits, machines, automations, commands, and event bindings.

Acceptance criteria:
- The GUI no longer presents registered commands as `factory run machine ...`.
- Registered commands and reusable machines are visually distinct.
- A user can create a circuit, machine, and automation from HTML controls.
- The UI shows validation errors without requiring the user to inspect Markdown or JSON files manually.
- Created primitives appear in the visual Factory map after creation.

### WO-05: Add builder regression tests and verification coverage

Add tests for model separation, persistence, validation, API behavior, and GUI command/machine affordances.

Acceptance criteria:
- Tests cover command-vs-machine visualization behavior.
- Tests cover create/list/read/update behavior for builder primitives.
- Tests cover invalid machine-to-circuit and automation-to-machine references.
- Tests cover removal or de-emphasis of `evolve-factory-cli` from the default product flow.
- `go test ./...` passes.

## Known Ordering Or Dependency Constraints

1. Clarify model boundaries and default seed state first.
2. Implement builder persistence and validation before API or GUI creation flows.
3. Add API endpoints before wiring interactive GUI forms.
4. Update GUI command/machine behavior alongside the builder UI so users do not see invalid run commands.
5. Finish with regression tests and full verification.

## Completion Signal

The package is complete when a user can open the local Factory GUI and create the first real Factory primitives through visual controls:

- create a circuit
- create a machine that references that circuit
- create an automation that references that machine
- see the primitives represented correctly in the Factory map
- receive visible validation errors for invalid references
- no longer see command-backed nodes masquerading as runnable machines

## Lifecycle Status

refining

## Produced Work Orders

None yet.

## Current Status

- WO-01: candidate
- WO-02: candidate
- WO-03: candidate
- WO-04: candidate
- WO-05: candidate
