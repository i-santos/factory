# Refresh Factory Derived Index

## Status

release-ready

## Outcome

Rebuilt the derived Factory index projection after lane reconciliation.

## Index Result

- Status: fresh
- Record count: 25
- Warning count: 0
- Generated at: 2026-05-11T14:58:42.853Z

## Source Work Order

`.factory/03-shop-floor/a-input-buffer/20260511T145712Z-refresh-factory-derived-index.dir/work-order.md`

## Verification

- `node .factory/00-control-room/d-scripts/factory-index.mjs rebuild --workspace-root .factory --index-dir .factory/00-control-room/g-index`
- `node .factory/00-control-room/d-scripts/factory-index.mjs status --workspace-root .factory --index-dir .factory/00-control-room/g-index`

## Notes

The index is derived state under `.factory/00-control-room/g-index/` and remains rebuildable from canonical Factory records.
