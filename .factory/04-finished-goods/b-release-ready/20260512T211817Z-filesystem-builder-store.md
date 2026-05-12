# Filesystem BuilderStore

## Status

release-ready

## Outcome

Implemented current `.factory/` persistence as `FilesystemBuilderStore` while preserving compatibility wrappers and file layout.

## Package Branch

`factory/package/20260512T211254Z-builder-persistence-abstraction`

## Work Order

`.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T211817Z-filesystem-builder-store.dir/work-order.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
