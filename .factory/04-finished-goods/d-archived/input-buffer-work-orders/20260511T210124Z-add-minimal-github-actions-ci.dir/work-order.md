# Add Minimal GitHub Actions CI

## Source

Revision follow-up for `https://github.com/i-santos/factory/pull/1`.

## Requested Outcome

The package PR has a minimal GitHub Actions CI workflow so `land-pr` can observe a green PR check before auto-landing.

## Scope

- Add `.github/workflows/ci.yml`.
- Run the workflow on `pull_request`.
- Run the workflow on pushes to `master`.
- Check out the repository.
- Set up Go from `go.mod`.
- Run `go test ./...`.

## Excludes

- Deployment automation.
- Release publishing.
- Branch protection configuration.
- Non-Go test matrix expansion.

## Completion Signal

- `.github/workflows/ci.yml` exists.
- Local verification passes with `go test ./...`.
- The current package PR branch is pushed so GitHub Actions can create the missing PR check.
