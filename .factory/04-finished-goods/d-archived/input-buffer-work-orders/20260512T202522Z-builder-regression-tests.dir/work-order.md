# Builder Regression Tests

## Status

integrated

## Package

`.factory/02-yard/b-discharged/20260512T201334Z-factory-builder-foundation.dir/work-package.md`

## Outcome

Added tests for builder persistence, invalid references, API behavior, GUI builder rendering, and command-vs-machine visualization semantics.

## Evidence

- `internal/factory/builder_test.go`
- `internal/factory/gui_test.go`
- `internal/factory/visualize_test.go`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
