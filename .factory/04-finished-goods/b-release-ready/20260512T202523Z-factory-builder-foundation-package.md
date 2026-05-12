# Factory Builder Foundation Package

## Status

release-ready

## Outcome

Completed the first Factory Builder foundation slice.

Factory now has a deterministic builder service, API endpoints, and a local GUI builder surface for creating concrete Factory primitives. Commands are distinct from machines in the graph and GUI, and the old `evolve-factory-cli` seed command is archived outside the active command registry.

## Integrated Work Orders

- WO-01: Clarified Factory primitive model and default seed state.
- WO-02: Added builder persistence and validation services.
- WO-03: Exposed Factory Builder API endpoints.
- WO-04: Built the first visual Factory Builder UI.
- WO-05: Added builder regression tests.

## Package Branch

`factory/package/20260512T135022Z-autonomous-visual-runs`

## Release-Ready Child Evidence

- `.factory/04-finished-goods/b-release-ready/20260512T202518Z-factory-primitive-model.md`
- `.factory/04-finished-goods/b-release-ready/20260512T202519Z-builder-persistence-validation.md`
- `.factory/04-finished-goods/b-release-ready/20260512T202520Z-builder-api-endpoints.md`
- `.factory/04-finished-goods/b-release-ready/20260512T202521Z-visual-builder-ui.md`
- `.factory/04-finished-goods/b-release-ready/20260512T202522Z-builder-regression-tests.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`

## Next Step

Operator-controlled package shipping:

```text
<factory ship>
{ "release_ready_paths": [".factory/04-finished-goods/b-release-ready/20260512T202523Z-factory-builder-foundation-package.md"] }
</factory>
```
