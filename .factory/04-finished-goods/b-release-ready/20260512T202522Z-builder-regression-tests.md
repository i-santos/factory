# Builder Regression Tests

## Status

release-ready

## Outcome

Regression tests cover builder persistence, invalid references, API creation behavior, GUI builder rendering, and command-vs-machine run affordances.

## Package Branch

`factory/package/20260512T135022Z-autonomous-visual-runs`

## Work Order

`.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T202522Z-builder-regression-tests.dir/work-order.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
