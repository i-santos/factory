package factory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildFactoryGraphIncludesWorkspaceCommandsBindingsAndSectors(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, cfg.Sectors.RootDir, "product", "actions"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, cfg.Sectors.RootDir, "product", "actions", "refine-experience.md"), []byte("action"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateCommand(root, cfg, CommandDefinition{Name: "load-intake"}, "load product refine-experience"); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateCommand(root, cfg, CommandDefinition{Name: "drain-work-package"}, "drain"); err != nil {
		t.Fatal(err)
	}
	if _, err := AddEventBinding(root, cfg, EventBinding{
		On:    "command.completed",
		Where: map[string]string{"command": "load-intake"},
		Run:   "drain-work-package",
		Mode:  "auto",
	}); err != nil {
		t.Fatal(err)
	}
	graph, err := BuildFactoryGraph(VisualizationOptions{ProjectRoot: root, WorkspaceRoot: DefaultWorkspaceRoot})
	if err != nil {
		t.Fatal(err)
	}
	assertNode(t, graph, "dock", "lane")
	assertNode(t, graph, "automation-init", "automation")
	assertNode(t, graph, "machine-init", "machine")
	assertNode(t, graph, "circuit-init", "circuit")
	assertNode(t, graph, "cmd-load-intake", "machine")
	assertNode(t, graph, "cmd-drain-work-package", "machine")
	assertNode(t, graph, "sector-product", "sector")
	assertNode(t, graph, "sector-product-action-refine-experience", "sector-action")
	assertEdge(t, graph, "dock", "yard", "workspace-flow")
	assertEdge(t, graph, "automation-init", "machine-init", "automation-machine")
	assertEdge(t, graph, "machine-init", "circuit-init", "machine-circuit")
	assertEdge(t, graph, "cmd-load-intake", "cmd-drain-work-package", "event-binding")
	assertEdge(t, graph, "cmd-load-intake", "sector-product", "sector-reference")
	assertEdge(t, graph, "cmd-load-intake", "sector-product-action-refine-experience", "sector-action-reference")
}

func TestRenderFactoryGraphHTMLIncludesMapAndRunAffordances(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	graph, err := BuildFactoryGraph(VisualizationOptions{ProjectRoot: root, WorkspaceRoot: DefaultWorkspaceRoot})
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := RenderFactoryGraphHTML(graph)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(rendered, "Factory Map") {
		t.Fatalf("expected GUI title, got:\n%s", rendered)
	}
	if !strings.Contains(rendered, "factory run automation ") {
		t.Fatalf("expected automation run affordance, got:\n%s", rendered)
	}
	if !strings.Contains(rendered, "oklch(") {
		t.Fatalf("expected OKLCH design tokens, got:\n%s", rendered)
	}
}

func TestRenderFactoryGraphMermaid(t *testing.T) {
	graph := FactoryGraph{
		SchemaVersion: 1,
		Nodes: []GraphNode{
			{ID: "dock", Label: "01 Dock", Type: "lane"},
			{ID: "cmd-load-intake", Label: "load-intake", Type: "machine"},
		},
		Edges: []GraphEdge{{From: "dock", To: "cmd-load-intake", Type: "contains", Label: "machine"}},
	}
	rendered := RenderFactoryGraphMermaid(graph)
	if !strings.Contains(rendered, "flowchart LR") {
		t.Fatalf("expected mermaid flowchart, got:\n%s", rendered)
	}
	if !strings.Contains(rendered, "cmd_load_intake[\"load-intake\"]") {
		t.Fatalf("expected command node, got:\n%s", rendered)
	}
}

func assertNode(t *testing.T, graph FactoryGraph, id, nodeType string) {
	t.Helper()
	for _, node := range graph.Nodes {
		if node.ID == id && node.Type == nodeType {
			return
		}
	}
	t.Fatalf("expected node %s of type %s in %#v", id, nodeType, graph.Nodes)
}

func assertEdge(t *testing.T, graph FactoryGraph, from, to, edgeType string) {
	t.Helper()
	for _, edge := range graph.Edges {
		if edge.From == from && edge.To == to && edge.Type == edgeType {
			return
		}
	}
	t.Fatalf("expected edge %s -> %s of type %s in %#v", from, to, edgeType, graph.Edges)
}
