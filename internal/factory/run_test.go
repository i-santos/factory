package factory

import (
	"context"
	"encoding/json"
	"testing"
)

type sequenceRunner struct {
	results []CommandResult
	calls   []string
}

func (r *sequenceRunner) Run(_ context.Context, command CommandDefinition, input map[string]interface{}) (CommandResult, error) {
	r.calls = append(r.calls, command.Name)
	next := r.results[0]
	r.results = r.results[1:]
	return next, nil
}

func TestRunCommandChainFollowsEventBinding(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateCommand(root, cfg, CommandDefinition{Name: "load-intake"}, "load"); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateCommand(root, cfg, CommandDefinition{Name: "drain-work-package"}, "drain"); err != nil {
		t.Fatal(err)
	}
	bindings := EventBindingsFile{
		SchemaVersion: 1,
		Bindings: []EventBinding{
			{
				On: "command.completed",
				Where: map[string]string{
					"command":       "load-intake",
					"result.status": "succeeded",
				},
				Run: "drain-work-package",
				Input: map[string]interface{}{
					"package_path": "$result.data.package_path",
				},
				Mode: "auto",
			},
		},
	}
	if err := writeJSON(join(root, cfg.Events.BindingsPath), bindings); err != nil {
		t.Fatal(err)
	}
	runner := &sequenceRunner{
		results: []CommandResult{
			{
				Status:     "succeeded",
				ResultType: "work-package-created",
				Data:       map[string]interface{}{"package_path": ".factory/03-shop-floor/d-work-packages/pkg.dir"},
			},
			{
				Status:     "succeeded",
				ResultType: "work-package-drained",
			},
		},
	}
	records, err := RunCommandChain(context.Background(), RunOptions{
		ProjectRoot:       root,
		WorkspaceRoot:     DefaultWorkspaceRoot,
		Command:           "load-intake",
		MaxAutoIterations: 5,
		Runner:            runner,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if runner.calls[0] != "load-intake" || runner.calls[1] != "drain-work-package" {
		t.Fatalf("unexpected calls: %#v", runner.calls)
	}
	got := records[1].Input["package_path"]
	if got != ".factory/03-shop-floor/d-work-packages/pkg.dir" {
		t.Fatalf("unexpected chained input: %#v", got)
	}
}

func TestParseCommandResultRequiresEnvelopeStatus(t *testing.T) {
	_, err := ParseCommandResult([]byte(`{"summary":"missing status"}`))
	if err == nil {
		t.Fatal("expected missing status error")
	}
	var out map[string]interface{}
	data, _ := json.Marshal(CommandResult{Status: "succeeded"})
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
}

func join(root, rel string) string {
	return root + "/" + rel
}
