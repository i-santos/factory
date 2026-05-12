# Filesystem BuilderStore

## Status

integrated

## Package

`.factory/02-yard/b-discharged/20260512T211254Z-builder-persistence-abstraction.dir/work-package.md`

## Outcome

Moved current `.factory/` filesystem persistence behind `FilesystemBuilderStore` while preserving wrapper functions and the existing on-disk layout.

## Evidence

- `internal/factory/builder.go`
- `internal/factory/builder_test.go`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
