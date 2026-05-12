# Builder Persistence And Validation

## Status

release-ready

## Outcome

Factory has deterministic builder service functions for persisted circuits, machines, automations, commands, and event bindings, including reference validation before write.

## Package Branch

`factory/package/20260512T135022Z-autonomous-visual-runs`

## Work Order

`.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T202519Z-builder-persistence-validation.dir/work-order.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
