# Configurable Agent Factory Pivot

## Disposition

Status: `superseded-archived`

This package is superseded by the gamified Factory triangulation package completed on 2026-05-11:

`.factory/04-finished-goods/b-release-ready/20260511T141300Z-gamified-factory-triangulation-package.md`

The record is preserved for historical context. It should not be processed as active work unless a future operator explicitly reopens selected ideas under the current gamified Factory product model.

## Source Intake

Operator request received on 2026-04-28:

- Pivot the app away from hard-coded workflow structure.
- Rename current `commands` concept to `agents`.
- Agents have pseudocode prompts.
- Keep events, sectors, and findings.
- Introduce departments as configurable groups of agents.
- Departments exist for organization and visualization only.
- Current dock, yard, shop-floor, and finished-goods should become an initial template, not built-in application behavior.
- Do not preserve compatibility with the current early-development model.
- Document the final version later, not the pivot history.

## Shared Outcome

Factory becomes a configurable runtime where users define agents, departments, events, sectors, and findings through CLI-managed project state. Visualization must render the configured factory structure, not hard-coded docks, yards, shop floors, or finished goods.

## Why This Is A Work Package

The pivot crosses the CLI, internal data model, workspace initialization, runtime execution, visualization, tests, and current dogfood workspace state. It is coherent under one architecture outcome, but too broad for one delivery work order.

## Why This Is Not A Bulk Load

The intake is already shaped enough to identify delivery slices. No discovery pass is required before creating executable work orders.

## Intended Work Orders

### WO-01: Rename Commands To Agents In The Core Model

- **Outcome:** The app model exposes agents as the primary executable unit.
- **Scope:** Replace command-facing types, registry fields, CLI nouns, and runtime identifiers with agent-facing equivalents.
- **Includes:**
  - `CommandDefinition` becomes `AgentDefinition`.
  - command registry becomes agent registry.
  - command prompt paths become agent prompt paths.
  - CLI surface changes from `factory command ...` and `factory run <command>` to agent-first naming.
  - Runtime result/event fields use agent terminology.
- **Excludes:** Department visualization and init templates.
- **Completion Signal:** Tests prove an agent can be created, resolved by alias, and executed through the runtime chain.

### WO-02: Introduce Configurable Departments

- **Outcome:** Departments exist as first-class project configuration and group agents without affecting runtime execution.
- **Scope:** Add department storage, types, create/update/list operations, and CLI commands.
- **Includes:**
  - Department definition with name, description, optional aliases, and ordered agent membership.
  - Validation that department membership references existing agents.
  - CLI operations to create departments and assign/unassign agents.
- **Excludes:** Mermaid rendering details.
- **Completion Signal:** Tests prove users can create arbitrary departments and assign agents without any built-in department names.

### WO-03: Replace Hard-Coded Visualization With Configured Structure

- **Outcome:** `factory visualize` renders only configured departments, agents, event bindings, sectors, findings, and their relationships.
- **Scope:** Remove dock/yard/shop-floor/finished-goods assumptions from the graph builder.
- **Includes:**
  - Department nodes generated from project config.
  - Agent nodes grouped by configured department membership.
  - Event-binding edges between agents.
  - Sector and finding references discovered from configured project state.
  - JSON and Mermaid output.
- **Excludes:** New runtime execution behavior.
- **Completion Signal:** Tests fail if visualization emits departments or workflow lanes that were not configured by the user or template.

### WO-04: Rebuild Workspace Initialization Around A Template

- **Outcome:** `factory init` creates a useful starter factory through normal configuration files, not hard-coded runtime rules.
- **Scope:** Update default workspace layout and seed config.
- **Includes:**
  - Initial departments representing dock, yard, shop-floor, and finished-goods as template data.
  - Initial starter agents only if they are represented as normal agent definitions.
  - Initial event bindings only if stored in normal event binding config.
  - No app logic that assumes those departments exist.
- **Excludes:** Historical migration from the old workspace.
- **Completion Signal:** A new initialized workspace can delete or rename all starter departments and still pass validation.

### WO-05: Align Events, Sectors, And Findings With Agents

- **Outcome:** Events, sectors, and findings continue to work after the command-to-agent pivot.
- **Scope:** Update event binding matchers, emitted event payloads, sector action references, and finding discovery to use agent terminology where execution identity matters.
- **Includes:**
  - Event binding `where` examples and tests use `agent`.
  - Runtime events include agent identity.
  - Sector action references can point to agents or sector-local actions.
  - Findings remain retrieval context, not executable behavior.
- **Excludes:** Department grouping behavior.
- **Completion Signal:** Tests prove event bindings can chain agents and sector/finding references can be visualized.

### WO-06: Refresh Tests And Final User Documentation

- **Outcome:** The repository describes the final configurable-agent model only.
- **Scope:** Update tests, README, architecture docs, and examples after implementation work is complete.
- **Includes:**
  - Replace command-centric examples with agent-centric examples.
  - Document departments as optional organization-only groups.
  - Document starter template behavior.
  - Remove obsolete pivot narrative.
- **Excludes:** Compatibility notes for the old model.
- **Completion Signal:** README and architecture docs describe the same model the CLI implements, and `go test ./...` passes.

## Known Ordering Or Dependency Constraints

1. WO-01 should happen before WO-05 because event payload terminology depends on the executable unit name.
2. WO-02 should happen before WO-03 because visualization must consume configured departments.
3. WO-04 should happen after WO-01 and WO-02 so the template seeds the final model.
4. WO-06 should happen last so documentation describes the final state rather than the transition.
5. WO-03 can proceed after the department config shape exists, even before the full init template is refreshed.

## Current Status

| Work Order | Status | Readiness |
| --- | --- | --- |
| WO-01 | candidate | ready |
| WO-02 | candidate | ready |
| WO-03 | candidate | waiting-for WO-02 |
| WO-04 | candidate | waiting-for WO-01 and WO-02 |
| WO-05 | candidate | waiting-for WO-01 |
| WO-06 | candidate | waiting-for implementation work |

## Completion Signal

This package is complete when:

- The CLI has no hard-coded workflow organization.
- Users can define agents, departments, event bindings, sectors, and findings through project state.
- The default dock/yard/shop-floor/finished-goods organization exists only as initial template data.
- `factory visualize` reflects configured state only.
- The final docs describe the agent and department model without legacy command terminology.

## Next Step

```text
Factory.call("process-work-package", {
  package_path: ".factory/03-shop-floor/d-work-packages/20260428T012116Z-configurable-agent-factory-pivot.dir"
})
```
