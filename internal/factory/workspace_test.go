package factory

import (
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
