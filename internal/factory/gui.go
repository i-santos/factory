package factory

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type GUIOptions struct {
	ProjectRoot   string
	WorkspaceRoot string
	Addr          string
}

func ServeGUI(opts GUIOptions) error {
	addr := opts.Addr
	if strings.TrimSpace(addr) == "" {
		addr = "127.0.0.1:8765"
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Printf("Factory GUI listening on http://%s\n", listener.Addr().String())
	return http.Serve(listener, NewGUIHandler(opts))
}

func NewGUIHandler(opts GUIOptions) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		graph, err := BuildFactoryGraph(VisualizationOptions{ProjectRoot: opts.ProjectRoot, WorkspaceRoot: opts.WorkspaceRoot})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		runs, err := ListRunSessions(opts.ProjectRoot, opts.WorkspaceRoot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		html, err := RenderFactoryGUIHTML(graph, runs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	})
	mux.HandleFunc("/api/graph", func(w http.ResponseWriter, r *http.Request) {
		graph, err := BuildFactoryGraph(VisualizationOptions{ProjectRoot: opts.ProjectRoot, WorkspaceRoot: opts.WorkspaceRoot})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(graph)
	})
	mux.HandleFunc("/api/runs", func(w http.ResponseWriter, r *http.Request) {
		runs, err := ListRunSessions(opts.ProjectRoot, opts.WorkspaceRoot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(runs)
	})
	mux.HandleFunc("/api/builder", builderCollectionHandler(opts))
	mux.HandleFunc("/api/builder/", builderPrimitiveHandler(opts))
	return mux
}

func builderCollectionHandler(opts GUIOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/builder" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		inventory, err := ListBuilderInventory(BuilderOptions{ProjectRoot: opts.ProjectRoot, WorkspaceRoot: opts.WorkspaceRoot})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, http.StatusOK, inventory)
	}
}

func builderPrimitiveHandler(opts GUIOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/builder/")
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) == 0 || parts[0] == "" {
			http.NotFound(w, r)
			return
		}
		kind := parts[0]
		id := ""
		if len(parts) > 1 {
			id = parts[1]
		}
		builderOpts := BuilderOptions{ProjectRoot: opts.ProjectRoot, WorkspaceRoot: opts.WorkspaceRoot}
		switch r.Method {
		case http.MethodGet:
			handleBuilderGet(w, builderOpts, kind, id)
		case http.MethodPost:
			if id != "" {
				http.Error(w, "POST targets a primitive collection", http.StatusBadRequest)
				return
			}
			handleBuilderWrite(w, r, builderOpts, kind, "", false)
		case http.MethodPut:
			if id == "" {
				http.Error(w, "PUT requires a primitive id", http.StatusBadRequest)
				return
			}
			handleBuilderWrite(w, r, builderOpts, kind, id, true)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleBuilderGet(w http.ResponseWriter, opts BuilderOptions, kind, id string) {
	inventory, err := ListBuilderInventory(opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if id == "" {
		switch kind {
		case "circuits":
			writeJSONResponse(w, http.StatusOK, inventory.Circuits)
		case "machines":
			writeJSONResponse(w, http.StatusOK, inventory.Machines)
		case "automations":
			writeJSONResponse(w, http.StatusOK, inventory.Automations)
		case "commands":
			writeJSONResponse(w, http.StatusOK, inventory.Commands)
		case "bindings":
			writeJSONResponse(w, http.StatusOK, inventory.Bindings)
		default:
			http.Error(w, "unsupported builder primitive", http.StatusNotFound)
		}
		return
	}
	if kind == "bindings" {
		index, err := strconv.Atoi(id)
		if err != nil || index < 0 || index >= len(inventory.Bindings) {
			http.Error(w, "binding not found", http.StatusNotFound)
			return
		}
		writeJSONResponse(w, http.StatusOK, inventory.Bindings[index])
		return
	}
	value, err := ReadBuilderPrimitive(opts, kind, id)
	if errors.Is(err, os.ErrNotExist) {
		http.Error(w, "primitive not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, value)
}

func handleBuilderWrite(w http.ResponseWriter, r *http.Request, opts BuilderOptions, kind, id string, update bool) {
	switch kind {
	case "circuits":
		var input BuilderCircuitInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid circuit JSON", http.StatusBadRequest)
			return
		}
		if update {
			input.ID = id
		}
		value, err := CreateOrUpdateCircuit(opts, input)
		writeBuilderResult(w, value, err)
	case "machines":
		var input MachineDefinition
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid machine JSON", http.StatusBadRequest)
			return
		}
		if update {
			input.ID = id
		}
		value, err := CreateOrUpdateMachine(opts, input)
		writeBuilderResult(w, value, err)
	case "automations":
		var input AutomationDefinition
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid automation JSON", http.StatusBadRequest)
			return
		}
		if update {
			input.ID = id
		}
		value, err := CreateOrUpdateAutomation(opts, input)
		writeBuilderResult(w, value, err)
	case "commands":
		var input BuilderCommandInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid command JSON", http.StatusBadRequest)
			return
		}
		if update {
			input.Name = id
		}
		value, err := CreateOrUpdateCommand(opts, input)
		writeBuilderResult(w, value, err)
	case "bindings":
		var input EventBinding
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid binding JSON", http.StatusBadRequest)
			return
		}
		if update {
			index, err := strconv.Atoi(id)
			if err != nil {
				http.Error(w, "binding id must be a numeric index", http.StatusBadRequest)
				return
			}
			value, err := UpdateEventBinding(opts, index, input)
			writeBuilderResult(w, value, err)
			return
		}
		value, err := CreateEventBinding(opts, input)
		writeBuilderResult(w, value, err)
	default:
		http.Error(w, "unsupported builder primitive", http.StatusNotFound)
	}
}

func writeBuilderResult(w http.ResponseWriter, value interface{}, err error) {
	if errors.Is(err, os.ErrNotExist) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONResponse(w, http.StatusOK, value)
}

func writeJSONResponse(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func RenderFactoryGraphHTML(graph FactoryGraph) (string, error) {
	return RenderFactoryGUIHTML(graph, nil)
}

func RenderFactoryGUIHTML(graph FactoryGraph, runs []RunSession) (string, error) {
	if runs == nil {
		runs = []RunSession{}
	}
	graphJSON, err := json.Marshal(graph)
	if err != nil {
		return "", err
	}
	runsJSON, err := json.Marshal(runs)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	err = guiTemplate.Execute(&out, struct {
		Graph     FactoryGraph
		GraphJSON template.JS
		RunsJSON  template.JS
	}{
		Graph:     graph,
		GraphJSON: template.JS(graphJSON),
		RunsJSON:  template.JS(runsJSON),
	})
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

var guiTemplate = template.Must(template.New("factory-gui").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Factory Map</title>
<style>
:root {
  color-scheme: light;
  --bg: oklch(97.5% 0.009 105);
  --surface: oklch(99% 0.006 105);
  --panel: oklch(94.5% 0.012 110);
  --ink: oklch(23% 0.025 120);
  --muted: oklch(47% 0.025 115);
  --line: oklch(84% 0.025 115);
  --accent: oklch(54% 0.13 150);
  --accent-soft: oklch(91% 0.055 150);
  --warn: oklch(72% 0.13 78);
  --done: oklch(63% 0.12 152);
  --shadow: 0 18px 50px color-mix(in oklch, var(--ink) 13%, transparent);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif;
}
* { box-sizing: border-box; }
body {
  margin: 0;
  background: var(--bg);
  color: var(--ink);
}
button, code { font: inherit; }
.shell {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
}
.sidebar {
  background: var(--surface);
  border-right: 1px solid var(--line);
  padding: 22px;
}
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 24px;
}
.mark {
  width: 30px;
  height: 30px;
  border-radius: 7px;
  background:
    linear-gradient(135deg, transparent 45%, color-mix(in oklch, var(--accent) 70%, var(--surface)) 45%),
    var(--accent-soft);
  border: 1px solid color-mix(in oklch, var(--accent) 45%, var(--line));
}
.brand h1 {
  font-size: 1rem;
  margin: 0;
  letter-spacing: 0;
}
.meta {
  display: grid;
  gap: 12px;
  margin-bottom: 28px;
}
.meta div {
  padding-bottom: 12px;
  border-bottom: 1px solid var(--line);
}
.meta span {
  display: block;
  color: var(--muted);
  font-size: .78rem;
  margin-bottom: 3px;
}
.meta strong {
  font-size: .92rem;
  overflow-wrap: anywhere;
}
.legend {
  display: grid;
  gap: 9px;
  font-size: .86rem;
}
.legend-row {
  display: flex;
  align-items: center;
  gap: 9px;
}
.dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: var(--accent);
}
.dot.automation { background: oklch(58% 0.12 275); }
.dot.circuit { background: oklch(61% 0.13 36); }
.dot.command { background: oklch(54% 0.11 15); }
.dot.sector { background: oklch(60% 0.10 205); }
.main {
  padding: 24px;
  display: grid;
  grid-template-rows: auto minmax(0, 1fr);
  gap: 18px;
}
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}
.toolbar h2 {
  margin: 0;
  font-size: 1.35rem;
  letter-spacing: 0;
}
.toolbar p {
  margin: 4px 0 0;
  max-width: 72ch;
  color: var(--muted);
  font-size: .94rem;
}
.map {
  min-height: 620px;
  overflow: auto;
  border: 1px solid var(--line);
  background: color-mix(in oklch, var(--surface) 72%, var(--panel));
  box-shadow: var(--shadow);
  border-radius: 8px;
  padding: 18px;
}
.lanes {
  display: grid;
  grid-template-columns: repeat(5, minmax(210px, 1fr));
  gap: 14px;
  min-width: 1100px;
}
.lane {
  min-height: 560px;
  background: color-mix(in oklch, var(--surface) 82%, var(--panel));
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 14px;
}
.lane h3 {
  margin: 0 0 12px;
  font-size: .88rem;
  color: var(--muted);
  font-weight: 700;
}
.node {
  display: grid;
  gap: 8px;
  padding: 12px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--surface);
  margin-bottom: 10px;
}
.node header {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: space-between;
}
.node-title {
  font-weight: 700;
  font-size: .94rem;
  overflow-wrap: anywhere;
}
.badge {
  font-size: .7rem;
  text-transform: uppercase;
  letter-spacing: .04em;
  color: var(--muted);
}
.node p {
  margin: 0;
  color: var(--muted);
  font-size: .82rem;
  line-height: 1.4;
}
.run {
  border: 1px solid color-mix(in oklch, var(--accent) 50%, var(--line));
  background: var(--accent-soft);
  color: var(--ink);
  padding: 7px 9px;
  border-radius: 7px;
  cursor: pointer;
  text-align: left;
  overflow-wrap: anywhere;
}
.run:hover { border-color: var(--accent); }
.run:focus-visible { outline: 3px solid color-mix(in oklch, var(--accent) 30%, transparent); outline-offset: 2px; }
.builder {
  display: grid;
  gap: 12px;
  border: 1px solid var(--line);
  background: var(--surface);
  border-radius: 8px;
  padding: 14px;
}
.builder-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}
.builder-form {
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 12px;
  display: grid;
  gap: 8px;
}
.builder-form h3 {
  margin: 0;
  font-size: .9rem;
}
.builder-form label {
  display: grid;
  gap: 5px;
  color: var(--muted);
  font-size: .78rem;
}
.builder-form input,
.builder-form textarea {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 7px;
  padding: 7px;
  background: var(--surface);
  color: var(--ink);
}
.builder-form textarea { min-height: 72px; resize: vertical; }
.builder-form button {
  border: 1px solid color-mix(in oklch, var(--accent) 50%, var(--line));
  border-radius: 7px;
  background: var(--accent-soft);
  color: var(--ink);
  padding: 8px;
  cursor: pointer;
}
.builder-message {
  min-height: 18px;
  color: var(--muted);
  font-size: .8rem;
  overflow-wrap: anywhere;
}
.links {
  border-top: 1px solid var(--line);
  padding-top: 8px;
  color: var(--muted);
  font-size: .76rem;
  line-height: 1.45;
}
.empty {
  border: 1px dashed var(--line);
  border-radius: 8px;
  padding: 12px;
  color: var(--muted);
  font-size: .84rem;
}
.run-monitor {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 12px;
}
.run-card {
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--surface);
  padding: 12px;
  display: grid;
  gap: 8px;
}
.run-card[data-status="needs-human-action"] {
  border-color: color-mix(in oklch, var(--warn) 70%, var(--line));
  background: color-mix(in oklch, var(--warn) 18%, var(--surface));
}
.status-pill {
  display: inline-flex;
  width: fit-content;
  border-radius: 999px;
  border: 1px solid var(--line);
  padding: 3px 8px;
  font-size: .72rem;
  color: var(--muted);
}
.decision {
  border-top: 1px solid var(--line);
  padding-top: 8px;
  display: grid;
  gap: 8px;
}
.decision-control {
  display: grid;
  gap: 5px;
}
.decision-control label {
  font-size: .78rem;
  color: var(--muted);
}
.decision-control input,
.decision-control textarea,
.decision-control select {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 7px;
  padding: 7px;
  background: var(--surface);
  color: var(--ink);
}
@media (max-width: 860px) {
  .shell { grid-template-columns: 1fr; }
  .sidebar { border-right: 0; border-bottom: 1px solid var(--line); }
  .main { padding: 14px; }
}
</style>
</head>
<body>
<div class="shell">
  <aside class="sidebar">
    <div class="brand"><div class="mark" aria-hidden="true"></div><h1>Factory Map</h1></div>
    <div class="meta">
      <div><span>Workspace</span><strong>{{.Graph.WorkspaceRoot}}</strong></div>
      <div><span>Pieces</span><strong id="piece-count">0</strong></div>
      <div><span>Relations</span><strong>{{len .Graph.Edges}}</strong></div>
    </div>
    <div class="legend" aria-label="Map legend">
      <div class="legend-row"><span class="dot automation"></span>Automations</div>
      <div class="legend-row"><span class="dot"></span>Machines</div>
      <div class="legend-row"><span class="dot circuit"></span>Circuits</div>
      <div class="legend-row"><span class="dot command"></span>Commands</div>
      <div class="legend-row"><span class="dot sector"></span>Sectors</div>
    </div>
  </aside>
  <main class="main">
    <section class="toolbar">
      <div>
        <h2>Production System</h2>
        <p>Inspect the current factory as lanes, reusable runtime pieces, and dispatch paths. Run buttons copy explicit CLI commands so execution stays visible.</p>
      </div>
    </section>
    <section class="map" aria-label="Factory production map">
      <div class="lanes" id="lanes"></div>
    </section>
    <section class="builder" aria-label="Factory builder">
      <div class="toolbar">
        <div>
          <h2>Factory Builder</h2>
          <p>Create concrete circuits, machines, automations, commands, and event bindings through structured controls.</p>
        </div>
      </div>
      <div class="builder-grid">
        <form class="builder-form" data-kind="circuits">
          <h3>Circuit</h3>
          <label>ID<input name="id" placeholder="plan-work"></label>
          <label>Name<input name="name" placeholder="Plan work"></label>
          <label>Description<textarea name="description"></textarea></label>
          <button type="submit">Create circuit</button>
          <div class="builder-message"></div>
        </form>
        <form class="builder-form" data-kind="machines">
          <h3>Machine</h3>
          <label>ID<input name="id" placeholder="build-package"></label>
          <label>Name<input name="name" placeholder="Build package"></label>
          <label>Circuits<input name="circuits" placeholder="plan-work, implement-work"></label>
          <button type="submit">Create machine</button>
          <div class="builder-message"></div>
        </form>
        <form class="builder-form" data-kind="automations">
          <h3>Automation</h3>
          <label>ID<input name="id" placeholder="process-package"></label>
          <label>Name<input name="name" placeholder="Process package"></label>
          <label>Machines<input name="machines" placeholder="build-package, release-package"></label>
          <button type="submit">Create automation</button>
          <div class="builder-message"></div>
        </form>
        <form class="builder-form" data-kind="commands">
          <h3>Command</h3>
          <label>Name<input name="name" placeholder="load-intake"></label>
          <label>Description<textarea name="description"></textarea></label>
          <button type="submit">Create command</button>
          <div class="builder-message"></div>
        </form>
        <form class="builder-form" data-kind="bindings">
          <h3>Event Binding</h3>
          <label>On<input name="on" placeholder="command.completed"></label>
          <label>Source command<input name="source" placeholder="load-intake"></label>
          <label>Run command<input name="run" placeholder="drain-work-package"></label>
          <button type="submit">Create binding</button>
          <div class="builder-message"></div>
        </form>
      </div>
    </section>
    <section aria-label="Run monitor">
      <div class="toolbar">
        <div>
          <h2>Run Monitor</h2>
          <p>Watch durable run state and respond when the factory needs human input.</p>
        </div>
      </div>
      <div class="run-monitor" id="run-monitor"></div>
    </section>
  </main>
</div>
<script>
const graph = {{.GraphJSON}};
const runs = {{.RunsJSON}};
const laneOrder = ["dock", "yard", "shop-floor", "finished-goods", "sectors"];
const laneTargets = new Map(laneOrder.map((id) => [id, []]));
const byId = new Map(graph.nodes.map((node) => [node.id, node]));
const children = new Map();
for (const edge of graph.edges) {
  if (!children.has(edge.from)) children.set(edge.from, []);
  children.get(edge.from).push(edge);
  if (edge.type === "contains" && laneTargets.has(edge.from)) {
    laneTargets.get(edge.from).push(byId.get(edge.to));
  }
}
document.getElementById("piece-count").textContent = String(graph.nodes.filter((node) => node.type !== "lane").length);
const lanesEl = document.getElementById("lanes");
for (const laneId of laneOrder) {
  const lane = byId.get(laneId);
  const section = document.createElement("section");
  section.className = "lane";
  const title = document.createElement("h3");
  title.textContent = lane ? lane.label : laneId;
  section.append(title);
  const nodes = laneTargets.get(laneId).filter(Boolean);
  if (nodes.length === 0) {
    const empty = document.createElement("div");
    empty.className = "empty";
    empty.textContent = "No visible pieces in this lane yet.";
    section.append(empty);
  }
  for (const node of nodes) section.append(renderNode(node));
  lanesEl.append(section);
}
wireBuilderForms();
renderRuns();
function renderNode(node) {
  const item = document.createElement("article");
  item.className = "node";
  const header = document.createElement("header");
  const title = document.createElement("div");
  title.className = "node-title";
  title.textContent = node.label;
  const badge = document.createElement("span");
  badge.className = "badge";
  badge.textContent = node.type;
  header.append(title, badge);
  item.append(header);
  if (node.description) {
    const desc = document.createElement("p");
    desc.textContent = node.description;
    item.append(desc);
  }
  const command = runCommandFor(node);
  if (command) {
    const button = document.createElement("button");
    button.className = "run";
    button.type = "button";
    button.textContent = command;
    button.title = "Copy run command";
    button.addEventListener("click", async () => {
      await navigator.clipboard.writeText(command);
      button.textContent = "Copied: " + command;
      window.setTimeout(() => { button.textContent = command; }, 1600);
    });
    item.append(button);
  }
  const outgoing = (children.get(node.id) || []).filter((edge) => edge.type !== "contains");
  if (outgoing.length) {
    const links = document.createElement("div");
    links.className = "links";
    links.textContent = outgoing.map((edge) => edge.label + " " + ((byId.get(edge.to) || {}).label || edge.to)).join(" | ");
    item.append(links);
  }
  return item;
}
function renderRuns() {
  const monitor = document.getElementById("run-monitor");
  if (!runs.length) {
    const empty = document.createElement("div");
    empty.className = "empty";
    empty.textContent = "No durable runs recorded yet.";
    monitor.append(empty);
    return;
  }
  for (const run of runs) {
    const card = document.createElement("article");
    card.className = "run-card";
    card.dataset.status = run.status || "";
    const title = document.createElement("div");
    title.className = "node-title";
    title.textContent = run.targetId || run.runId;
    const status = document.createElement("span");
    status.className = "status-pill";
    status.textContent = run.status || "unknown";
    const summary = document.createElement("p");
    summary.textContent = (run.steps || []).length + " step(s) recorded";
    card.append(title, status, summary);
    if (run.decisionRequest) card.append(renderDecision(run.decisionRequest));
    monitor.append(card);
  }
}
function renderDecision(decision) {
  const wrap = document.createElement("div");
  wrap.className = "decision";
  const title = document.createElement("strong");
  title.textContent = decision.title || "Human action required";
  wrap.append(title);
  if (decision.description) {
    const desc = document.createElement("p");
    desc.textContent = decision.description;
    wrap.append(desc);
  }
  for (const control of decision.controls || []) {
    wrap.append(renderDecisionControl(control));
  }
  return wrap;
}
function renderDecisionControl(control) {
  const field = document.createElement("div");
  field.className = "decision-control";
  const label = document.createElement("label");
  label.textContent = control.label || control.name || control.type;
  field.append(label);
  if (control.type === "option-list") {
    const select = document.createElement("select");
    for (const option of control.options || []) {
      const item = document.createElement("option");
      item.value = option.value;
      item.textContent = option.label || option.value;
      select.append(item);
    }
    field.append(select);
  } else if (control.type === "textarea") {
    field.append(document.createElement("textarea"));
  } else if (control.type === "checkbox") {
    const input = document.createElement("input");
    input.type = "checkbox";
    field.append(input);
  } else {
    const input = document.createElement("input");
    input.type = control.type === "button" ? "button" : "text";
    input.value = control.type === "button" ? (control.label || "Select") : "";
    field.append(input);
  }
  return field;
}
function runCommandFor(node) {
  const id = node.id.replace(/^(command|machine|automation)-/, "");
  if (node.type === "command") return "factory run " + id;
  if (node.type === "machine") return "factory run machine " + id;
  if (node.type === "automation") return "factory run automation " + id;
  return "";
}
function wireBuilderForms() {
  for (const form of document.querySelectorAll(".builder-form")) {
    form.addEventListener("submit", async (event) => {
      event.preventDefault();
      const message = form.querySelector(".builder-message");
      message.textContent = "Saving...";
      const kind = form.dataset.kind;
      const body = builderPayload(kind, new FormData(form));
      try {
        const response = await fetch("/api/builder/" + kind, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(body)
        });
        const text = await response.text();
        if (!response.ok) {
          message.textContent = text.trim() || "Validation failed.";
          return;
        }
        message.textContent = "Saved. Refresh to update the map.";
        form.reset();
      } catch (error) {
        message.textContent = String(error);
      }
    });
  }
}
function builderPayload(kind, data) {
  const values = Object.fromEntries(data.entries());
  if (kind === "circuits") {
    return { id: clean(values.id), name: clean(values.name), description: clean(values.description) };
  }
  if (kind === "machines") {
    return {
      id: clean(values.id),
      name: clean(values.name),
      circuits: splitList(values.circuits).map((circuit) => ({ circuit }))
    };
  }
  if (kind === "automations") {
    return {
      id: clean(values.id),
      name: clean(values.name),
      machines: splitList(values.machines).map((machine) => ({ machine }))
    };
  }
  if (kind === "commands") {
    return { name: clean(values.name), description: clean(values.description) };
  }
  return {
    on: clean(values.on),
    where: clean(values.source) ? { command: clean(values.source) } : {},
    run: clean(values.run),
    mode: "auto"
  };
}
function splitList(value) {
  return String(value || "").split(",").map((item) => item.trim()).filter(Boolean);
}
function clean(value) {
  return String(value || "").trim();
}
</script>
</body>
</html>`))
