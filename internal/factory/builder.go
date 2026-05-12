package factory

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var builderIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

type BuilderOptions struct {
	ProjectRoot   string
	WorkspaceRoot string
}

type BuilderInventory struct {
	Circuits    []CircuitDefinition    `json:"circuits"`
	Machines    []MachineDefinition    `json:"machines"`
	Automations []AutomationDefinition `json:"automations"`
	Commands    []CommandDefinition    `json:"commands"`
	Bindings    []EventBinding         `json:"bindings"`
}

type BuilderCircuitInput struct {
	ID            string `json:"id"`
	Name          string `json:"name,omitempty"`
	Program       string `json:"program,omitempty"`
	ProgramBody   string `json:"programBody,omitempty"`
	RuntimeKernel string `json:"runtimeKernel,omitempty"`
	Description   string `json:"description,omitempty"`
}

type BuilderCommandInput struct {
	CommandDefinition
	PromptBody string `json:"promptBody,omitempty"`
}

func ListBuilderInventory(opts BuilderOptions) (BuilderInventory, error) {
	opts = normalizeBuilderOptions(opts)
	cfg, err := LoadConfig(opts.ProjectRoot, opts.WorkspaceRoot)
	if err != nil {
		return BuilderInventory{}, err
	}
	circuits, err := loadAllCircuits(opts.ProjectRoot, cfg.RootDir)
	if err != nil {
		return BuilderInventory{}, err
	}
	machines, err := loadAllMachines(opts.ProjectRoot, cfg.RootDir)
	if err != nil {
		return BuilderInventory{}, err
	}
	automations, err := loadAllAutomations(opts.ProjectRoot, cfg.RootDir)
	if err != nil {
		return BuilderInventory{}, err
	}
	reg, err := LoadRegistry(opts.ProjectRoot, cfg)
	if err != nil {
		return BuilderInventory{}, err
	}
	bindings, err := LoadEventBindings(opts.ProjectRoot, cfg)
	if err != nil {
		return BuilderInventory{}, err
	}
	return BuilderInventory{
		Circuits:    sortedValues(circuits),
		Machines:    sortedValues(machines),
		Automations: sortedValues(automations),
		Commands:    sortedValues(reg.Commands),
		Bindings:    bindings.Bindings,
	}, nil
}

func CreateOrUpdateCircuit(opts BuilderOptions, input BuilderCircuitInput) (CircuitDefinition, error) {
	opts = normalizeBuilderOptions(opts)
	if err := validateBuilderID(input.ID, "circuit"); err != nil {
		return CircuitDefinition{}, err
	}
	cfg, err := LoadConfig(opts.ProjectRoot, opts.WorkspaceRoot)
	if err != nil {
		return CircuitDefinition{}, err
	}
	def := CircuitDefinition{
		ID:            input.ID,
		Name:          strings.TrimSpace(input.Name),
		Program:       strings.TrimSpace(input.Program),
		RuntimeKernel: strings.TrimSpace(input.RuntimeKernel),
		Description:   strings.TrimSpace(input.Description),
	}
	if def.Program == "" {
		def.Program = "project://" + filepath.ToSlash(filepath.Join(cfg.RootDir, "circuits", def.ID, "program.md"))
	}
	if def.RuntimeKernel == "" {
		def.RuntimeKernel = "project://" + filepath.ToSlash(filepath.Join(cfg.RootDir, "runtime", "runtime-kernel.md"))
	}
	if input.ProgramBody != "" {
		if err := writeProjectResource(opts.ProjectRoot, def.Program, input.ProgramBody); err != nil {
			return CircuitDefinition{}, err
		}
	} else if path, ok := projectResourcePath(opts.ProjectRoot, def.Program); ok {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := writeProjectResource(opts.ProjectRoot, def.Program, defaultBuilderCircuitProgram(def.ID)); err != nil {
				return CircuitDefinition{}, err
			}
		} else if err != nil {
			return CircuitDefinition{}, err
		}
	}
	if err := writeJSON(filepath.Join(opts.ProjectRoot, cfg.RootDir, "circuits", def.ID, "circuit.json"), def); err != nil {
		return CircuitDefinition{}, err
	}
	return def, nil
}

func CreateOrUpdateMachine(opts BuilderOptions, def MachineDefinition) (MachineDefinition, error) {
	opts = normalizeBuilderOptions(opts)
	if err := validateBuilderID(def.ID, "machine"); err != nil {
		return MachineDefinition{}, err
	}
	if len(def.Circuits) == 0 {
		return MachineDefinition{}, fmt.Errorf("machine circuits are required")
	}
	cfg, err := LoadConfig(opts.ProjectRoot, opts.WorkspaceRoot)
	if err != nil {
		return MachineDefinition{}, err
	}
	circuits, err := loadAllCircuits(opts.ProjectRoot, cfg.RootDir)
	if err != nil {
		return MachineDefinition{}, err
	}
	for _, step := range def.Circuits {
		if strings.TrimSpace(step.Circuit) == "" {
			return MachineDefinition{}, fmt.Errorf("machine circuit reference is required")
		}
		if _, ok := circuits[step.Circuit]; !ok {
			return MachineDefinition{}, fmt.Errorf("machine references missing circuit %q", step.Circuit)
		}
	}
	if err := writeJSON(filepath.Join(opts.ProjectRoot, cfg.RootDir, "machines", def.ID+".json"), def); err != nil {
		return MachineDefinition{}, err
	}
	return def, nil
}

func CreateOrUpdateAutomation(opts BuilderOptions, def AutomationDefinition) (AutomationDefinition, error) {
	opts = normalizeBuilderOptions(opts)
	if err := validateBuilderID(def.ID, "automation"); err != nil {
		return AutomationDefinition{}, err
	}
	if len(def.Machines) == 0 {
		return AutomationDefinition{}, fmt.Errorf("automation machines are required")
	}
	cfg, err := LoadConfig(opts.ProjectRoot, opts.WorkspaceRoot)
	if err != nil {
		return AutomationDefinition{}, err
	}
	machines, err := loadAllMachines(opts.ProjectRoot, cfg.RootDir)
	if err != nil {
		return AutomationDefinition{}, err
	}
	for _, step := range def.Machines {
		if strings.TrimSpace(step.Machine) == "" {
			return AutomationDefinition{}, fmt.Errorf("automation machine reference is required")
		}
		if _, ok := machines[step.Machine]; !ok {
			return AutomationDefinition{}, fmt.Errorf("automation references missing machine %q", step.Machine)
		}
	}
	if err := writeJSON(filepath.Join(opts.ProjectRoot, cfg.RootDir, "automations", def.ID+".json"), def); err != nil {
		return AutomationDefinition{}, err
	}
	return def, nil
}

func CreateOrUpdateCommand(opts BuilderOptions, input BuilderCommandInput) (CommandDefinition, error) {
	opts = normalizeBuilderOptions(opts)
	def := input.CommandDefinition
	if err := validateBuilderID(def.Name, "command"); err != nil {
		return CommandDefinition{}, err
	}
	cfg, err := LoadConfig(opts.ProjectRoot, opts.WorkspaceRoot)
	if err != nil {
		return CommandDefinition{}, err
	}
	if def.ResultSchema == "" {
		def.ResultSchema = ResultContractV1
	}
	if len(def.Emits) == 0 {
		def.Emits = []string{"command.completed", "command.blocked", "command.failed", "command.waiting", "command.needs-human-action"}
	}
	if len(def.Consumes) == 0 {
		def.Consumes = []string{"operator.requested"}
	}
	if def.Prompt == "" {
		def.Prompt = "project://" + filepath.ToSlash(filepath.Join(cfg.Commands.ProjectCommandsDir, def.Name+".md"))
	}
	reg, err := LoadRegistry(opts.ProjectRoot, cfg)
	if err != nil {
		return CommandDefinition{}, err
	}
	for name, existing := range reg.Commands {
		if name == def.Name {
			continue
		}
		for _, alias := range def.Aliases {
			if alias == name || contains(existing.Aliases, alias) {
				return CommandDefinition{}, fmt.Errorf("alias %q conflicts with command %q", alias, name)
			}
		}
	}
	if input.PromptBody != "" {
		if err := writeProjectResource(opts.ProjectRoot, def.Prompt, input.PromptBody); err != nil {
			return CommandDefinition{}, err
		}
	} else if err := writePromptIfMissing(opts.ProjectRoot, def.Prompt, defaultBuilderCommandPrompt(def.Name)); err != nil {
		return CommandDefinition{}, err
	}
	reg.Commands[def.Name] = def
	return def, SaveRegistry(opts.ProjectRoot, cfg, reg)
}

func CreateEventBinding(opts BuilderOptions, binding EventBinding) (EventBinding, error) {
	opts = normalizeBuilderOptions(opts)
	cfg, err := LoadConfig(opts.ProjectRoot, opts.WorkspaceRoot)
	if err != nil {
		return EventBinding{}, err
	}
	reg, err := LoadRegistry(opts.ProjectRoot, cfg)
	if err != nil {
		return EventBinding{}, err
	}
	if _, ok := ResolveCommand(reg, binding.Run); !ok {
		return EventBinding{}, fmt.Errorf("binding target command %q is not registered", binding.Run)
	}
	if command := binding.Where["command"]; command != "" {
		if _, ok := ResolveCommand(reg, command); !ok {
			return EventBinding{}, fmt.Errorf("binding source command %q is not registered", command)
		}
	}
	return AddEventBinding(opts.ProjectRoot, cfg, binding)
}

func UpdateEventBinding(opts BuilderOptions, index int, binding EventBinding) (EventBinding, error) {
	opts = normalizeBuilderOptions(opts)
	cfg, err := LoadConfig(opts.ProjectRoot, opts.WorkspaceRoot)
	if err != nil {
		return EventBinding{}, err
	}
	reg, err := LoadRegistry(opts.ProjectRoot, cfg)
	if err != nil {
		return EventBinding{}, err
	}
	if _, ok := ResolveCommand(reg, binding.Run); !ok {
		return EventBinding{}, fmt.Errorf("binding target command %q is not registered", binding.Run)
	}
	if command := binding.Where["command"]; command != "" {
		if _, ok := ResolveCommand(reg, command); !ok {
			return EventBinding{}, fmt.Errorf("binding source command %q is not registered", command)
		}
	}
	if binding.On == "" {
		return EventBinding{}, fmt.Errorf("binding event type is required")
	}
	if binding.Mode == "" {
		binding.Mode = "auto"
	}
	file, err := LoadEventBindings(opts.ProjectRoot, cfg)
	if err != nil {
		return EventBinding{}, err
	}
	if index < 0 || index >= len(file.Bindings) {
		return EventBinding{}, os.ErrNotExist
	}
	file.Bindings[index] = binding
	return binding, SaveEventBindings(opts.ProjectRoot, cfg, file)
}

func ReadBuilderPrimitive(opts BuilderOptions, kind, id string) (interface{}, error) {
	opts = normalizeBuilderOptions(opts)
	cfg, err := LoadConfig(opts.ProjectRoot, opts.WorkspaceRoot)
	if err != nil {
		return nil, err
	}
	switch kind {
	case "circuits":
		return LoadCircuit(opts.ProjectRoot, cfg.RootDir, id)
	case "machines":
		return LoadMachine(opts.ProjectRoot, cfg.RootDir, id)
	case "automations":
		return LoadAutomation(opts.ProjectRoot, cfg.RootDir, id)
	case "commands":
		reg, err := LoadRegistry(opts.ProjectRoot, cfg)
		if err != nil {
			return nil, err
		}
		def, ok := ResolveCommand(reg, id)
		if !ok {
			return nil, os.ErrNotExist
		}
		return def, nil
	default:
		return nil, fmt.Errorf("unsupported builder primitive %q", kind)
	}
}

func normalizeBuilderOptions(opts BuilderOptions) BuilderOptions {
	if opts.ProjectRoot == "" {
		opts.ProjectRoot = "."
	}
	if opts.WorkspaceRoot == "" {
		opts.WorkspaceRoot = DefaultWorkspaceRoot
	}
	return opts
}

func validateBuilderID(id, label string) error {
	if !builderIDPattern.MatchString(strings.TrimSpace(id)) {
		return fmt.Errorf("%s id must use lowercase letters, numbers, and hyphens", label)
	}
	return nil
}

func writeProjectResource(projectRoot, resource, body string) error {
	path, ok := projectResourcePath(projectRoot, resource)
	if !ok {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(body), 0o644)
}

func defaultBuilderCircuitProgram(id string) string {
	return fmt.Sprintf(`# %s

---
factory_runtime_contract: 1
factory_prompt_contract: project-circuit.v1
---

<comment>
Goal: define this project-created circuit.
</comment>

<code>
return {
  status: "succeeded",
  summary: "Circuit completed.",
  data: Runtime.input optional default {}
}
</code>
`, id)
}

func defaultBuilderCommandPrompt(name string) string {
	return fmt.Sprintf(`# %s

---
factory_runtime_contract: 1
factory_prompt_contract: project-command.v1
---

<comment>
Goal: define this project-created command.
</comment>

<code>
return {
  status: "succeeded",
  resultType: "command.completed",
  summary: "Command completed.",
  events: [],
  data: Runtime.input optional default {}
}
</code>
`, name)
}

func sortedValues[T interface{ GetID() string }](values map[string]T) []T {
	out := make([]T, 0, len(values))
	for _, value := range values {
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetID() < out[j].GetID() })
	return out
}

func (d CircuitDefinition) GetID() string    { return d.ID }
func (d MachineDefinition) GetID() string    { return d.ID }
func (d AutomationDefinition) GetID() string { return d.ID }
func (d CommandDefinition) GetID() string    { return d.Name }
