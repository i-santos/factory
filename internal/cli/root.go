package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"factory/internal/factory"
	"github.com/spf13/cobra"
)

type rootOptions struct {
	projectRoot   string
	workspaceRoot string
}

func NewRootCommand() *cobra.Command {
	opts := &rootOptions{projectRoot: ".", workspaceRoot: factory.DefaultWorkspaceRoot}
	cmd := &cobra.Command{
		Use:   "factory",
		Short: "Create and run configurable agent factories",
	}
	cmd.PersistentFlags().StringVar(&opts.projectRoot, "project-root", ".", "Project root containing the factory workspace")
	cmd.PersistentFlags().StringVar(&opts.workspaceRoot, "workspace-root", factory.DefaultWorkspaceRoot, "Factory workspace root")
	cmd.AddCommand(newInitCommand(opts))
	cmd.AddCommand(newCommandCommand(opts))
	cmd.AddCommand(newEventCommand(opts))
	cmd.AddCommand(newRunCommand(opts))
	cmd.AddCommand(newVisualizeCommand(opts))
	return cmd
}

func newInitCommand(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a project-local factory workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := factory.InitWorkspace(opts.projectRoot, opts.workspaceRoot); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "initialized %s\n", opts.workspaceRoot)
			return nil
		},
	}
}

func newCommandCommand(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "command",
		Short: "Manage project-created factory commands",
	}
	cmd.AddCommand(newCommandCreateCommand(opts))
	return cmd
}

func newCommandCreateCommand(opts *rootOptions) *cobra.Command {
	var aliases []string
	var description string
	var prompt string
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create and register a project command",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := factory.LoadConfig(opts.projectRoot, opts.workspaceRoot)
			if err != nil {
				return err
			}
			def := factory.CommandDefinition{
				Name:        args[0],
				Prompt:      prompt,
				Aliases:     aliases,
				Description: description,
			}
			created, err := factory.CreateCommand(opts.projectRoot, cfg, def, defaultPrompt(args[0]))
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(created)
		},
	}
	cmd.Flags().StringSliceVar(&aliases, "alias", nil, "Command alias; may be repeated or comma-separated")
	cmd.Flags().StringVar(&description, "description", "", "Human-facing command description")
	cmd.Flags().StringVar(&prompt, "prompt", "", "Existing prompt resource path")
	return cmd
}

func newEventCommand(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event",
		Short: "Manage factory event bindings",
	}
	cmd.AddCommand(newEventBindCommand(opts))
	return cmd
}

func newEventBindCommand(opts *rootOptions) *cobra.Command {
	var on string
	var where []string
	var run string
	var input string
	var mode string
	cmd := &cobra.Command{
		Use:   "bind create",
		Short: "Create an event binding that can trigger another command",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := factory.LoadConfig(opts.projectRoot, opts.workspaceRoot)
			if err != nil {
				return err
			}
			inputMap := map[string]interface{}{}
			if strings.TrimSpace(input) != "" {
				if err := json.Unmarshal([]byte(input), &inputMap); err != nil {
					return fmt.Errorf("--input must be a JSON object: %w", err)
				}
			}
			binding, err := factory.AddEventBinding(opts.projectRoot, cfg, factory.EventBinding{
				On:    on,
				Where: parseWhere(where),
				Run:   run,
				Input: inputMap,
				Mode:  mode,
			})
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(binding)
		},
	}
	cmd.Flags().StringVar(&on, "on", "command.completed", "Event type to match")
	cmd.Flags().StringArrayVar(&where, "where", nil, "Matcher in key=value form; may be repeated")
	cmd.Flags().StringVar(&run, "run", "", "Command to run when the binding matches")
	cmd.Flags().StringVar(&input, "input", "{}", "JSON object used as next command input; supports binding expressions")
	cmd.Flags().StringVar(&mode, "mode", "auto", "Binding mode")
	_ = cmd.MarkFlagRequired("run")
	return cmd
}

func newRunCommand(opts *rootOptions) *cobra.Command {
	var data string
	var maxAutoIterations int
	var agentCommand string
	cmd := &cobra.Command{
		Use:   "run <command>",
		Short: "Run a registered command and follow matching event bindings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := map[string]interface{}{}
			if strings.TrimSpace(data) != "" {
				if err := json.Unmarshal([]byte(data), &input); err != nil {
					return fmt.Errorf("--data must be a JSON object: %w", err)
				}
			}
			runner := factory.CodexRunner{
				AgentCommand: splitCommand(agentCommand),
				WorkDir:      opts.projectRoot,
			}
			records, err := factory.RunCommandChain(context.Background(), factory.RunOptions{
				ProjectRoot:       opts.projectRoot,
				WorkspaceRoot:     opts.workspaceRoot,
				Command:           args[0],
				Input:             input,
				MaxAutoIterations: maxAutoIterations,
				Runner:            runner,
			})
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(map[string]interface{}{
				"status":  "completed",
				"records": records,
			})
		},
	}
	cmd.Flags().StringVar(&data, "data", "{}", "JSON object passed as command input")
	cmd.Flags().IntVar(&maxAutoIterations, "max-auto-iterations", 0, "Maximum automatic event-binding continuations")
	cmd.Flags().StringVar(&agentCommand, "agent-command", "codex exec", "Agent command used to execute Factory.call")
	return cmd
}

func newVisualizeCommand(opts *rootOptions) *cobra.Command {
	var format string
	var output string
	cmd := &cobra.Command{
		Use:   "visualize",
		Short: "Visualize the current factory workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			graph, err := factory.BuildFactoryGraph(factory.VisualizationOptions{
				ProjectRoot:   opts.projectRoot,
				WorkspaceRoot: opts.workspaceRoot,
			})
			if err != nil {
				return err
			}
			var data []byte
			switch format {
			case "mermaid":
				data = []byte(factory.RenderFactoryGraphMermaid(graph))
			case "json":
				data, err = factory.RenderFactoryGraphJSON(graph)
				if err != nil {
					return err
				}
				data = append(data, '\n')
			default:
				return fmt.Errorf("--format must be mermaid or json")
			}
			if strings.TrimSpace(output) == "" {
				_, err = cmd.OutOrStdout().Write(data)
				return err
			}
			if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
				return err
			}
			return os.WriteFile(output, data, 0o644)
		},
	}
	cmd.Flags().StringVar(&format, "format", "mermaid", "Output format: mermaid or json")
	cmd.Flags().StringVar(&output, "output", "", "Write visualization to a file instead of stdout")
	return cmd
}

func defaultPrompt(name string) string {
	return fmt.Sprintf(`# %s

---
factory_runtime_contract: 1
factory_prompt_contract: project-command.v1
---

<comment>
Goal: define this project-created command.
</comment>

<code>
request = Runtime.input optional

return {
  status: "succeeded",
  resultType: "%s.completed",
  summary: "Command completed.",
  artifacts: [],
  events: [],
  data: {}
}
</code>
`, name, name)
}

func splitCommand(value string) []string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return []string{"codex", "exec"}
	}
	return fields
}

func parseWhere(values []string) map[string]string {
	out := map[string]string{}
	for _, value := range values {
		key, val, ok := strings.Cut(value, "=")
		if !ok {
			continue
		}
		out[strings.TrimSpace(key)] = strings.TrimSpace(val)
	}
	return out
}

func ExecuteForTest(args ...string) error {
	cmd := NewRootCommand()
	cmd.SetArgs(args)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	return cmd.Execute()
}
