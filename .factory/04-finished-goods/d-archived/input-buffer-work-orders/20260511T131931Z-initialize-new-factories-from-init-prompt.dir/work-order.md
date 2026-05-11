# Initialize New Factories From An Init Prompt

## Source Work Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir`

## Shared Outcome

Factory becomes a gamified, Factorio-like application where users visualize and operate factories made from reusable circuits, machines, sectors, and automations. The product must support both CLI and GUI workflows, and each factory is a filesystem folder/project with its own initialized workspace, triangulation, reusable pieces, runtime programs, and exportable templates.

## Requested Outcome

A user can initialize a new Factory workspace that is immediately inspectable, versioned, and ready to run an explicit init automation that produces starter triangulation.

## Product Decision

`factory init` scaffolds by default. It must not execute Codex automatically.

The explicit run path should be one of:

- `factory run automation init`
- a dedicated follow-up flag such as `factory init --run`, only if the implementation keeps the non-running scaffold path as the default

## Scope

- Update or extend workspace initialization so a new workspace is created at the latest supported schema version.
- Persist explicit workspace version metadata, including `workspaceSchema`, `runtimeContract`, and the factory skill/runtime version needed by the project.
- Scaffold starter runtime files needed by circuit execution, including the runtime kernel reference or project-local runtime file according to the current resource model.
- Scaffold a starter init circuit that can run a prompt-program through the circuit runtime.
- Scaffold a starter init machine and init automation that wrap the init circuit.
- Ensure the default `factory init` command leaves the workspace in an inspectable, non-running state.
- Add or document an explicit command path to run the generated init automation.
- Produce starter triangulation artifacts when the init automation is run, not during default scaffolding.
- Add update planning from the beginning: define how `factory update` detects old workspace schemas, applies idempotent migrations, and preserves migration evidence/backups before structural changes.

## Excludes

- Rich GUI onboarding.
- Full GUI visualization work.
- Agentic orchestrator behavior.
- Permanent side-by-side active workspace roots such as `.factory/v2` or `.factory/v3`.
- Running generated init automation implicitly during default `factory init`.

## Ordering Constraints

- Depends on the triangulation baseline, circuit runtime kernel contract, reusable factory project model, and serial orchestrator.
- Must preserve the invariant that the orchestrator is deterministic application code and circuits are the only agent-executed runtime unit.
- Must keep one canonical active workspace root, normally `.factory/`, with schema-versioned migration behavior rather than permanent versioned active roots.

## Completion Signal

- `factory init` creates a latest-schema workspace in an empty target folder without invoking Codex.
- The generated workspace contains inspectable starter circuit, machine, automation, runtime/config metadata, and update/migration metadata.
- An explicit run command can execute the starter init automation through the serial orchestrator.
- Running the starter init automation produces initial triangulation artifacts.
- Re-running initialization or update behavior is idempotent or blocks with a clear recovery message.
- Tests cover default scaffold-only initialization, explicit init automation execution using an injectable/mock runner, and schema/update metadata behavior.
