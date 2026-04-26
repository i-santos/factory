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

type Runner interface {
	Run(ctx context.Context, command CommandDefinition, input map[string]interface{}) (CommandResult, error)
}

type CodexRunner struct {
	AgentCommand []string
	Env          []string
	WorkDir      string
}

func (r CodexRunner) Run(ctx context.Context, command CommandDefinition, input map[string]interface{}) (CommandResult, error) {
	if mock := os.Getenv("FACTORY_AGENT_MOCK_RESULT"); mock != "" {
		return ParseCommandResult([]byte(mock))
	}
	agentCommand := r.AgentCommand
	if len(agentCommand) == 0 {
		agentCommand = []string{"codex", "exec"}
	}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return CommandResult{}, err
	}
	call := fmt.Sprintf("Factory.call(%q, %s)", command.Name, string(inputJSON))
	args := append(agentCommand[1:], call)
	cmd := exec.CommandContext(ctx, agentCommand[0], args...)
	cmd.Dir = r.WorkDir
	cmd.Env = append(os.Environ(), r.Env...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return CommandResult{}, fmt.Errorf("agent command failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return ParseCommandResult(stdout.Bytes())
}

func ParseCommandResult(data []byte) (CommandResult, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return CommandResult{}, fmt.Errorf("agent returned empty output")
	}
	var result CommandResult
	if err := json.Unmarshal(trimmed, &result); err != nil {
		return CommandResult{}, fmt.Errorf("agent output must be a command result JSON envelope: %w", err)
	}
	result.Raw = append([]byte(nil), trimmed...)
	if result.Status == "" {
		return result, fmt.Errorf("command result missing status")
	}
	switch result.Status {
	case "succeeded", "blocked", "failed", "waiting":
	default:
		return result, fmt.Errorf("unsupported command result status %q", result.Status)
	}
	return result, nil
}

type RunOptions struct {
	ProjectRoot       string
	WorkspaceRoot     string
	Command           string
	Input             map[string]interface{}
	MaxAutoIterations int
	Runner            Runner
}

func RunCommandChain(ctx context.Context, opts RunOptions) ([]RunRecord, error) {
	if opts.WorkspaceRoot == "" {
		opts.WorkspaceRoot = DefaultWorkspaceRoot
	}
	if opts.ProjectRoot == "" {
		opts.ProjectRoot = "."
	}
	cfg, err := LoadConfig(opts.ProjectRoot, opts.WorkspaceRoot)
	if err != nil {
		return nil, err
	}
	limit := opts.MaxAutoIterations
	if limit <= 0 {
		limit = cfg.Events.MaxAutoIterations
	}
	if limit <= 0 {
		limit = 25
	}
	runner := opts.Runner
	if runner == nil {
		runner = CodexRunner{WorkDir: opts.ProjectRoot}
	}
	currentCommand := opts.Command
	currentInput := cloneMap(opts.Input)
	records := []RunRecord{}
	for iteration := 0; iteration < limit; iteration++ {
		reg, err := LoadRegistry(opts.ProjectRoot, cfg)
		if err != nil {
			return records, err
		}
		def, ok := ResolveCommand(reg, currentCommand)
		if !ok {
			return records, fmt.Errorf("command %q is not registered", currentCommand)
		}
		result, err := runner.Run(ctx, def, currentInput)
		if err != nil {
			return records, err
		}
		event := BuildCommandEvent(def.Name, result)
		if err := AppendEvent(opts.ProjectRoot, cfg, event); err != nil {
			return records, err
		}
		record := RunRecord{Command: def.Name, Input: currentInput, Result: result, Event: event}
		records = append(records, record)
		if result.Status != "succeeded" {
			return records, nil
		}
		next, nextInput, ok, err := NextCommandFromEvent(opts.ProjectRoot, cfg, event, result)
		if err != nil {
			return records, err
		}
		if !ok {
			return records, nil
		}
		currentCommand = next
		currentInput = nextInput
	}
	return records, fmt.Errorf("auto iteration limit reached after %d iteration(s)", limit)
}

func BuildCommandEvent(command string, result CommandResult) FactoryEvent {
	eventType := map[string]string{
		"succeeded": "command.completed",
		"blocked":   "command.blocked",
		"failed":    "command.failed",
		"waiting":   "command.waiting",
	}[result.Status]
	data := cloneMap(result.Data)
	if data == nil {
		data = map[string]interface{}{}
	}
	data["resultType"] = result.ResultType
	return FactoryEvent{
		EventID:      fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		EventType:    eventType,
		FactoryID:    "default",
		Command:      command,
		Status:       result.Status,
		ArtifactRefs: append([]string(nil), result.Artifacts...),
		Data:         data,
	}
}

func AppendEvent(projectRoot string, cfg Config, event FactoryEvent) error {
	path := filepath.Join(projectRoot, cfg.Events.LogDir, "events.jsonl")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(data, '\n'))
	return err
}

func NextCommandFromEvent(projectRoot string, cfg Config, event FactoryEvent, result CommandResult) (string, map[string]interface{}, bool, error) {
	file, err := LoadEventBindings(projectRoot, cfg)
	if err != nil {
		return "", nil, false, err
	}
	for _, binding := range file.Bindings {
		if binding.Mode != "" && binding.Mode != "auto" {
			continue
		}
		if binding.On != event.EventType {
			continue
		}
		if !bindingMatches(binding, event, result) {
			continue
		}
		input, err := renderBindingInput(binding.Input, event, result)
		if err != nil {
			return "", nil, false, err
		}
		return binding.Run, input, true, nil
	}
	return "", nil, false, nil
}

func bindingMatches(binding EventBinding, event FactoryEvent, result CommandResult) bool {
	for key, want := range binding.Where {
		switch key {
		case "command":
			if event.Command != want {
				return false
			}
		case "status", "result.status":
			if result.Status != want {
				return false
			}
		case "event.status":
			if event.Status != want {
				return false
			}
		default:
			if strings.HasPrefix(key, "result.data.") {
				if fmt.Sprint(result.Data[strings.TrimPrefix(key, "result.data.")]) != want {
					return false
				}
			} else if strings.HasPrefix(key, "event.data.") {
				if fmt.Sprint(event.Data[strings.TrimPrefix(key, "event.data.")]) != want {
					return false
				}
			}
		}
	}
	return true
}

func renderBindingInput(template map[string]interface{}, event FactoryEvent, result CommandResult) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	for key, value := range template {
		if str, ok := value.(string); ok && strings.HasPrefix(str, "$") {
			resolved, err := resolveBindingValue(str, event, result)
			if err != nil {
				return nil, err
			}
			out[key] = resolved
			continue
		}
		out[key] = value
	}
	return out, nil
}

func resolveBindingValue(expr string, event FactoryEvent, result CommandResult) (interface{}, error) {
	switch {
	case strings.HasPrefix(expr, "$event.data."):
		return event.Data[strings.TrimPrefix(expr, "$event.data.")], nil
	case strings.HasPrefix(expr, "$result.data."):
		return result.Data[strings.TrimPrefix(expr, "$result.data.")], nil
	case expr == "$event.command":
		return event.Command, nil
	case expr == "$result.status":
		return result.Status, nil
	default:
		return nil, fmt.Errorf("unsupported binding expression %q", expr)
	}
}

func cloneMap(in map[string]interface{}) map[string]interface{} {
	if in == nil {
		return map[string]interface{}{}
	}
	out := make(map[string]interface{}, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
