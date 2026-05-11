# Gamified Factory Use Cases

## UC-01: Visualize A Factory

The user opens a factory project and sees a gamified factory map. Machines, circuits, sectors, and automation state are visible as production objects.

Acceptance criteria:
- the UI shows machines and circuits as first-class visual objects
- execution state is visible without reading logs
- sectors help organization but do not imply runtime constraints

## UC-02: Create A Circuit

The user creates a reusable circuit with a prompt-program.

Acceptance criteria:
- the circuit stores metadata, input/output expectations, and a program path
- the program follows `runtime-kernel.md`
- the circuit can be validated independently

## UC-03: Create A Machine

The user creates a machine from one or more circuits.

Acceptance criteria:
- circuit order is explicit
- output handoff rules are configurable
- the machine can run end-to-end through the orchestrator

## UC-04: Run A Circuit

The orchestrator invokes a circuit by running `codex exec` with the circuit program and runtime context.

Acceptance criteria:
- every circuit loads the runtime kernel
- the program uses structured pseudocode
- circuit output is captured as structured state

## UC-05: Run A Machine

The orchestrator runs a machine by executing its circuits serially.

Acceptance criteria:
- the orchestrator waits for each circuit to finish
- failures stop the machine with a visible blocked state
- later circuits can consume earlier circuit output

## UC-06: Run An Automation

The user dispatches an automation from CLI or GUI.

Acceptance criteria:
- automation defines an ordered list of machines
- the orchestrator runs each machine in order
- the automation records progress and final status

## UC-07: Organize With Sectors

The user groups machines or reusable pieces into sectors for navigation.

Acceptance criteria:
- sectors do not restrict automation links
- automations reference machines, not sectors
- the GUI uses sectors for organization and visualization

## UC-08: Create A New Factory

The user creates a new factory as a folder/project.

Acceptance criteria:
- the app initializes the workspace
- the app runs an init prompt
- the app creates initial triangulation
- starter pieces are stored as reusable project state

## UC-09: Reuse Pieces

The user reuses circuits, machines, automations, sectors, and templates across factories.

Acceptance criteria:
- pieces can be exported and imported
- templates preserve reusable factory structure
- imported pieces can be validated before running

## UC-10: Keep Runtime Scope Narrow

The user composes many focused circuits instead of one broad prompt.

Acceptance criteria:
- machines can contain one circuit or many circuits
- circuit scope is intentionally narrow
- orchestration supplies grain-control over the workflow
