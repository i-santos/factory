# Builder Persistence Abstraction

## Source Intake

The operator asked for an implementation work package to introduce pluggable Factory Builder persistence.

The current Builder implementation writes directly to `.factory/` files from concrete functions in `internal/factory/builder.go`, and the GUI/API handlers call those concrete functions directly. That was acceptable for the first Builder slice, but the Builder needs a strategy/repository-style boundary before more functionality is added.

## Shared Outcome

Factory Builder persistence is behind a clean Go abstraction. The current filesystem-backed `.factory/` layout remains the default implementation, while the GUI/API layer depends on the abstraction so future SQLite, Postgres, or remote-service implementations can be added without rewriting request handlers.

## Why This Is A Work Package

This is larger than one isolated edit because it coordinates several implementation surfaces:

- persistence interface design
- filesystem store extraction
- validation placement
- GUI/API dependency inversion
- compatibility with current workspace layout
- regression tests for filesystem and fake/in-memory store behavior

Those areas should ship together as one cohesive package, but can be refined into focused work orders.

## Intended Work Orders

### WO-01: Define BuilderStore persistence boundary

Introduce a `BuilderStore` interface or equivalent abstraction for circuits, machines, automations, commands, and event bindings.

Acceptance criteria:
- The interface supports list/inventory behavior.
- The interface supports read where the API exposes read behavior.
- The interface supports create/update for circuits, machines, automations, commands, and event bindings.
- Interface naming and method signatures leave room for future SQLite, Postgres, or remote implementations.

### WO-02: Extract filesystem BuilderStore implementation

Move current direct filesystem behavior behind a `FilesystemBuilderStore` implementation.

Acceptance criteria:
- Existing `.factory/` persistence layout is unchanged.
- Circuit creation still writes `circuit.json` and `program.md`.
- Machine, automation, command, and binding persistence remain compatible with existing files.
- Validation remains deterministic and rejects invalid references before writing invalid definitions.

### WO-03: Update GUI/API to use BuilderStore

Update the GUI/API builder handlers so they depend on the store abstraction instead of calling direct filesystem functions.

Acceptance criteria:
- `NewGUIHandler` or the GUI options can receive or construct a builder store.
- Builder endpoints use the abstraction for list, read, create, and update.
- Public API behavior remains compatible with the current web UI.
- Error responses remain clear for validation failures.

### WO-04: Add fake or in-memory BuilderStore tests

Add a fake or in-memory store test path that proves the API/service layer is decoupled from filesystem persistence.

Acceptance criteria:
- At least one test drives builder API behavior through a non-filesystem store.
- Existing filesystem builder tests continue to pass.
- Tests prove the GUI/API layer can use a substituted store.

### WO-05: Documentation and verification

Document the persistence boundary and verify the full package.

Acceptance criteria:
- Project docs or README explain that filesystem persistence is one `BuilderStore` implementation.
- Future persistence implementations are described as an intended extension point.
- `go test ./...` passes.

## Known Ordering Or Dependency Constraints

1. Define the interface first.
2. Extract the current filesystem behavior behind that interface.
3. Update GUI/API construction and handlers to depend on the abstraction.
4. Add fake/in-memory tests after the API has a store injection path.
5. Update docs and run full verification last.

## Completion Signal

The package is complete when the Builder API and GUI no longer directly depend on concrete filesystem helper functions, while all current filesystem behavior and tests remain compatible.

## Lifecycle Status

release-ready

## Produced Work Orders

- `.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T211816Z-define-builder-store-boundary.dir/work-order.md`
- `.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T211817Z-filesystem-builder-store.dir/work-order.md`
- `.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T211818Z-gui-api-builder-store.dir/work-order.md`
- `.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T211819Z-fake-builder-store-tests.dir/work-order.md`
- `.factory/04-finished-goods/d-archived/input-buffer-work-orders/20260512T211820Z-builder-store-docs-verification.dir/work-order.md`

## Release-Ready Evidence

- `.factory/04-finished-goods/b-release-ready/20260512T211816Z-builder-store-boundary.md`
- `.factory/04-finished-goods/b-release-ready/20260512T211817Z-filesystem-builder-store.md`
- `.factory/04-finished-goods/b-release-ready/20260512T211818Z-builder-store-gui-api.md`
- `.factory/04-finished-goods/b-release-ready/20260512T211819Z-fake-builder-store-tests.md`
- `.factory/04-finished-goods/b-release-ready/20260512T211820Z-builder-store-docs.md`
- `.factory/04-finished-goods/b-release-ready/20260512T211821Z-builder-persistence-abstraction-package.md`

## Current Status

- WO-01: integrated
- WO-02: integrated
- WO-03: integrated
- WO-04: integrated
- WO-05: integrated
