package factory

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGUIHandlerServesHTMLAndGraphAPI(t *testing.T) {
	root := t.TempDir()
	if err := InitWorkspace(root, DefaultWorkspaceRoot); err != nil {
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

	graphResp := httptest.NewRecorder()
	handler.ServeHTTP(graphResp, httptest.NewRequest(http.MethodGet, "/api/graph", nil))
	if graphResp.Code != http.StatusOK {
		t.Fatalf("expected graph status 200, got %d", graphResp.Code)
	}
	if !strings.Contains(graphResp.Body.String(), `"automation-init"`) {
		t.Fatalf("expected init automation in graph API, got:\n%s", graphResp.Body.String())
	}
}
