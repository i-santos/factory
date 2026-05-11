package factory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateWorkspaceAcceptsInitializedSharedModel(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	result, err := ValidateWorkspace(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "valid" {
		t.Fatalf("expected valid workspace, got %#v", result)
	}
	if result.Pieces.Circuits != 1 || result.Pieces.Machines != 1 || result.Pieces.Automations != 1 {
		t.Fatalf("unexpected inventory: %#v", result.Pieces)
	}
}

func TestValidateWorkspaceRejectsMissingCircuitReferences(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".factory", "machines", "broken.json")
	if err := os.WriteFile(path, []byte(`{"id":"broken","circuits":[{"circuit":"missing"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := ValidateWorkspace(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "invalid" {
		t.Fatalf("expected invalid workspace, got %#v", result)
	}
	if len(result.Issues) != 1 || !strings.Contains(result.Issues[0].Message, "missing circuit") {
		t.Fatalf("expected missing circuit issue, got %#v", result.Issues)
	}
	if _, err := BuildFactoryGraph(VisualizationOptions{ProjectRoot: root, WorkspaceRoot: DefaultWorkspaceRoot}); err == nil {
		t.Fatal("expected graph build to fail validation")
	}
}
