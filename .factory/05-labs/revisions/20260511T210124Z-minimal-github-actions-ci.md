# Minimal GitHub Actions CI Revision

## Target

- Branch: `factory/package/20260509T033121Z-triangulation-strategy`
- Worktree: `/home/igor/code/factory-cli`
- Pull request: `https://github.com/i-santos/factory/pull/1`

## Feedback

Auto-land is blocked because GitHub reports no PR checks. Add a minimal CI workflow so `land-pr` can require and observe a green check.

## Disposition

follow-up-order-on-active-branch

## Follow-Up Work

Add `.github/workflows/ci.yml` with a minimal Go test job for `pull_request` and pushes to `master`.
