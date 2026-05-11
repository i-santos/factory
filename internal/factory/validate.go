package factory

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ValidationResult struct {
	Status string              `json:"status"`
	Issues []ValidationIssue   `json:"issues,omitempty"`
	Pieces ValidationInventory `json:"pieces"`
}

type ValidationIssue struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type ValidationInventory struct {
	Circuits    int `json:"circuits"`
	Machines    int `json:"machines"`
	Automations int `json:"automations"`
	Sectors     int `json:"sectors"`
}

func ValidateWorkspace(projectRoot, workspaceRoot string) (ValidationResult, error) {
	if projectRoot == "" {
		projectRoot = "."
	}
	if workspaceRoot == "" {
		workspaceRoot = DefaultWorkspaceRoot
	}
	cfg, err := LoadConfig(projectRoot, workspaceRoot)
	if err != nil {
		return ValidationResult{}, err
	}
	result := ValidationResult{Status: "valid"}
	circuits, err := loadAllCircuits(projectRoot, cfg.RootDir)
	if err != nil {
		return result, err
	}
	machines, err := loadAllMachines(projectRoot, cfg.RootDir)
	if err != nil {
		return result, err
	}
	automations, err := loadAllAutomations(projectRoot, cfg.RootDir)
	if err != nil {
		return result, err
	}
	sectors, err := discoverSectors(projectRoot, cfg)
	if err != nil {
		return result, err
	}
	result.Pieces = ValidationInventory{
		Circuits:    len(circuits),
		Machines:    len(machines),
		Automations: len(automations),
		Sectors:     len(sectors),
	}
	for id, circuit := range circuits {
		if strings.TrimSpace(circuit.Program) == "" {
			result.Issues = append(result.Issues, ValidationIssue{Path: circuitPath(cfg.RootDir, id), Message: "circuit program is required"})
			continue
		}
		if path, ok := projectResourcePath(projectRoot, circuit.Program); ok {
			if _, err := os.Stat(path); err != nil {
				result.Issues = append(result.Issues, ValidationIssue{Path: circuit.Program, Message: "circuit program file is missing"})
			}
		}
	}
	for id, machine := range machines {
		for _, step := range machine.Circuits {
			if _, ok := circuits[step.Circuit]; !ok {
				result.Issues = append(result.Issues, ValidationIssue{Path: machinePath(cfg.RootDir, id), Message: fmt.Sprintf("machine references missing circuit %q", step.Circuit)})
			}
		}
	}
	for id, automation := range automations {
		for _, step := range automation.Machines {
			if _, ok := machines[step.Machine]; !ok {
				result.Issues = append(result.Issues, ValidationIssue{Path: automationPath(cfg.RootDir, id), Message: fmt.Sprintf("automation references missing machine %q", step.Machine)})
			}
		}
	}
	if len(result.Issues) > 0 {
		result.Status = "invalid"
	}
	sort.Slice(result.Issues, func(i, j int) bool {
		return result.Issues[i].Path+result.Issues[i].Message < result.Issues[j].Path+result.Issues[j].Message
	})
	return result, nil
}

func loadAllCircuits(projectRoot, workspaceRoot string) (map[string]CircuitDefinition, error) {
	out := map[string]CircuitDefinition{}
	root := filepath.Join(projectRoot, workspaceRoot, "circuits")
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return out, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		def, err := LoadCircuit(projectRoot, workspaceRoot, entry.Name())
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return out, err
		}
		out[def.ID] = def
	}
	return out, nil
}

func loadAllMachines(projectRoot, workspaceRoot string) (map[string]MachineDefinition, error) {
	out := map[string]MachineDefinition{}
	root := filepath.Join(projectRoot, workspaceRoot, "machines")
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return out, err
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		def, err := LoadMachine(projectRoot, workspaceRoot, id)
		if err != nil {
			return out, err
		}
		out[def.ID] = def
	}
	return out, nil
}

func loadAllAutomations(projectRoot, workspaceRoot string) (map[string]AutomationDefinition, error) {
	out := map[string]AutomationDefinition{}
	root := filepath.Join(projectRoot, workspaceRoot, "automations")
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return out, err
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		def, err := LoadAutomation(projectRoot, workspaceRoot, id)
		if err != nil {
			return out, err
		}
		out[def.ID] = def
	}
	return out, nil
}

func circuitPath(workspaceRoot, id string) string {
	return filepath.ToSlash(filepath.Join(workspaceRoot, "circuits", id, "circuit.json"))
}

func machinePath(workspaceRoot, id string) string {
	return filepath.ToSlash(filepath.Join(workspaceRoot, "machines", id+".json"))
}

func automationPath(workspaceRoot, id string) string {
	return filepath.ToSlash(filepath.Join(workspaceRoot, "automations", id+".json"))
}
