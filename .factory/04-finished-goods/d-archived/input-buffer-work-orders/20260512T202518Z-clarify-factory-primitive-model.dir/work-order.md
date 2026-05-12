# Clarify Factory Primitive Model

## Status

integrated

## Package

`.factory/02-yard/b-discharged/20260512T201334Z-factory-builder-foundation.dir/work-package.md`

## Outcome

Factory now distinguishes registered commands from reusable machines in the model, graph, GUI run affordances, and project context. The old `evolve-factory-cli` seed command is archived and removed from the active command registry.

## Evidence

- `.factory/00-control-room/a-config/commands.json`
- `.factory/00-control-room/b-commands/archive/evolve-factory-cli.md`
- `docs/project-model.md`
- `docs/architecture.md`
- `README.md`
- `internal/factory/visualize.go`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
