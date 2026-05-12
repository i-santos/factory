package factory

import (
	"context"
	"encoding/json"
	"strings"
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

func TestRunCommandChainFollowsSupportedNextAction(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateCommand(root, cfg, CommandDefinition{Name: "prepare"}, "prepare"); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateCommand(root, cfg, CommandDefinition{Name: "finish"}, "finish"); err != nil {
		t.Fatal(err)
	}
	runner := &sequenceRunner{
		results: []CommandResult{
			{
				Status: "succeeded",
				NextAction: &NextAction{
					Kind:    "command",
					Command: "finish",
					Input:   map[string]interface{}{"value": "ok"},
				},
			},
			{Status: "succeeded"},
		},
	}
	records, err := RunCommandChain(context.Background(), RunOptions{
		ProjectRoot:       root,
		WorkspaceRoot:     DefaultWorkspaceRoot,
		Command:           "prepare",
		MaxAutoIterations: 5,
		RunID:             "run_supported_next",
		Runner:            runner,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if got := records[1].Input["value"]; got != "ok" {
		t.Fatalf("unexpected next-action input: %#v", got)
	}
	runs, err := ListRunSessions(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 1 || runs[0].Status != RunStatusSucceeded {
		t.Fatalf("unexpected persisted runs: %#v", runs)
	}
}

func TestRunCommandChainBlocksUnsupportedNextAction(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateCommand(root, cfg, CommandDefinition{Name: "prepare"}, "prepare"); err != nil {
		t.Fatal(err)
	}
	runner := &sequenceRunner{
		results: []CommandResult{
			{
				Status: "succeeded",
				NextAction: &NextAction{
					Kind: "unsupported",
				},
			},
		},
	}
	records, err := RunCommandChain(context.Background(), RunOptions{
		ProjectRoot:       root,
		WorkspaceRoot:     DefaultWorkspaceRoot,
		Command:           "prepare",
		MaxAutoIterations: 5,
		RunID:             "run_blocked_next",
		Runner:            runner,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	runs, err := ListRunSessions(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 persisted run, got %#v", runs)
	}
	if runs[0].Status != RunStatusBlocked {
		t.Fatalf("expected blocked run, got %#v", runs[0])
	}
	if !strings.Contains(runs[0].StopReason, "unsupported next action") {
		t.Fatalf("unexpected stop reason: %q", runs[0].StopReason)
	}
}

func TestRunCommandChainStopsForHumanDecision(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateCommand(root, cfg, CommandDefinition{Name: "prepare"}, "prepare"); err != nil {
		t.Fatal(err)
	}
	runner := &sequenceRunner{
		results: []CommandResult{
			{
				Status: "succeeded",
				NextAction: &NextAction{
					Kind:          "command",
					Command:       "finish",
					RequiresHuman: true,
				},
				DecisionRequest: &DecisionRequest{
					ID:    "decision_1",
					Title: "Choose next step",
					Controls: []DecisionControl{
						{
							Type: "option-list",
							Name: "action",
							Options: []DecisionControlOption{
								{Value: "continue", Label: "Continue"},
							},
						},
					},
				},
			},
		},
	}
	records, err := RunCommandChain(context.Background(), RunOptions{
		ProjectRoot:       root,
		WorkspaceRoot:     DefaultWorkspaceRoot,
		Command:           "prepare",
		MaxAutoIterations: 5,
		RunID:             "run_human",
		Runner:            runner,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	runs, err := ListRunSessions(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(runs) != 1 || runs[0].Status != RunStatusNeedsHumanAction {
		t.Fatalf("unexpected persisted human run: %#v", runs)
	}
	if runs[0].DecisionRequest == nil || runs[0].DecisionRequest.Title != "Choose next step" {
		t.Fatalf("missing decision request: %#v", runs[0])
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

func TestParseCommandResultAcceptsNeedsHumanAction(t *testing.T) {
	result, err := ParseCommandResult([]byte(`{"status":"needs-human-action"}`))
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != string(RunStatusNeedsHumanAction) {
		t.Fatalf("unexpected status: %q", result.Status)
	}
}

func join(root, rel string) string {
	return root + "/" + rel
}
