# Autonomous Visual Runs Package

## Status

shipping-pr-pending

## Pull Request

pending-until-created

## Outcome

Completed the first usable implementation slice for autonomous visual Factory runs.

The application can now represent durable visual run state, continue from supported deterministic next actions, stop visibly for unsafe continuation, persist run sessions, and expose run state plus decision requests through the local GUI/API.

## Integrated Work Orders

- WO-01: Added run-state, next-action, and decision-request models.
- WO-02: Implemented deterministic continuation rules.
- WO-03: Persisted and loaded run records.
- WO-04: Exposed run state through GUI/API.
- WO-05: Added autonomous visual run tests.

## Package Branch

`factory/package/20260512T135022Z-autonomous-visual-runs`

## Release-Ready Child Evidence

- `.factory/04-finished-goods/b-release-ready/20260512T141646Z-run-state-next-action-models.md`
- `.factory/04-finished-goods/b-release-ready/20260512T141647Z-deterministic-continuation-rules.md`
- `.factory/04-finished-goods/b-release-ready/20260512T141648Z-run-record-persistence.md`
- `.factory/04-finished-goods/b-release-ready/20260512T141649Z-run-state-gui-api.md`
- `.factory/04-finished-goods/b-release-ready/20260512T141650Z-autonomous-visual-run-tests.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`

## Next Step

Operator-controlled package shipping:

```text
<factory ship>
{ "release_ready_paths": [".factory/04-finished-goods/b-release-ready/20260512T141651Z-autonomous-visual-runs-package.md"] }
</factory>
```
