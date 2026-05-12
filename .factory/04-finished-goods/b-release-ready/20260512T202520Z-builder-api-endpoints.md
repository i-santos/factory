# Builder API Endpoints

## Status

release-ready

## Outcome

The GUI server exposes structured builder endpoints for inventory, list, read, create, and update behavior for Factory primitives.

## Package Branch

`factory/package/20260512T135022Z-autonomous-visual-runs`

## Work Order

`.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T202520Z-builder-api-endpoints.dir/work-order.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
