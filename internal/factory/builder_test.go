package factory

import (
	"path/filepath"
	"testing"
)

func TestBuilderCreatesAndListsPrimitives(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	opts := BuilderOptions{ProjectRoot: root, WorkspaceRoot: DefaultWorkspaceRoot}
	if _, err := CreateOrUpdateCircuit(opts, BuilderCircuitInput{
		ID:          "plan-work",
		Name:        "Plan Work",
		Description: "Creates a plan.",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateOrUpdateMachine(opts, MachineDefinition{
		ID:   "build-work",
		Name: "Build Work",
		Circuits: []MachineCircuitStep{
			{Circuit: "plan-work"},
		},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateOrUpdateAutomation(opts, AutomationDefinition{
		ID:   "process-work",
		Name: "Process Work",
		Machines: []AutomationMachineStep{
			{Machine: "build-work"},
		},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateOrUpdateCommand(opts, BuilderCommandInput{
		CommandDefinition: CommandDefinition{Name: "load-intake"},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateEventBinding(opts, EventBinding{
		On:    "command.completed",
		Where: map[string]string{"command": "load-intake"},
		Run:   "load-intake",
	}); err != nil {
		t.Fatal(err)
	}
	inventory, err := ListBuilderInventory(opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(inventory.Circuits) != 2 || len(inventory.Machines) != 2 || len(inventory.Automations) != 2 {
		t.Fatalf("expected init plus created primitives, got %#v", inventory)
	}
	if len(inventory.Commands) != 1 || inventory.Commands[0].Name != "load-intake" {
		t.Fatalf("expected created command, got %#v", inventory.Commands)
	}
	if len(inventory.Bindings) != 1 {
		t.Fatalf("expected created binding, got %#v", inventory.Bindings)
	}
	assertPathExists(t, root, filepath.ToSlash(filepath.Join(DefaultWorkspaceRoot, "circuits", "plan-work", "program.md")))
}

func TestBuilderRejectsInvalidReferences(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	opts := BuilderOptions{ProjectRoot: root, WorkspaceRoot: DefaultWorkspaceRoot}
	if _, err := CreateOrUpdateMachine(opts, MachineDefinition{
		ID: "bad-machine",
		Circuits: []MachineCircuitStep{
			{Circuit: "missing-circuit"},
		},
	}); err == nil {
		t.Fatal("expected missing circuit reference to fail")
	}
	if _, err := CreateOrUpdateAutomation(opts, AutomationDefinition{
		ID: "bad-automation",
		Machines: []AutomationMachineStep{
			{Machine: "missing-machine"},
		},
	}); err == nil {
		t.Fatal("expected missing machine reference to fail")
	}
	if _, err := CreateEventBinding(opts, EventBinding{
		On:  "command.completed",
		Run: "missing-command",
	}); err == nil {
		t.Fatal("expected missing command binding target to fail")
	}
}
