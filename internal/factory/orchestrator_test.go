package factory

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

type fakeCircuitRunner struct {
	calls []string
}

func (r *fakeCircuitRunner) RunCircuit(_ context.Context, circuit CircuitDefinition, input map[string]interface{}) (CircuitResult, error) {
	r.calls = append(r.calls, circuit.ID)
	return CircuitResult{
		Status:  "succeeded",
		Summary: circuit.ID + " done",
		Data: map[string]interface{}{
			"from": circuit.ID,
			"prev": input["from"],
		},
	}, nil
}

func TestRunMachineRunsCircuitsSerially(t *testing.T) {
	root := t.TempDir()
	writeFactoryPiece(t, root, ".factory/circuits/init/circuit.json", `{"id":"init","program":".factory/circuits/init/program.md"}`)
	writeFactoryPiece(t, root, ".factory/circuits/plan/circuit.json", `{"id":"plan","program":".factory/circuits/plan/program.md"}`)
	writeFactoryPiece(t, root, ".factory/machines/build.json", `{"id":"build","circuits":[{"circuit":"init"},{"circuit":"plan","input":"previous.output"}]}`)

	runner := &fakeCircuitRunner{}
	result, err := RunMachine(context.Background(), OrchestratorOptions{
		ProjectRoot:   root,
		WorkspaceRoot: ".factory",
		Runner:        runner,
	}, "build", map[string]interface{}{"from": "input"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "succeeded" {
		t.Fatalf("unexpected status %q", result.Status)
	}
	if got := runner.calls; len(got) != 2 || got[0] != "init" || got[1] != "plan" {
		t.Fatalf("unexpected circuit order: %#v", got)
	}
	if result.Circuits[1].Input.(map[string]interface{})["from"] != "init" {
		t.Fatalf("expected second circuit to receive previous output, got %#v", result.Circuits[1].Input)
	}
}

func TestRunAutomationRunsMachinesSerially(t *testing.T) {
	root := t.TempDir()
	writeFactoryPiece(t, root, ".factory/circuits/init/circuit.json", `{"id":"init","program":".factory/circuits/init/program.md"}`)
	writeFactoryPiece(t, root, ".factory/circuits/ship/circuit.json", `{"id":"ship","program":".factory/circuits/ship/program.md"}`)
	writeFactoryPiece(t, root, ".factory/machines/build.json", `{"id":"build","circuits":[{"circuit":"init"}]}`)
	writeFactoryPiece(t, root, ".factory/machines/release.json", `{"id":"release","circuits":[{"circuit":"ship"}]}`)
	writeFactoryPiece(t, root, ".factory/automations/process.json", `{"id":"process","machines":[{"machine":"build"},{"machine":"release","input":"previous.output"}]}`)

	runner := &fakeCircuitRunner{}
	result, err := RunAutomation(context.Background(), OrchestratorOptions{
		ProjectRoot:   root,
		WorkspaceRoot: ".factory",
		Runner:        runner,
	}, "process", map[string]interface{}{"from": "input"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "succeeded" {
		t.Fatalf("unexpected status %q", result.Status)
	}
	if got := runner.calls; len(got) != 2 || got[0] != "init" || got[1] != "ship" {
		t.Fatalf("unexpected circuit order: %#v", got)
	}
}

func writeFactoryPiece(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
