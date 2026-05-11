package factory

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
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
		html, err := RenderFactoryGraphHTML(graph)
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
	return mux
}

func RenderFactoryGraphHTML(graph FactoryGraph) (string, error) {
	graphJSON, err := json.Marshal(graph)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	err = guiTemplate.Execute(&out, struct {
		Graph     FactoryGraph
		GraphJSON template.JS
	}{
		Graph:     graph,
		GraphJSON: template.JS(graphJSON),
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
  </main>
</div>
<script>
const graph = {{.GraphJSON}};
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
function runCommandFor(node) {
  const id = node.id.replace(/^(machine|automation)-/, "");
  if (node.type === "machine") return "factory run machine " + id;
  if (node.type === "automation") return "factory run automation " + id;
  return "";
}
</script>
</body>
</html>`))
