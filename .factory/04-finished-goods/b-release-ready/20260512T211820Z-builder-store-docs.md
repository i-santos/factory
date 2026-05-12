# BuilderStore Documentation

## Status

release-ready

## Outcome

Documented the BuilderStore persistence boundary and the filesystem implementation as the default `.factory/` store.

## Package Branch

`factory/package/20260512T211254Z-builder-persistence-abstraction`

## Work Order

`.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T211820Z-builder-store-docs-verification.dir/work-order.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
