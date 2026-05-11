# Initialize New Factories From An Init Prompt

## Status

release-ready

## Outcome

Implemented scaffold-only new factory initialization with explicit starter automation execution.

`factory init` now creates an inspectable latest-schema workspace without invoking Codex. The initialized workspace includes starter runtime/config metadata, an init circuit, an init machine, an init automation, migration evidence directories, and an empty triangulation directory for explicit later generation.

## Product Decisions Preserved

- Default initialization scaffolds only.
- Starter triangulation is produced only when the user explicitly runs the generated init automation.
- Workspace updates use one canonical active workspace root with schema metadata and migration evidence, not permanent active `.factory/vN` roots.
- The orchestrator remains deterministic Go application code; agentic behavior stays inside circuits.

## Artifacts

- `internal/factory/workspace.go`
- `internal/factory/types.go`
- `internal/factory/workspace_test.go`
- `internal/cli/root.go`
- `README.md`

## Behavior

- Adds latest workspace/runtime constants.
- Adds migration metadata to workspace config.
- Adds `factory update` command.
- Scaffolds:
  - `.factory/runtime/runtime-kernel.md`
  - `.factory/circuits/init/circuit.json`
  - `.factory/circuits/init/program.md`
  - `.factory/machines/init.json`
  - `.factory/automations/init.json`
  - `.factory/triangulation/`
  - `.factory/00-control-room/d-migrations/`
- Keeps `factory init` idempotent and scaffold-only.
- Supports custom workspace roots in scaffolded circuit resource references.
- Records update evidence in `.factory/00-control-room/d-migrations/applied.jsonl`.

## Source Work Order

`.factory/03-shop-floor/a-input-buffer/20260511T131931Z-initialize-new-factories-from-init-prompt.dir/work-order.md`

## Source Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir/work-package.md`

## Branches

- Package branch: `factory/package/20260509T033121Z-triangulation-strategy`
- Work-order branch: `factory-order/20260511T131931Z-initialize-new-factories-from-init-prompt`
- Worktree: `.worktrees/default/20260511T132321Z-initialize-new-factories-from-init-prompt`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go build -o /tmp/factory-cli ./cmd/factory`
- `/tmp/factory-cli --project-root /tmp/factory-cli-init-check-fresh init`
- `/tmp/factory-cli --project-root /tmp/factory-cli-init-check-custom --workspace-root .custom-factory init`
- `FACTORY_CIRCUIT_MOCK_RESULT='{"status":"succeeded","summary":"mock init","data":{"triangulation":"created"},"artifacts":[".factory/triangulation/roadmap.md"]}' /tmp/factory-cli --project-root /tmp/factory-cli-init-check-fresh run automation init`

## Quality Control

Result: green.

The full Go test suite and build pass. The smoke checks confirm default init scaffolds without producing triangulation files and that the generated init automation can be dispatched explicitly through the serial orchestrator with an injectable circuit result.

## Audit

Tier: checklist.

- Scope stayed inside workspace initialization, CLI command registration, documentation, and focused tests.
- No agentic behavior was added to the orchestrator.
- Default `factory init` does not run Codex.
- Custom workspace root resource references were verified.
- Update behavior is idempotent for the current schema and records migration evidence.

## Next Step

Operator-controlled shipping remains explicit:

```text
<factory ship>
{ "release_ready_paths": [".factory/04-finished-goods/b-release-ready/20260511T133000Z-initialize-new-factories-from-init-prompt.md"] }
</factory>
```
