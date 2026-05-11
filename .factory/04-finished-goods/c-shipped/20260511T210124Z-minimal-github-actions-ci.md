# Minimal GitHub Actions CI

## Status

shipping-pr-pending

## Pull Request

https://github.com/i-santos/factory/pull/1

## Outcome

Added a minimal GitHub Actions CI workflow for the current package release.

The workflow runs on pull requests and pushes to `master`, checks out the repository, sets up Go from `go.mod`, and runs `go test ./...`.

## Release Context

- Pull request: `https://github.com/i-santos/factory/pull/1`
- Package branch: `factory/package/20260509T033121Z-triangulation-strategy`
- Work order: `.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260511T210124Z-add-minimal-github-actions-ci.dir/work-order.md`

## Verification

- `env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...`

## Notes

This is a release-blocking follow-up for the shipped package PR. Auto-landing was blocked because GitHub reported no checks on the package branch.
