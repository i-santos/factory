package factory

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGUIHandlerServesHTMLAndGraphAPI(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
		t.Fatal(err)
	}
	runDir := filepath.Join(root, DefaultWorkspaceRoot, "runs")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	run := RunSession{
		RunID:      "run_gui",
		TargetType: "command-chain",
		TargetID:   "prepare",
		Status:     RunStatusNeedsHumanAction,
		DecisionRequest: &DecisionRequest{
			ID:          "decision_gui",
			Title:       "Choose route",
			Description: "The factory needs a human decision.",
			Controls: []DecisionControl{
				{
					Type: "option-list",
					Name: "route",
					Options: []DecisionControlOption{
						{Value: "continue", Label: "Continue"},
					},
				},
			},
		},
	}
	data, err := json.Marshal(run)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "run_gui.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	handler := NewGUIHandler(GUIOptions{ProjectRoot: root, WorkspaceRoot: DefaultWorkspaceRoot})

	htmlResp := httptest.NewRecorder()
	handler.ServeHTTP(htmlResp, httptest.NewRequest(http.MethodGet, "/", nil))
	if htmlResp.Code != http.StatusOK {
		t.Fatalf("expected HTML status 200, got %d", htmlResp.Code)
	}
	if !strings.Contains(htmlResp.Body.String(), "Production System") {
		t.Fatalf("expected GUI shell, got:\n%s", htmlResp.Body.String())
	}
	if !strings.Contains(htmlResp.Body.String(), "Run Monitor") {
		t.Fatalf("expected run monitor, got:\n%s", htmlResp.Body.String())
	}
	if !strings.Contains(htmlResp.Body.String(), "Choose route") {
		t.Fatalf("expected decision request in GUI, got:\n%s", htmlResp.Body.String())
	}

	graphResp := httptest.NewRecorder()
	handler.ServeHTTP(graphResp, httptest.NewRequest(http.MethodGet, "/api/graph", nil))
	if graphResp.Code != http.StatusOK {
		t.Fatalf("expected graph status 200, got %d", graphResp.Code)
	}
	if !strings.Contains(graphResp.Body.String(), `"automation-init"`) {
		t.Fatalf("expected init automation in graph API, got:\n%s", graphResp.Body.String())
	}

	runsResp := httptest.NewRecorder()
	handler.ServeHTTP(runsResp, httptest.NewRequest(http.MethodGet, "/api/runs", nil))
	if runsResp.Code != http.StatusOK {
		t.Fatalf("expected runs status 200, got %d", runsResp.Code)
	}
	if !strings.Contains(runsResp.Body.String(), `"needs-human-action"`) {
		t.Fatalf("expected human action run in API, got:\n%s", runsResp.Body.String())
	}
}
