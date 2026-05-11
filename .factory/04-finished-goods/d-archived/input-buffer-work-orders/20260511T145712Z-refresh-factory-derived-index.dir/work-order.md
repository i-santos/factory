# Refresh Factory Derived Index

## Source Work Package

`.factory/02-yard/a-refining/20260511T142314Z-factory-workspace-organization-cleanup.dir`

## Requested Outcome

The derived Factory index reflects canonical records after completed lane reconciliation.

## Scope

- Rebuild `.factory/00-control-room/g-index/` from canonical Factory records.
- Verify the rebuilt index reports fresh.
- Record release-ready cleanup evidence.

## Excludes

- Changing canonical records outside derived index maintenance.
- Deciding older pivot package disposition.

## Completion Signal

- `factory-index.mjs status` reports `fresh`.
- Release-ready evidence records the rebuilt index result.
