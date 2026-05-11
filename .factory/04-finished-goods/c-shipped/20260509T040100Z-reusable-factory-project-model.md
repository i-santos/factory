# Reusable Factory Project Model

## Status

shipping-pr-pending

## Pull Request

https://github.com/i-santos/factory/pull/1

## Outcome

Defined the reusable local factory model shared by CLI and GUI.

## Artifacts

- `docs/project-model.md`
- `docs/architecture.md`
- `docs/triangulation.md`
- `README.md`

## Source Work Order

`.factory/03-shop-floor/a-input-buffer/20260509T035806Z-model-reusable-factory-pieces.dir/work-order.md`

## Source Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir/work-package.md`

## Branches

- Package branch: `factory/package/20260509T033121Z-triangulation-strategy`
- Work-order branch: `factory/order/20260509T035806Z-model-reusable-factory-pieces`
- Merge commit: `b43306f492b18bf4be56d9e48ec9e483ae58c934`

## Verification

- `GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
