# Builder API Endpoints

## Status

integrated

## Package

`.factory/02-yard/b-discharged/20260512T201334Z-factory-builder-foundation.dir/work-package.md`

## Outcome

Added `/api/builder` inventory and `/api/builder/{primitive}` endpoints for list, read, create, and update behavior across circuits, machines, automations, commands, and bindings.

## Evidence

- `internal/factory/gui.go`
- `internal/factory/gui_test.go`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
