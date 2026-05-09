# Circuit Runtime Kernel Contract

This document defines the prompt-program contract loaded by every circuit. It is intentionally small: deterministic enough for orchestration, but pragmatic enough for agent execution.

## Runtime Boundary

A circuit is the only agent-executed unit in Factory. The application orchestrator invokes one circuit by running `codex exec` with:

- the circuit program
- the runtime kernel
- structured input
- execution metadata

The orchestrator is not an agent. It does not interpret prompt-specific meaning. It starts circuits, waits for completion, captures structured output, records state, and decides the next configured circuit or machine step.

## Program Shape

A circuit program may contain prose, but executable pseudocode must live in a `<code>` block. Guidance belongs in `<comment>`.

```text
<comment>
Explain intent, boundary, assumptions, and recovery rules.
</comment>

<code>
input = Runtime.input
Assert input.task exists
result = Pseudo.DoUsefulWork(input.task)
return { status: "succeeded", result }
</code>
```

## Tags

### `<code>`

Defines executable pseudocode. The runtime treats statements in this block as ordered instructions.

Rules:
- Execute top to bottom unless explicit control flow changes order.
- Unknown unprefixed operations are invalid.
- Return exactly one final result object.

### `<comment>`

Defines non-executable guidance.

Rules:
- Use it for goals, boundaries, rationale, and examples.
- Do not hide required execution steps in comments.
- Comments may guide `Pseudo.*` interpretation but cannot override code.

## Core Primitives

### `Runtime.input`

Structured input supplied by the orchestrator.

Rules:
- Read-only.
- Must be validated with `Assert` before required fields are used.
- Optional fields must declare defaults in the code.

### `Assert condition`

Deterministic guard.

Rules:
- If the condition is false, the circuit blocks.
- Assertions should be specific enough to produce useful recovery information.
- Assertions must not mutate state.

### `Tool.call(tool_ref, input)`

External tool invocation.

Rules:
- `tool_ref` must name an allowed script, CLI, or host tool.
- `input` must be structured.
- Tool results must be checked before use.
- Tool calls are side-effecting unless documented otherwise.

### `Trace.event(event)`

Durable execution trace.

Rules:
- Use for start, decisions, tool handoffs, blocking conditions, and finish.
- Events must include enough information to inspect execution later.
- Events must not contain secrets.

### `Pseudo.*`

Agent-inferred pseudocode operation.

Rules:
- `Pseudo.*` is allowed only inside the active circuit boundary.
- It must follow surrounding comments, assertions, and data contracts.
- It may not ignore explicit control flow.
- It must block when materially different interpretations are plausible.
- It may mutate state only when the active circuit owns that mutation.

## Determinism Model

Factory does not require every internal reasoning step to be deterministic. It requires deterministic boundaries:

- known input shape
- explicit assertions
- ordered execution
- structured tool calls
- structured return object
- traceable decisions
- block instead of guessing on ambiguous state

The orchestrator can rely on the final circuit output because every circuit must return a structured result.

## Result Contract

Every circuit returns one object:

```json
{
  "status": "succeeded|blocked|failed|waiting",
  "summary": "Short operator-facing summary.",
  "data": {},
  "artifacts": [],
  "events": [],
  "errors": [],
  "next": null
}
```

Rules:
- `status` is required.
- `summary` is required.
- `data` is circuit-owned structured output.
- `artifacts` contains file paths or durable references.
- `next` is advisory data only; automatic continuation belongs to the orchestrator configuration.

## Example: Transform Input

```text
<comment>
Normalize one incoming task into a short slug.
</comment>

<code>
input = Runtime.input
Assert input.title exists

slug = Pseudo.Slugify(input.title)

return {
  status: "succeeded",
  summary: "Task title normalized.",
  data: { slug },
  artifacts: [],
  events: [],
  errors: [],
  next: null
}
</code>
```

## Example: Tool-Backed Circuit

```text
<comment>
Run a validation script and return its result.
</comment>

<code>
input = Runtime.input
Assert input.workspace_root exists

Trace.event({
  eventType: "start",
  summary: "Validation circuit started."
})

validation = Tool.call("./scripts/validate-factory", {
  workspace_root: input.workspace_root
})

Assert validation.status == "passed"

return {
  status: "succeeded",
  summary: "Factory validation passed.",
  data: validation,
  artifacts: validation.artifacts optional default [],
  events: [],
  errors: [],
  next: null
}
</code>
```

## Blocking Rules

The circuit must return `blocked` when:

- required input is missing
- required tool is unavailable
- `Pseudo.*` interpretation is ambiguous
- a required assertion fails
- the circuit would need to mutate state outside its ownership boundary

## Validation Checklist

A circuit program is valid when:

- it has one `<code>` block for executable logic
- required inputs are asserted
- tool calls are explicit
- `Pseudo.*` operations are scoped and recoverable
- it returns one structured result
- comments do not contain hidden execution requirements
