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
	cfg, err := LoadConfig(root, DefaultWorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateCommand(root, cfg, CommandDefinition{Name: "load-intake"}, "prompt"); err != nil {
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
	if !strings.Contains(htmlResp.Body.String(), "Factory Builder") {
		t.Fatalf("expected builder panel, got:\n%s", htmlResp.Body.String())
	}
	if !strings.Contains(htmlResp.Body.String(), `node.type === "command"`) {
		t.Fatalf("expected command-specific run affordance, got:\n%s", htmlResp.Body.String())
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

	createCircuitResp := httptest.NewRecorder()
	handler.ServeHTTP(createCircuitResp, httptest.NewRequest(http.MethodPost, "/api/builder/circuits", strings.NewReader(`{"id":"plan-work","name":"Plan Work"}`)))
	if createCircuitResp.Code != http.StatusOK {
		t.Fatalf("expected create circuit status 200, got %d: %s", createCircuitResp.Code, createCircuitResp.Body.String())
	}
	createMachineResp := httptest.NewRecorder()
	handler.ServeHTTP(createMachineResp, httptest.NewRequest(http.MethodPost, "/api/builder/machines", strings.NewReader(`{"id":"build-work","circuits":[{"circuit":"plan-work"}]}`)))
	if createMachineResp.Code != http.StatusOK {
		t.Fatalf("expected create machine status 200, got %d: %s", createMachineResp.Code, createMachineResp.Body.String())
	}
	invalidMachineResp := httptest.NewRecorder()
	handler.ServeHTTP(invalidMachineResp, httptest.NewRequest(http.MethodPost, "/api/builder/machines", strings.NewReader(`{"id":"bad-work","circuits":[{"circuit":"missing"}]}`)))
	if invalidMachineResp.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid machine status 400, got %d", invalidMachineResp.Code)
	}
}

func TestBuilderAPIUsesInjectedStore(t *testing.T) {
	store := &fakeBuilderStore{
		inventory: BuilderInventory{
			Circuits: []CircuitDefinition{{ID: "from-store", Name: "From Store"}},
		},
	}
	handler := NewGUIHandler(GUIOptions{
		ProjectRoot:   t.TempDir(),
		WorkspaceRoot: DefaultWorkspaceRoot,
		BuilderStore:  store,
	})

	listResp := httptest.NewRecorder()
	handler.ServeHTTP(listResp, httptest.NewRequest(http.MethodGet, "/api/builder/circuits", nil))
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list status 200, got %d", listResp.Code)
	}
	if !strings.Contains(listResp.Body.String(), "from-store") {
		t.Fatalf("expected injected store response, got:\n%s", listResp.Body.String())
	}

	createResp := httptest.NewRecorder()
	handler.ServeHTTP(createResp, httptest.NewRequest(http.MethodPost, "/api/builder/circuits", strings.NewReader(`{"id":"api-circuit","name":"API Circuit"}`)))
	if createResp.Code != http.StatusOK {
		t.Fatalf("expected create status 200, got %d: %s", createResp.Code, createResp.Body.String())
	}
	if store.savedCircuit.ID != "api-circuit" {
		t.Fatalf("expected API to call injected store, got %#v", store.savedCircuit)
	}
}

type fakeBuilderStore struct {
	inventory    BuilderInventory
	savedCircuit CircuitDefinition
}

func (s *fakeBuilderStore) ListInventory() (BuilderInventory, error) {
	return s.inventory, nil
}

func (s *fakeBuilderStore) ReadPrimitive(kind, id string) (interface{}, error) {
	if kind == "circuits" && id == "from-store" {
		return s.inventory.Circuits[0], nil
	}
	return nil, os.ErrNotExist
}

func (s *fakeBuilderStore) SaveCircuit(input BuilderCircuitInput) (CircuitDefinition, error) {
	s.savedCircuit = CircuitDefinition{ID: input.ID, Name: input.Name}
	s.inventory.Circuits = append(s.inventory.Circuits, s.savedCircuit)
	return s.savedCircuit, nil
}

func (s *fakeBuilderStore) SaveMachine(def MachineDefinition) (MachineDefinition, error) {
	return def, nil
}

func (s *fakeBuilderStore) SaveAutomation(def AutomationDefinition) (AutomationDefinition, error) {
	return def, nil
}

func (s *fakeBuilderStore) SaveCommand(input BuilderCommandInput) (CommandDefinition, error) {
	return input.CommandDefinition, nil
}

func (s *fakeBuilderStore) CreateEventBinding(binding EventBinding) (EventBinding, error) {
	return binding, nil
}

func (s *fakeBuilderStore) UpdateEventBinding(_ int, binding EventBinding) (EventBinding, error) {
	return binding, nil
}
