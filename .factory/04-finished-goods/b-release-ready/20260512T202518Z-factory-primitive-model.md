# Factory Primitive Model

## Status

release-ready

## Outcome

Commands, circuits, machines, automations, and event bindings are now represented as distinct Factory primitives. Registered commands no longer masquerade as reusable machines in the graph or GUI run affordances.

## Package Branch

`factory/package/20260512T135022Z-autonomous-visual-runs`

## Work Order

`.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T202518Z-clarify-factory-primitive-model.dir/work-order.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
