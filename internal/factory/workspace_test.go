package factory

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateCommandRegistersProjectCommand(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	def, err := CreateCommand(root, cfg, CommandDefinition{
		Name:        "load-intake",
		Aliases:     []string{"load"},
		Description: "Load intake from operator input.",
	}, "prompt body")
	if err != nil {
		t.Fatal(err)
	}
	if def.ResultSchema != ResultContractV1 {
		t.Fatalf("expected default result schema, got %q", def.ResultSchema)
	}
	reg, err := LoadRegistry(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := ResolveCommand(reg, "load"); !ok {
		t.Fatal("expected alias to resolve")
	}
	promptPath := filepath.Join(root, ".factory", "00-control-room", "b-commands", "load-intake.md")
	data, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "prompt body" {
		t.Fatalf("unexpected prompt body: %q", string(data))
	}
}

func TestInitWorkspaceScaffoldsInspectableInitAutomation(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.FactorySkill.WorkspaceSchema != LatestWorkspaceSchema {
		t.Fatalf("unexpected workspace schema %d", cfg.FactorySkill.WorkspaceSchema)
	}
	if cfg.FactorySkill.RuntimeContract != LatestRuntimeContract {
		t.Fatalf("unexpected runtime contract %d", cfg.FactorySkill.RuntimeContract)
	}
	assertPathExists(t, root, ".factory/runtime/runtime-kernel.md")
	assertPathExists(t, root, ".factory/circuits/init/circuit.json")
	assertPathExists(t, root, ".factory/circuits/init/program.md")
	assertPathExists(t, root, ".factory/machines/init.json")
	assertPathExists(t, root, ".factory/automations/init.json")
	assertPathMissing(t, root, ".factory/triangulation/roadmap.md")
}

func TestInitWorkspaceUsesConfiguredWorkspaceRootInCircuitResources(t *testing.T) {
	root := t.TempDir()
	workspaceRoot := ".custom-factory"
	if err := InitWorkspace(root, workspaceRoot); err != nil {
		t.Fatal(err)
	}
	circuit, err := LoadCircuit(root, workspaceRoot, "init")
	if err != nil {
		t.Fatal(err)
	}
	if circuit.RuntimeKernel != "project://.custom-factory/runtime/runtime-kernel.md" {
		t.Fatalf("unexpected runtime kernel %q", circuit.RuntimeKernel)
	}
}

func TestInitAutomationRunsOnlyWhenExplicitlyDispatched(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	runner := &fakeInitRunner{}
	result, err := RunAutomation(context.Background(), OrchestratorOptions{
		ProjectRoot:   root,
		WorkspaceRoot: DefaultWorkspaceRoot,
		Runner:        runner,
	}, "init", map[string]interface{}{"requested_outcome": "Create a test factory."})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "succeeded" {
		t.Fatalf("unexpected status %q", result.Status)
	}
	if len(runner.calls) != 1 || runner.calls[0] != "init" {
		t.Fatalf("expected explicit init dispatch, got %#v", runner.calls)
	}
}

func TestUpdateWorkspaceRecordsMigrationEvidence(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	if err := UpdateWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, ".factory", "00-control-room", "d-migrations", "applied.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var event map[string]interface{}
	if err := json.Unmarshal(data[:len(data)-1], &event); err != nil {
		t.Fatal(err)
	}
	if event["status"] != "already-current" {
		t.Fatalf("unexpected migration event: %#v", event)
	}
	cfg, err := LoadConfig(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	configData, err := os.ReadFile(filepath.Join(root, ".factory", "00-control-room", "a-config", "factory.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(configData) {
		t.Fatal("expected persisted config to remain valid JSON")
	}
	if cfg.Migrations.LogPath == "" || cfg.FactorySkill.WorkspaceSchema != LatestWorkspaceSchema {
		t.Fatalf("expected persisted update metadata, got %#v", cfg)
	}
}

type fakeInitRunner struct {
	calls []string
}

func (r *fakeInitRunner) RunCircuit(_ context.Context, circuit CircuitDefinition, input map[string]interface{}) (CircuitResult, error) {
	r.calls = append(r.calls, circuit.ID)
	return CircuitResult{
		Status:  "succeeded",
		Summary: "init done",
		Data: map[string]interface{}{
			"requested_outcome": input["requested_outcome"],
		},
	}, nil
}

func assertPathExists(t *testing.T, root, rel string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
		t.Fatalf("expected %s to exist: %v", rel, err)
	}
}

func assertPathMissing(t *testing.T, root, rel string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be missing, got %v", rel, err)
	}
}
