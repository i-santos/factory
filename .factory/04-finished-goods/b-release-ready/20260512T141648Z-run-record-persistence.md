# Run Record Persistence

## Status

release-ready

## Outcome

Persisted command-chain run sessions under `.factory/runs/` and added helpers to load recorded runs.

## Package Branch

`factory/package/20260512T135022Z-autonomous-visual-runs`

## Work Order

`.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T141648Z-persist-load-run-records.dir/work-order.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
