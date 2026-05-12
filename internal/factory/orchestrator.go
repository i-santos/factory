package factory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type CircuitDefinition struct {
	ID            string `json:"id"`
	Name          string `json:"name,omitempty"`
	Program       string `json:"program"`
	RuntimeKernel string `json:"runtimeKernel,omitempty"`
	Description   string `json:"description,omitempty"`
}

type MachineCircuitStep struct {
	Circuit string      `json:"circuit"`
	Input   interface{} `json:"input,omitempty"`
}

type MachineDefinition struct {
	ID          string               `json:"id"`
	Name        string               `json:"name,omitempty"`
	Circuits    []MachineCircuitStep `json:"circuits"`
	Description string               `json:"description,omitempty"`
}

type AutomationMachineStep struct {
	Machine string      `json:"machine"`
	Input   interface{} `json:"input,omitempty"`
}

type AutomationDefinition struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name,omitempty"`
	Machines    []AutomationMachineStep `json:"machines"`
	Description string                  `json:"description,omitempty"`
}

type OrchestratorOptions struct {
	ProjectRoot   string
	WorkspaceRoot string
	Runner        CircuitRunner
}

type CircuitRunner interface {
	RunCircuit(ctx context.Context, circuit CircuitDefinition, input map[string]interface{}) (CircuitResult, error)
}

type CircuitResult struct {
	Status          string                 `json:"status"`
	Summary         string                 `json:"summary,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
	Artifacts       []string               `json:"artifacts,omitempty"`
	Errors          []string               `json:"errors,omitempty"`
	NextAction      *NextAction            `json:"nextAction,omitempty"`
	DecisionRequest *DecisionRequest       `json:"decisionRequest,omitempty"`
	Raw             json.RawMessage        `json:"-"`
}

type MachineRunResult struct {
	Status   string             `json:"status"`
	Machine  string             `json:"machine"`
	Circuits []CircuitRunRecord `json:"circuits"`
	Started  string             `json:"startedAt"`
	Finished string             `json:"finishedAt,omitempty"`
}

type CircuitRunRecord struct {
	Circuit string        `json:"circuit"`
	Input   interface{}   `json:"input,omitempty"`
	Result  CircuitResult `json:"result"`
}

type AutomationRunResult struct {
	Status     string             `json:"status"`
	Automation string             `json:"automation"`
	Machines   []MachineRunResult `json:"machines"`
	Started    string             `json:"startedAt"`
	Finished   string             `json:"finishedAt,omitempty"`
}

type CodexCircuitRunner struct {
	AgentCommand []string
	WorkDir      string
	Env          []string
}

func (r CodexCircuitRunner) RunCircuit(ctx context.Context, circuit CircuitDefinition, input map[string]interface{}) (CircuitResult, error) {
	if mock := os.Getenv("FACTORY_CIRCUIT_MOCK_RESULT"); mock != "" {
		return ParseCircuitResult([]byte(mock))
	}
	agentCommand := r.AgentCommand
	if len(agentCommand) == 0 {
		agentCommand = []string{"codex", "exec"}
	}
	program, err := os.ReadFile(filepath.Join(r.WorkDir, filepath.FromSlash(strings.TrimPrefix(circuit.Program, "project://"))))
	if err != nil {
		return CircuitResult{}, err
	}
	kernel := ""
	if circuit.RuntimeKernel != "" {
		data, err := os.ReadFile(filepath.Join(r.WorkDir, filepath.FromSlash(strings.TrimPrefix(circuit.RuntimeKernel, "project://"))))
		if err != nil {
			return CircuitResult{}, err
		}
		kernel = string(data)
	}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return CircuitResult{}, err
	}
	prompt := fmt.Sprintf("%s\n\nRuntime input:\n```json\n%s\n```\n\n%s", kernel, string(inputJSON), string(program))
	args := append(agentCommand[1:], prompt)
	cmd := exec.CommandContext(ctx, agentCommand[0], args...)
	cmd.Dir = r.WorkDir
	cmd.Env = append(os.Environ(), r.Env...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return CircuitResult{}, fmt.Errorf("circuit command failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return ParseCircuitResult(stdout.Bytes())
}

func ParseCircuitResult(data []byte) (CircuitResult, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return CircuitResult{}, fmt.Errorf("circuit returned empty output")
	}
	var result CircuitResult
	if err := json.Unmarshal(trimmed, &result); err != nil {
		return CircuitResult{}, fmt.Errorf("circuit output must be a JSON result envelope: %w", err)
	}
	result.Raw = append([]byte(nil), trimmed...)
	if result.Status == "" {
		return result, fmt.Errorf("circuit result missing status")
	}
	switch result.Status {
	case "succeeded", "blocked", "failed", "waiting", string(RunStatusNeedsHumanAction):
	default:
		return result, fmt.Errorf("unsupported circuit result status %q", result.Status)
	}
	return result, nil
}

func RunMachine(ctx context.Context, opts OrchestratorOptions, machineID string, input map[string]interface{}) (MachineRunResult, error) {
	opts = normalizeOrchestratorOptions(opts)
	runner := opts.Runner
	if runner == nil {
		runner = CodexCircuitRunner{WorkDir: opts.ProjectRoot}
	}
	machine, err := LoadMachine(opts.ProjectRoot, opts.WorkspaceRoot, machineID)
	if err != nil {
		return MachineRunResult{}, err
	}
	result := MachineRunResult{
		Status:  "succeeded",
		Machine: machine.ID,
		Started: time.Now().UTC().Format(time.RFC3339),
	}
	currentInput := cloneMap(input)
	for _, step := range machine.Circuits {
		circuit, err := LoadCircuit(opts.ProjectRoot, opts.WorkspaceRoot, step.Circuit)
		if err != nil {
			return result, err
		}
		stepInput := resolveStepInput(step.Input, currentInput)
		circuitResult, err := runner.RunCircuit(ctx, circuit, stepInput)
		if err != nil {
			return result, err
		}
		result.Circuits = append(result.Circuits, CircuitRunRecord{
			Circuit: circuit.ID,
			Input:   stepInput,
			Result:  circuitResult,
		})
		if circuitResult.Status != "succeeded" {
			result.Status = circuitResult.Status
			result.Finished = time.Now().UTC().Format(time.RFC3339)
			return result, nil
		}
		currentInput = cloneMap(circuitResult.Data)
	}
	result.Finished = time.Now().UTC().Format(time.RFC3339)
	return result, nil
}

func RunAutomation(ctx context.Context, opts OrchestratorOptions, automationID string, input map[string]interface{}) (AutomationRunResult, error) {
	opts = normalizeOrchestratorOptions(opts)
	automation, err := LoadAutomation(opts.ProjectRoot, opts.WorkspaceRoot, automationID)
	if err != nil {
		return AutomationRunResult{}, err
	}
	result := AutomationRunResult{
		Status:     "succeeded",
		Automation: automation.ID,
		Started:    time.Now().UTC().Format(time.RFC3339),
	}
	currentInput := cloneMap(input)
	for _, step := range automation.Machines {
		stepInput := resolveStepInput(step.Input, currentInput)
		machineResult, err := RunMachine(ctx, opts, step.Machine, stepInput)
		if err != nil {
			return result, err
		}
		result.Machines = append(result.Machines, machineResult)
		if machineResult.Status != "succeeded" {
			result.Status = machineResult.Status
			result.Finished = time.Now().UTC().Format(time.RFC3339)
			return result, nil
		}
		if len(machineResult.Circuits) > 0 {
			currentInput = cloneMap(machineResult.Circuits[len(machineResult.Circuits)-1].Result.Data)
		}
	}
	result.Finished = time.Now().UTC().Format(time.RFC3339)
	return result, nil
}

func LoadCircuit(projectRoot, workspaceRoot, id string) (CircuitDefinition, error) {
	var def CircuitDefinition
	if err := readJSON(filepath.Join(projectRoot, workspaceRoot, "circuits", id, "circuit.json"), &def); err != nil {
		return def, err
	}
	if def.ID == "" {
		def.ID = id
	}
	return def, nil
}

func LoadMachine(projectRoot, workspaceRoot, id string) (MachineDefinition, error) {
	var def MachineDefinition
	if err := readJSON(filepath.Join(projectRoot, workspaceRoot, "machines", id+".json"), &def); err != nil {
		return def, err
	}
	if def.ID == "" {
		def.ID = id
	}
	if len(def.Circuits) == 0 {
		return def, fmt.Errorf("machine %q has no circuits", id)
	}
	return def, nil
}

func LoadAutomation(projectRoot, workspaceRoot, id string) (AutomationDefinition, error) {
	var def AutomationDefinition
	if err := readJSON(filepath.Join(projectRoot, workspaceRoot, "automations", id+".json"), &def); err != nil {
		return def, err
	}
	if def.ID == "" {
		def.ID = id
	}
	if len(def.Machines) == 0 {
		return def, fmt.Errorf("automation %q has no machines", id)
	}
	return def, nil
}

func normalizeOrchestratorOptions(opts OrchestratorOptions) OrchestratorOptions {
	if opts.ProjectRoot == "" {
		opts.ProjectRoot = "."
	}
	if opts.WorkspaceRoot == "" {
		opts.WorkspaceRoot = DefaultWorkspaceRoot
	}
	return opts
}

func resolveStepInput(configured interface{}, previous map[string]interface{}) map[string]interface{} {
	switch value := configured.(type) {
	case nil:
		return cloneMap(previous)
	case string:
		if value == "previous.output" || value == "automation.input" || value == "machine.input" {
			return cloneMap(previous)
		}
		return map[string]interface{}{"value": value}
	case map[string]interface{}:
		return cloneMap(value)
	default:
		return map[string]interface{}{"value": value}
	}
}
