# CLI And GUI Project Model Alignment

## Status

shipping-pr-pending

## Pull Request

pending-until-created

## Outcome

Aligned CLI and GUI inspection around one shared Factory project model.

The CLI now exposes validation for reusable circuits, machines, automations, sectors, and cross-references. Graph-backed CLI visualization and GUI rendering use the same validated workspace data, so invalid project model references fail before inspection surfaces diverge.

## Artifacts

- `internal/factory/validate.go`
- `internal/factory/validate_test.go`
- `internal/factory/visualize.go`
- `internal/cli/root.go`
- `README.md`

## Behavior

- Adds `factory validate`.
- Validates circuit program paths.
- Validates machine-to-circuit references.
- Validates automation-to-machine references.
- Reports shared model inventory.
- Gates graph building on shared validation.
- Documents validation usage.

## Source Work Order

`.factory/03-shop-floor/a-input-buffer/20260511T140842Z-align-cli-and-gui-project-model.dir/work-order.md`

## Source Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir/work-package.md`

## Branches

- Package branch: `factory/package/20260509T033121Z-triangulation-strategy`
- Work-order branch: `factory-order/20260511T140842Z-align-cli-and-gui-project-model`
- Worktree: `.worktrees/default/20260511T140842Z-align-cli-and-gui-project-model`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go build -o /tmp/factory-cli ./cmd/factory`
- `/tmp/factory-cli --project-root /tmp/factory-cli-align-check init`
- `/tmp/factory-cli --project-root /tmp/factory-cli-align-check validate`
- `/tmp/factory-cli --project-root /tmp/factory-cli-align-check visualize --format json`

## Quality Control

Result: green.

The full Go test suite and build pass. Smoke checks prove an initialized CLI workspace validates and renders the same reusable model through graph JSON.

## Audit

Tier: checklist.

- No GUI-only schema was introduced.
- Validation uses the same loaders as orchestration and graph rendering.
- Invalid references block graph construction with a clear validation failure.
- The work preserves one canonical active workspace root.
