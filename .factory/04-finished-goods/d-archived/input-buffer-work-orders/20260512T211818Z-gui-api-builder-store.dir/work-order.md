# GUI API BuilderStore Injection

## Status

integrated

## Package

`.factory/02-yard/b-discharged/20260512T211254Z-builder-persistence-abstraction.dir/work-package.md`

## Outcome

Updated GUI/API builder handlers to use `BuilderStore`, with filesystem store construction as the default and explicit store injection through `GUIOptions`.

## Evidence

- `internal/factory/gui.go`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`
- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go run ./cmd/factory validate`
