# Visual Builder UI

## Status

integrated

## Package

`.factory/02-yard/b-discharged/20260512T201334Z-factory-builder-foundation.dir/work-package.md`

## Outcome

Added the first Factory Builder panel to the GUI with structured HTML controls for creating circuits, machines, automations, commands, and event bindings. Registered commands now display as commands and copy `factory run <command>` instead of invalid machine commands.

## Evidence

- `internal/factory/gui.go`
- `internal/factory/visualize.go`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
