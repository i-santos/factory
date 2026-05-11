# Gamified Visualization Foundation

## Status

shipping-pr-pending

## Pull Request

pending-until-created

## Outcome

Implemented the first local browser GUI foundation for Factory.

The CLI now serves a graph-backed Factory map from the same workspace model used by `visualize`. The GUI renders lanes, reusable automations, machines, circuits, sectors, relationships, and explicit run-command affordances without adding a separate frontend build stack.

## Product Decisions Preserved

- The first GUI slice is a zero-build local web UI served by the Go CLI.
- The GUI reads the same project model as the CLI.
- Execution remains explicit through CLI run commands.
- The orchestrator remains deterministic Go application code.

## Artifacts

- `internal/factory/gui.go`
- `internal/factory/gui_test.go`
- `internal/factory/visualize.go`
- `internal/factory/visualize_test.go`
- `internal/cli/root.go`
- `README.md`

## Behavior

- Adds `factory gui`.
- Serves a local HTML Factory map.
- Adds `/api/graph` for graph-backed UI data.
- Extends `BuildFactoryGraph` to include reusable circuits, machines, and automations.
- Displays command-copy affordances for machine and automation execution.
- Documents GUI usage.

## Source Work Order

`.factory/03-shop-floor/a-input-buffer/20260511T140238Z-create-gamified-visualization-foundation.dir/work-order.md`

## Source Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir/work-package.md`

## Branches

- Package branch: `factory/package/20260509T033121Z-triangulation-strategy`
- Work-order branch: `factory-order/20260511T140238Z-create-gamified-visualization-foundation`
- Worktree: `.worktrees/default/20260511T140238Z-create-gamified-visualization-foundation`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go build -o /tmp/factory-cli ./cmd/factory`

## Quality Control

Result: green.

The full Go test suite and build pass. Tests cover graph-backed GUI rendering, the HTTP handler, graph API, reusable circuit/machine/automation nodes, and run-command affordances.

## Audit

Tier: checklist.

- Scope stayed inside the GUI foundation, shared graph model, CLI command, tests, and docs.
- No frontend toolchain was introduced.
- The GUI does not hide execution behind agentic browser behavior.
- UI copy and command affordances are in English.
- The first GUI slice directly enables CLI/GUI alignment work.
