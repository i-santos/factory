# Factory Workspace Organization Cleanup

## Source

Created by `factory organize` from audit report:

`.factory/05-labs/organization-audits/20260511T142314Z-organization-audit.md`

## Shared Outcome

The Factory workspace lanes accurately reflect current work state after the completed gamified Factory package, without deleting or rewriting historical records inline.

## Why This Is A Work Package

The cleanup spans multiple lane types and decisions:

- completed package and work-order records
- stale derived index projections
- an older package whose product direction may be superseded

These should be handled as explicit, reviewable cleanup work rather than as inline mutations during organization audit.

## Intended Work Orders

### WO-01: Reconcile Completed Package And Work-Order Lane Placement

- **Outcome:** Completed package and integrated work-order records no longer appear as active input-buffer/refining work.
- **Scope:** Move, archive, or annotate completed Factory records according to current lane policy.
- **Includes:**
  - gamified Factory triangulation package record
  - WO-01 through WO-07 input-buffer work-order records
  - release-ready package evidence references
- **Excludes:** Shipping or publishing release-ready work.
- **Completion Signal:** Active lanes contain only active or intentionally pending work.

### WO-02: Refresh Derived Factory Index Projection

- **Outcome:** The derived index reflects canonical records after cleanup.
- **Scope:** Rebuild or verify `.factory/00-control-room/g-index/`.
- **Includes:**
  - index manifest freshness
  - record count and relationship projections
  - release ledger consistency
- **Excludes:** Changing canonical records outside index maintenance.
- **Completion Signal:** `factory index` reports a fresh projection with no stale warnings.

### WO-03: Decide Older Pivot Package Disposition

- **Outcome:** The configurable-agent pivot package has an explicit disposition under the current gamified Factory direction.
- **Scope:** Review `.factory/03-shop-floor/d-work-packages/20260428T012116Z-configurable-agent-factory-pivot.dir/work-package.md`.
- **Includes:**
  - classify as superseded, archived, or convertible future work
  - preserve any useful ideas as future work if still relevant
  - avoid silently deleting historical context
- **Excludes:** Implementing a new product pivot.
- **Completion Signal:** The old pivot package no longer appears as ambiguous active/stale work.

## Known Ordering Or Dependency Constraints

1. WO-01 should happen before WO-02 so the derived index reflects final canonical lane placement.
2. WO-03 can happen independently, but should be resolved before considering the workspace fully organized.

## Current Status

| Work Order | Status | Readiness |
| --- | --- | --- |
| WO-01 | integrated | done |
| WO-02 | candidate | ready |
| WO-03 | candidate | ready |

## Produced Work Orders

- WO-01: `.factory/03-shop-floor/a-input-buffer/20260511T145301Z-reconcile-completed-factory-lanes.dir/work-order.md`

## Integrated Work Orders

- WO-01: Reconciled completed package and work-order lane placement.
  - Release-ready evidence: `.factory/04-finished-goods/b-release-ready/20260511T145600Z-reconcile-completed-factory-lanes.md`
  - Moved completed gamified Factory package to `.factory/02-yard/b-discharged/20260509T033121Z-gamified-factory-triangulation.dir/`.
  - Moved integrated source work orders to `.factory/04-finished-goods/d-archived/input-buffer-work-orders/`.

## Completion Signal

- Active Factory lanes no longer show completed package work as pending.
- The derived index is fresh.
- The older configurable-agent pivot package has an explicit disposition.
- Cleanup is recorded as reviewable Factory work rather than untracked manual mutation.
