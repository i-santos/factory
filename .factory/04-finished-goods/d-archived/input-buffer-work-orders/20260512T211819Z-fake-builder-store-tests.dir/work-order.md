# Fake BuilderStore Tests

## Status

integrated

## Package

`.factory/02-yard/b-discharged/20260512T211254Z-builder-persistence-abstraction.dir/work-package.md`

## Outcome

Added a fake BuilderStore test path proving the Builder API can use a non-filesystem store.

## Evidence

- `internal/factory/gui_test.go`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
