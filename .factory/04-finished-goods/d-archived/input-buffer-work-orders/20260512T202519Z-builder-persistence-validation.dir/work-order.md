# Builder Persistence And Validation

## Status

integrated

## Package

`.factory/02-yard/b-discharged/20260512T201334Z-factory-builder-foundation.dir/work-package.md`

## Outcome

Added deterministic builder services for creating, reading, updating, and listing circuits, machines, automations, commands, and event bindings. Invalid references are rejected before persisted definitions are written.

## Evidence

- `internal/factory/builder.go`
- `internal/factory/builder_test.go`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
