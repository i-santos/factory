# Visual Builder UI

## Status

release-ready

## Outcome

The local GUI includes a first Factory Builder surface with HTML controls to create circuits, machines, automations, commands, and event bindings. Validation feedback is shown in the UI surface.

## Package Branch

`factory/package/20260512T135022Z-autonomous-visual-runs`

## Work Order

`.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T202521Z-visual-builder-ui.dir/work-order.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
