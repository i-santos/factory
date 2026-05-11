# Define Circuit Runtime Kernel Contract

## Source Work Package

`.factory/02-yard/a-refining/20260509T033121Z-gamified-factory-triangulation.dir`

## Shared Outcome

Factory becomes a gamified, Factorio-like application where users assemble circuits into machines and automations. Circuits are the agent runtime abstraction.

## Requested Outcome

Define the minimal deterministic `runtime-kernel.md` contract that every circuit loads before executing a prompt-program.

## Scope

- Document the circuit runtime boundary.
- Define `<code>` and `<comment>` semantics.
- Define required primitives: `Runtime.input`, `Assert`, `Tool.call(...)`, `Trace.event(...)`, and `Pseudo.*`.
- Define deterministic execution expectations and blocking behavior.
- Include example circuit programs.
- Keep orchestrator implementation out of scope.

## Ordering Constraints

- Depends on WO-01 triangulation baseline.
- Must complete before serial orchestrator implementation.

## Completion Signal

- The repository has a circuit runtime kernel contract.
- Prompt-program examples can be checked against the contract.
- The contract distinguishes deterministic primitives from `Pseudo.*` inference.
