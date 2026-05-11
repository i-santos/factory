package factory

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultWorkspaceRoot    = ".factory"
	FactorySkillVersion     = "0.1.0"
	LatestRuntimeContract   = 1
	LatestWorkspaceSchema   = 1
	ResultContractV1        = "factory://schemas/command-result.v1.json"
	CircuitResultContractV1 = "factory://schemas/circuit-result.v1.json"
)

func DefaultConfig(root string) Config {
	var cfg Config
	cfg.SchemaVersion = 1
	cfg.RootDir = root
	cfg.FactorySkill.Name = "factory"
	cfg.FactorySkill.RequiredSkillVersion = FactorySkillVersion
	cfg.FactorySkill.RuntimeContract = LatestRuntimeContract
	cfg.FactorySkill.WorkspaceSchema = LatestWorkspaceSchema
	cfg.FactorySkill.ResourceMode = "installed-skill"
	cfg.Commands.RegistryPath = filepath.Join(root, "00-control-room", "a-config", "commands.json")
	cfg.Commands.ProjectCommandsDir = filepath.Join(root, "00-control-room", "b-commands")
	cfg.Events.BindingsPath = filepath.Join(root, "00-control-room", "a-config", "event-bindings.json")
	cfg.Events.LogDir = filepath.Join(root, "00-control-room", "d-events")
	cfg.Events.MaxAutoIterations = 25
	cfg.Migrations.LogPath = filepath.Join(root, "00-control-room", "d-migrations", "applied.jsonl")
	cfg.Migrations.BackupDir = filepath.Join(root, "00-control-room", "d-migrations", "backups")
	cfg.Migrations.StagingDir = filepath.Join(root, "00-control-room", "d-migrations", "staging")
	cfg.Sectors.RootDir = filepath.Join(root, "05-sectors")
	return cfg
}

func InitWorkspace(projectRoot, workspaceRoot string) error {
	cfg := DefaultConfig(workspaceRoot)
	dirs := []string{
		"00-control-room/a-config",
		"00-control-room/b-commands",
		"00-control-room/c-schemas",
		"00-control-room/d-events",
		"00-control-room/d-migrations",
		"00-control-room/d-migrations/backups",
		"00-control-room/d-migrations/staging",
		"00-control-room/e-state",
		"00-control-room/f-logs",
		"01-dock/a-incoming",
		"01-dock/b-triaged",
		"01-dock/c-rejected",
		"01-dock/d-immediate-response",
		"02-yard/a-pending",
		"02-yard/b-refining",
		"02-yard/c-discharged",
		"02-yard/d-archived",
		"03-shop-floor/a-input-buffer",
		"03-shop-floor/b-wip",
		"03-shop-floor/c-hold",
		"03-shop-floor/d-work-packages",
		"04-finished-goods/a-shipped",
		"04-finished-goods/b-archived",
		"05-sectors",
		"runtime",
		"circuits/init",
		"machines",
		"automations",
		"triangulation",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectRoot, workspaceRoot, filepath.FromSlash(dir)), 0o755); err != nil {
			return err
		}
	}
	if err := writeJSONIfMissing(filepath.Join(projectRoot, cfg.RootDir, "00-control-room", "a-config", "factory.json"), cfg); err != nil {
		return err
	}
	reg := CommandRegistry{SchemaVersion: 1, Commands: map[string]CommandDefinition{}}
	if err := writeJSONIfMissing(filepath.Join(projectRoot, cfg.Commands.RegistryPath), reg); err != nil {
		return err
	}
	bindings := EventBindingsFile{SchemaVersion: 1, Bindings: []EventBinding{}}
	if err := writeJSONIfMissing(filepath.Join(projectRoot, cfg.Events.BindingsPath), bindings); err != nil {
		return err
	}
	if err := writeTextIfMissing(filepath.Join(projectRoot, workspaceRoot, "runtime", "runtime-kernel.md"), defaultRuntimeKernel()); err != nil {
		return err
	}
	if err := writeTextIfMissing(filepath.Join(projectRoot, workspaceRoot, "circuits", "init", "program.md"), defaultInitCircuitProgram(workspaceRoot)); err != nil {
		return err
	}
	if err := writeJSONIfMissing(filepath.Join(projectRoot, workspaceRoot, "circuits", "init", "circuit.json"), CircuitDefinition{
		ID:            "init",
		Name:          "Initialize factory triangulation",
		Program:       filepath.ToSlash(filepath.Join(workspaceRoot, "circuits", "init", "program.md")),
		RuntimeKernel: "project://" + filepath.ToSlash(filepath.Join(workspaceRoot, "runtime", "runtime-kernel.md")),
		Description:   "Creates starter triangulation artifacts for a newly scaffolded factory.",
	}); err != nil {
		return err
	}
	if err := writeJSONIfMissing(filepath.Join(projectRoot, workspaceRoot, "machines", "init.json"), MachineDefinition{
		ID:          "init",
		Name:        "Initialize factory",
		Description: "Runs the starter init circuit.",
		Circuits: []MachineCircuitStep{
			{Circuit: "init"},
		},
	}); err != nil {
		return err
	}
	return writeJSONIfMissing(filepath.Join(projectRoot, workspaceRoot, "automations", "init.json"), AutomationDefinition{
		ID:          "init",
		Name:        "Initialize factory",
		Description: "Runs the starter factory initialization machine.",
		Machines: []AutomationMachineStep{
			{Machine: "init"},
		},
	})
}

func LoadConfig(projectRoot, workspaceRoot string) (Config, error) {
	path := filepath.Join(projectRoot, workspaceRoot, "00-control-room", "a-config", "factory.json")
	var cfg Config
	if err := readJSON(path, &cfg); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(workspaceRoot), nil
		}
		return cfg, err
	}
	if cfg.RootDir == "" {
		cfg.RootDir = workspaceRoot
	}
	normalizeConfigDefaults(&cfg, workspaceRoot)
	return cfg, nil
}

func UpdateWorkspace(projectRoot, workspaceRoot string) error {
	cfg, err := LoadConfig(projectRoot, workspaceRoot)
	if err != nil {
		return err
	}
	if cfg.FactorySkill.WorkspaceSchema > LatestWorkspaceSchema {
		return fmt.Errorf("workspace schema %d is newer than supported schema %d", cfg.FactorySkill.WorkspaceSchema, LatestWorkspaceSchema)
	}
	fromSchema := cfg.FactorySkill.WorkspaceSchema
	if err := InitWorkspace(projectRoot, workspaceRoot); err != nil {
		return err
	}
	cfg, err = LoadConfig(projectRoot, workspaceRoot)
	if err != nil {
		return err
	}
	cfg.FactorySkill.WorkspaceSchema = LatestWorkspaceSchema
	cfg.FactorySkill.RuntimeContract = LatestRuntimeContract
	if err := writeJSON(filepath.Join(projectRoot, cfg.RootDir, "00-control-room", "a-config", "factory.json"), cfg); err != nil {
		return err
	}
	return appendMigrationEvent(projectRoot, cfg, map[string]interface{}{
		"fromSchema": fromSchema,
		"toSchema":   LatestWorkspaceSchema,
		"status":     "already-current",
	})
}

func LoadRegistry(projectRoot string, cfg Config) (CommandRegistry, error) {
	reg := CommandRegistry{SchemaVersion: 1, Commands: map[string]CommandDefinition{}}
	err := readJSON(filepath.Join(projectRoot, cfg.Commands.RegistryPath), &reg)
	if errors.Is(err, os.ErrNotExist) {
		return reg, nil
	}
	if reg.Commands == nil {
		reg.Commands = map[string]CommandDefinition{}
	}
	return reg, err
}

func SaveRegistry(projectRoot string, cfg Config, reg CommandRegistry) error {
	if reg.SchemaVersion == 0 {
		reg.SchemaVersion = 1
	}
	if reg.Commands == nil {
		reg.Commands = map[string]CommandDefinition{}
	}
	return writeJSON(filepath.Join(projectRoot, cfg.Commands.RegistryPath), reg)
}

func CreateCommand(projectRoot string, cfg Config, def CommandDefinition, promptBody string) (CommandDefinition, error) {
	if def.Name == "" {
		return def, fmt.Errorf("command name is required")
	}
	if def.ResultSchema == "" {
		def.ResultSchema = ResultContractV1
	}
	if len(def.Emits) == 0 {
		def.Emits = []string{"command.completed", "command.blocked", "command.failed", "command.waiting"}
	}
	if len(def.Consumes) == 0 {
		def.Consumes = []string{"operator.requested"}
	}
	if def.Prompt == "" {
		rel := filepath.Join(cfg.Commands.ProjectCommandsDir, def.Name+".md")
		def.Prompt = "project://" + filepath.ToSlash(rel)
	}
	reg, err := LoadRegistry(projectRoot, cfg)
	if err != nil {
		return def, err
	}
	if _, exists := reg.Commands[def.Name]; exists {
		return def, fmt.Errorf("command %q already exists", def.Name)
	}
	for name, existing := range reg.Commands {
		for _, alias := range def.Aliases {
			if alias == name || contains(existing.Aliases, alias) {
				return def, fmt.Errorf("alias %q conflicts with command %q", alias, name)
			}
		}
	}
	if err := writePromptIfMissing(projectRoot, def.Prompt, promptBody); err != nil {
		return def, err
	}
	reg.Commands[def.Name] = def
	return def, SaveRegistry(projectRoot, cfg, reg)
}

func ResolveCommand(reg CommandRegistry, name string) (CommandDefinition, bool) {
	if def, ok := reg.Commands[name]; ok {
		return def, true
	}
	for _, def := range reg.Commands {
		if contains(def.Aliases, name) {
			return def, true
		}
	}
	return CommandDefinition{}, false
}

func LoadEventBindings(projectRoot string, cfg Config) (EventBindingsFile, error) {
	file := EventBindingsFile{SchemaVersion: 1, Bindings: []EventBinding{}}
	err := readJSON(filepath.Join(projectRoot, cfg.Events.BindingsPath), &file)
	if errors.Is(err, os.ErrNotExist) {
		return file, nil
	}
	return file, err
}

func SaveEventBindings(projectRoot string, cfg Config, file EventBindingsFile) error {
	if file.SchemaVersion == 0 {
		file.SchemaVersion = 1
	}
	if file.Bindings == nil {
		file.Bindings = []EventBinding{}
	}
	return writeJSON(filepath.Join(projectRoot, cfg.Events.BindingsPath), file)
}

func AddEventBinding(projectRoot string, cfg Config, binding EventBinding) (EventBinding, error) {
	if binding.On == "" {
		return binding, fmt.Errorf("binding event type is required")
	}
	if binding.Run == "" {
		return binding, fmt.Errorf("binding target command is required")
	}
	if binding.Mode == "" {
		binding.Mode = "auto"
	}
	file, err := LoadEventBindings(projectRoot, cfg)
	if err != nil {
		return binding, err
	}
	file.Bindings = append(file.Bindings, binding)
	return binding, SaveEventBindings(projectRoot, cfg, file)
}

func writePromptIfMissing(projectRoot, resource, body string) error {
	path, ok := projectResourcePath(projectRoot, resource)
	if !ok {
		return nil
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(body), 0o644)
}

func projectResourcePath(projectRoot, resource string) (string, bool) {
	const prefix = "project://"
	if len(resource) < len(prefix) || resource[:len(prefix)] != prefix {
		return "", false
	}
	return filepath.Join(projectRoot, filepath.FromSlash(resource[len(prefix):])), true
}

func writeJSONIfMissing(path string, value interface{}) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return writeJSON(path, value)
}

func writeTextIfMissing(path, value string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(value), 0o644)
}

func writeJSON(path string, value interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func readJSON(path string, value interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func normalizeConfigDefaults(cfg *Config, workspaceRoot string) {
	if cfg.FactorySkill.Name == "" {
		cfg.FactorySkill.Name = "factory"
	}
	if cfg.FactorySkill.RequiredSkillVersion == "" {
		cfg.FactorySkill.RequiredSkillVersion = FactorySkillVersion
	}
	if cfg.FactorySkill.RuntimeContract == 0 {
		cfg.FactorySkill.RuntimeContract = LatestRuntimeContract
	}
	if cfg.FactorySkill.WorkspaceSchema == 0 {
		cfg.FactorySkill.WorkspaceSchema = LatestWorkspaceSchema
	}
	if cfg.FactorySkill.ResourceMode == "" {
		cfg.FactorySkill.ResourceMode = "installed-skill"
	}
	if cfg.Migrations.LogPath == "" {
		cfg.Migrations.LogPath = filepath.Join(workspaceRoot, "00-control-room", "d-migrations", "applied.jsonl")
	}
	if cfg.Migrations.BackupDir == "" {
		cfg.Migrations.BackupDir = filepath.Join(workspaceRoot, "00-control-room", "d-migrations", "backups")
	}
	if cfg.Migrations.StagingDir == "" {
		cfg.Migrations.StagingDir = filepath.Join(workspaceRoot, "00-control-room", "d-migrations", "staging")
	}
}

func appendMigrationEvent(projectRoot string, cfg Config, event map[string]interface{}) error {
	path := filepath.Join(projectRoot, cfg.Migrations.LogPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	event["workspaceSchema"] = LatestWorkspaceSchema
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(data, '\n'))
	return err
}

func defaultRuntimeKernel() string {
	return `# Factory Runtime Kernel

This project-local kernel records the runtime contract expected by scaffolded
circuits. The installed Factory skill remains the source of executable runtime
behavior for agent sessions.

Required result envelope:

` + "```json" + `
{
  "status": "succeeded|blocked|failed|waiting",
  "summary": "Short operator-facing summary.",
  "data": {},
  "artifacts": [],
  "errors": []
}
` + "```" + `
`
}

func defaultInitCircuitProgram(workspaceRoot string) string {
	return fmt.Sprintf(`# Initialize Factory Triangulation

---
factory_runtime_contract: %d
factory_prompt_contract: factory-init-circuit.v1
---

<comment>
Goal: create starter triangulation artifacts for a newly scaffolded Factory
workspace. This circuit runs only when the user explicitly dispatches the init
automation, for example with `+"`factory run automation init`"+`.
</comment>

<code>
input = Runtime.input optional default {}

Trace.event({
  eventType: "start",
  summary: "Factory init triangulation circuit started."
})

triangulation = Pseudo.CreateStarterTriangulation({
  workspace_root: %q,
  requested_outcome: input.requested_outcome optional default "Initialize a reusable local factory."
})

return {
  status: "succeeded",
  summary: "Starter factory triangulation created.",
  data: {
    triangulation: triangulation
  },
  artifacts: triangulation.artifacts,
  errors: []
}
</code>
`, LatestRuntimeContract, workspaceRoot)
}

func contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}
