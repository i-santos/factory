# BuilderStore GUI API Integration

## Status

release-ready

## Outcome

Updated Builder GUI/API endpoints to use `BuilderStore`, with filesystem persistence as the default and test/future implementations injectable through `GUIOptions`.

## Package Branch

`factory/package/20260512T211254Z-builder-persistence-abstraction`

## Work Order

`.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T211818Z-gui-api-builder-store.dir/work-order.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
