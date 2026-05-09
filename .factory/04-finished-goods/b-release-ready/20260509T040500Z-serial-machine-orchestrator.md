# Serial Machine Orchestrator

## Status

release-ready

## Outcome

Implemented the first serial orchestrator slice for reusable machines and automations.

## Artifacts

- `internal/factory/orchestrator.go`
- `internal/factory/orchestrator_test.go`
- `internal/cli/root.go`

## Behavior

- Loads circuits from `.factory/circuits/<id>/circuit.json`.
- Loads machines from `.factory/machines/<id>.json`.
- Loads automations from `.factory/automations/<id>.json`.
- Runs machine circuits serially.
- Runs automation machines serially.
- Provides injectable circuit runner for tests.
- Adds CLI subcommands:
  - `factory run machine <name>`
  - `factory run automation <name>`

## Source Work Order

`.factory/03-shop-floor/a-input-buffer/20260509T040012Z-build-serial-machine-orchestrator.dir/work-order.md`

## Source Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir/work-package.md`

## Branches

- Package branch: `factory/package/20260509T033121Z-triangulation-strategy`
- Work-order branch: `factory/order/20260509T040012Z-build-serial-machine-orchestrator`
- Merge commit: `a7bd0e8d40cbc73ad403ec020ec56424fd9d7237`

## Verification

- `GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
