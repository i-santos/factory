package factory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type VisualizationOptions struct {
	ProjectRoot   string
	WorkspaceRoot string
}

type FactoryGraph struct {
	SchemaVersion int         `json:"schemaVersion"`
	WorkspaceRoot string      `json:"workspaceRoot"`
	Nodes         []GraphNode `json:"nodes"`
	Edges         []GraphEdge `json:"edges"`
}

type GraphNode struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Type        string   `json:"type"`
	Path        string   `json:"path,omitempty"`
	Description string   `json:"description,omitempty"`
	Aliases     []string `json:"aliases,omitempty"`
	Consumes    []string `json:"consumes,omitempty"`
	Emits       []string `json:"emits,omitempty"`
}

type GraphEdge struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Type  string `json:"type"`
	Label string `json:"label,omitempty"`
}

func BuildFactoryGraph(opts VisualizationOptions) (FactoryGraph, error) {
	if opts.ProjectRoot == "" {
		opts.ProjectRoot = "."
	}
	if opts.WorkspaceRoot == "" {
		opts.WorkspaceRoot = DefaultWorkspaceRoot
	}
	cfg, err := LoadConfig(opts.ProjectRoot, opts.WorkspaceRoot)
	if err != nil {
		return FactoryGraph{}, err
	}
	reg, err := LoadRegistry(opts.ProjectRoot, cfg)
	if err != nil {
		return FactoryGraph{}, err
	}
	bindings, err := LoadEventBindings(opts.ProjectRoot, cfg)
	if err != nil {
		return FactoryGraph{}, err
	}

	builder := graphBuilder{graph: FactoryGraph{SchemaVersion: 1, WorkspaceRoot: cfg.RootDir}}
	builder.addWorkspaceFlow(cfg.RootDir)
	validation, err := ValidateWorkspace(opts.ProjectRoot, opts.WorkspaceRoot)
	if err != nil {
		return FactoryGraph{}, err
	}
	if validation.Status != "valid" {
		return FactoryGraph{}, fmt.Errorf("workspace validation failed with %d issue(s)", len(validation.Issues))
	}
	if err := builder.addReusablePieces(opts.ProjectRoot, cfg); err != nil {
		return FactoryGraph{}, err
	}
	for _, name := range sortedCommandNames(reg.Commands) {
		def := reg.Commands[name]
		builder.addNode(GraphNode{
			ID:          commandNodeID(def.Name),
			Label:       def.Name,
			Type:        "machine",
			Path:        strings.TrimPrefix(def.Prompt, "project://"),
			Description: def.Description,
			Aliases:     append([]string(nil), def.Aliases...),
			Consumes:    append([]string(nil), def.Consumes...),
			Emits:       append([]string(nil), def.Emits...),
		})
		builder.addEdge("shop-floor", commandNodeID(def.Name), "contains", "machine")
	}
	for _, binding := range bindings.Bindings {
		for _, from := range matchingBindingCommands(reg.Commands, binding) {
			label := binding.On
			if binding.Mode != "" && binding.Mode != "auto" {
				label += " / " + binding.Mode
			}
			builder.addEdge(commandNodeID(from), commandNodeID(binding.Run), "event-binding", label)
		}
	}
	sectors, err := discoverSectors(opts.ProjectRoot, cfg)
	if err != nil {
		return FactoryGraph{}, err
	}
	for _, sector := range sectors {
		builder.addNode(GraphNode{ID: sectorNodeID(sector.Name), Label: sector.Name, Type: "sector", Path: sector.Path})
		builder.addEdge("sectors", sectorNodeID(sector.Name), "contains", "sector")
		for _, action := range sector.Actions {
			actionID := sectorActionNodeID(sector.Name, action.Name)
			builder.addNode(GraphNode{ID: actionID, Label: action.Name, Type: "sector-action", Path: action.Path})
			builder.addEdge(sectorNodeID(sector.Name), actionID, "contains", "action")
		}
	}
	if err := builder.addSectorReferences(opts.ProjectRoot, reg.Commands, sectors); err != nil {
		return FactoryGraph{}, err
	}
	builder.sort()
	return builder.graph, nil
}

func (b *graphBuilder) addReusablePieces(projectRoot string, cfg Config) error {
	circuits, err := loadAllCircuits(projectRoot, cfg.RootDir)
	if err != nil {
		return err
	}
	for _, id := range sortedMapKeys(circuits) {
		circuit := circuits[id]
		b.addNode(GraphNode{
			ID:          circuitNodeID(circuit.ID),
			Label:       displayName(circuit.ID, circuit.Name),
			Type:        "circuit",
			Path:        strings.TrimPrefix(circuit.Program, "project://"),
			Description: circuit.Description,
		})
		b.addEdge("shop-floor", circuitNodeID(circuit.ID), "contains", "circuit")
	}

	machines, err := loadAllMachines(projectRoot, cfg.RootDir)
	if err != nil {
		return err
	}
	for _, id := range sortedMapKeys(machines) {
		machine := machines[id]
		b.addNode(GraphNode{
			ID:          machineNodeID(machine.ID),
			Label:       displayName(machine.ID, machine.Name),
			Type:        "machine",
			Path:        machinePath(cfg.RootDir, id),
			Description: machine.Description,
		})
		b.addEdge("shop-floor", machineNodeID(machine.ID), "contains", "machine")
		for _, step := range machine.Circuits {
			b.addEdge(machineNodeID(machine.ID), circuitNodeID(step.Circuit), "machine-circuit", "runs")
		}
	}

	automations, err := loadAllAutomations(projectRoot, cfg.RootDir)
	if err != nil {
		return err
	}
	for _, id := range sortedMapKeys(automations) {
		automation := automations[id]
		b.addNode(GraphNode{
			ID:          automationNodeID(automation.ID),
			Label:       displayName(automation.ID, automation.Name),
			Type:        "automation",
			Path:        automationPath(cfg.RootDir, id),
			Description: automation.Description,
		})
		b.addEdge("yard", automationNodeID(automation.ID), "contains", "automation")
		for _, step := range automation.Machines {
			b.addEdge(automationNodeID(automation.ID), machineNodeID(step.Machine), "automation-machine", "dispatches")
		}
	}
	return nil
}

func sortedMapKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func RenderFactoryGraphMermaid(graph FactoryGraph) string {
	var out strings.Builder
	out.WriteString("flowchart LR\n")
	out.WriteString("  subgraph workspace[Factory Workspace]\n")
	for _, node := range graph.Nodes {
		if node.Type != "lane" {
			continue
		}
		fmt.Fprintf(&out, "    %s[%q]\n", mermaidID(node.ID), node.Label)
	}
	out.WriteString("  end\n")
	for _, node := range graph.Nodes {
		if node.Type == "lane" {
			continue
		}
		shape := "[%q]"
		if node.Type == "sector" {
			shape = "{{%q}}"
		}
		if node.Type == "sector-action" {
			shape = "([%q])"
		}
		fmt.Fprintf(&out, "  %s"+shape+"\n", mermaidID(node.ID), node.Label)
	}
	for _, edge := range graph.Edges {
		label := edge.Label
		if label == "" {
			label = edge.Type
		}
		fmt.Fprintf(&out, "  %s -->|%s| %s\n", mermaidID(edge.From), escapeMermaidLabel(label), mermaidID(edge.To))
	}
	return out.String()
}

func RenderFactoryGraphJSON(graph FactoryGraph) ([]byte, error) {
	return json.MarshalIndent(graph, "", "  ")
}

type graphBuilder struct {
	graph FactoryGraph
	seen  map[string]bool
	edges map[string]bool
}

func (b *graphBuilder) addNode(node GraphNode) {
	if b.seen == nil {
		b.seen = map[string]bool{}
	}
	if b.seen[node.ID] {
		return
	}
	b.graph.Nodes = append(b.graph.Nodes, node)
	b.seen[node.ID] = true
}

func (b *graphBuilder) addEdge(from, to, edgeType, label string) {
	if b.edges == nil {
		b.edges = map[string]bool{}
	}
	key := from + "\x00" + to + "\x00" + edgeType + "\x00" + label
	if b.edges[key] {
		return
	}
	b.graph.Edges = append(b.graph.Edges, GraphEdge{From: from, To: to, Type: edgeType, Label: label})
	b.edges[key] = true
}

func (b *graphBuilder) addWorkspaceFlow(root string) {
	lanes := []GraphNode{
		{ID: "dock", Label: "01 Dock", Type: "lane", Path: filepath.ToSlash(filepath.Join(root, "01-dock"))},
		{ID: "yard", Label: "02 Yard", Type: "lane", Path: filepath.ToSlash(filepath.Join(root, "02-yard"))},
		{ID: "shop-floor", Label: "03 Shop Floor", Type: "lane", Path: filepath.ToSlash(filepath.Join(root, "03-shop-floor"))},
		{ID: "finished-goods", Label: "04 Finished Goods", Type: "lane", Path: filepath.ToSlash(filepath.Join(root, "04-finished-goods"))},
		{ID: "sectors", Label: "05 Sectors", Type: "lane", Path: filepath.ToSlash(filepath.Join(root, "05-sectors"))},
	}
	for _, lane := range lanes {
		b.addNode(lane)
	}
	b.addEdge("dock", "yard", "workspace-flow", "triage")
	b.addEdge("yard", "shop-floor", "workspace-flow", "package")
	b.addEdge("shop-floor", "finished-goods", "workspace-flow", "ship")
	b.addEdge("shop-floor", "sectors", "workspace-integration", "consult")
}

func (b *graphBuilder) addSectorReferences(projectRoot string, commands map[string]CommandDefinition, sectors []sectorInfo) error {
	for _, def := range commands {
		promptPath, ok := projectResourcePath(projectRoot, def.Prompt)
		if !ok {
			continue
		}
		data, err := os.ReadFile(promptPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		body := string(data)
		for _, sector := range sectors {
			if strings.Contains(body, sector.Name) || strings.Contains(body, sector.Path) {
				b.addEdge(commandNodeID(def.Name), sectorNodeID(sector.Name), "sector-reference", "uses")
			}
			for _, action := range sector.Actions {
				if strings.Contains(body, action.Name) || strings.Contains(body, action.Path) {
					b.addEdge(commandNodeID(def.Name), sectorActionNodeID(sector.Name, action.Name), "sector-action-reference", "calls")
				}
			}
		}
	}
	return nil
}

func (b *graphBuilder) sort() {
	sort.Slice(b.graph.Nodes, func(i, j int) bool {
		if b.graph.Nodes[i].Type == "lane" && b.graph.Nodes[j].Type == "lane" {
			return laneRank(b.graph.Nodes[i].ID) < laneRank(b.graph.Nodes[j].ID)
		}
		if b.graph.Nodes[i].Type == b.graph.Nodes[j].Type {
			return b.graph.Nodes[i].ID < b.graph.Nodes[j].ID
		}
		return nodeTypeRank(b.graph.Nodes[i].Type) < nodeTypeRank(b.graph.Nodes[j].Type)
	})
	sort.Slice(b.graph.Edges, func(i, j int) bool {
		left := b.graph.Edges[i].From + b.graph.Edges[i].To + b.graph.Edges[i].Type + b.graph.Edges[i].Label
		right := b.graph.Edges[j].From + b.graph.Edges[j].To + b.graph.Edges[j].Type + b.graph.Edges[j].Label
		return left < right
	})
}

type sectorInfo struct {
	Name    string
	Path    string
	Actions []sectorActionInfo
}

type sectorActionInfo struct {
	Name string
	Path string
}

func discoverSectors(projectRoot string, cfg Config) ([]sectorInfo, error) {
	root := filepath.Join(projectRoot, cfg.Sectors.RootDir)
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var sectors []sectorInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		sector := sectorInfo{
			Name: entry.Name(),
			Path: filepath.ToSlash(filepath.Join(cfg.Sectors.RootDir, entry.Name())),
		}
		actionsRoot := filepath.Join(root, entry.Name(), "actions")
		actionEntries, err := os.ReadDir(actionsRoot)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		for _, actionEntry := range actionEntries {
			if actionEntry.IsDir() {
				continue
			}
			name := strings.TrimSuffix(actionEntry.Name(), filepath.Ext(actionEntry.Name()))
			sector.Actions = append(sector.Actions, sectorActionInfo{
				Name: name,
				Path: filepath.ToSlash(filepath.Join(cfg.Sectors.RootDir, entry.Name(), "actions", actionEntry.Name())),
			})
		}
		sort.Slice(sector.Actions, func(i, j int) bool { return sector.Actions[i].Name < sector.Actions[j].Name })
		sectors = append(sectors, sector)
	}
	sort.Slice(sectors, func(i, j int) bool { return sectors[i].Name < sectors[j].Name })
	return sectors, nil
}

func sortedCommandNames(commands map[string]CommandDefinition) []string {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func matchingBindingCommands(commands map[string]CommandDefinition, binding EventBinding) []string {
	if command := binding.Where["command"]; command != "" {
		return []string{command}
	}
	names := sortedCommandNames(commands)
	if len(names) == 0 {
		return []string{"unknown"}
	}
	return names
}

func commandNodeID(name string) string {
	return "cmd-" + name
}

func circuitNodeID(name string) string {
	return "circuit-" + name
}

func machineNodeID(name string) string {
	return "machine-" + name
}

func automationNodeID(name string) string {
	return "automation-" + name
}

func displayName(id, name string) string {
	if name != "" {
		return name
	}
	return id
}

func sectorNodeID(name string) string {
	return "sector-" + name
}

func sectorActionNodeID(sector, action string) string {
	return "sector-" + sector + "-action-" + action
}

func mermaidID(id string) string {
	replacer := strings.NewReplacer("-", "_", ".", "_", "/", "_", " ", "_")
	return replacer.Replace(id)
}

func escapeMermaidLabel(label string) string {
	return strings.ReplaceAll(label, "|", "/")
}

func nodeTypeRank(nodeType string) int {
	switch nodeType {
	case "lane":
		return 0
	case "automation":
		return 1
	case "machine":
		return 2
	case "circuit":
		return 3
	case "sector":
		return 4
	case "sector-action":
		return 5
	default:
		return 99
	}
}

func laneRank(id string) int {
	switch id {
	case "dock":
		return 0
	case "yard":
		return 1
	case "shop-floor":
		return 2
	case "finished-goods":
		return 3
	case "sectors":
		return 4
	default:
		return 99
	}
}
