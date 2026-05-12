# evolve-factory-cli

---
factory_runtime_contract: 1
factory_prompt_contract: project-command.v1
---

<comment>
Goal: evolve Factory itself by using the CLI-owned command registry and prompt materialization flow.

Boundary:
- New executable Factory behavior belongs in `.factory/00-control-room/b-commands/` and `commands.json`.
- This command should produce small, verifiable CLI improvements.
- The CLI must remain workflow-agnostic: users create commands and event bindings instead of changing Go code for each workflow.
</comment>

<code>
request = Runtime.input optional

target = request.target optional default "factory-cli"
slice = request.slice optional default "next-small-cli-improvement"

InspectFactoryCliState({
  command_registry: ".factory/00-control-room/a-config/commands.json",
  event_bindings: ".factory/00-control-room/a-config/event-bindings.json",
  cli_source: ["cmd/factory", "internal/cli", "internal/factory"]
})

ChooseOneConcreteImprovement({
  target,
  slice,
  constraints: [
    "do not add workflow-specific hard-coded commands",
    "preserve command prompt files as source of truth",
    "preserve structured registry and event binding contracts",
    "verify with Go tests"
  ]
})

ImplementAndVerifyImprovement()

return {
  status: "succeeded",
  resultType: "factory-cli.evolved",
  summary: "Factory CLI evolution slice completed.",
  artifacts: [
    ".factory/00-control-room/b-commands/evolve-factory-cli.md"
  ],
  events: [],
  data: {}
}
</code>
