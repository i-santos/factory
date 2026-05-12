package factory

import "encoding/json"

type Config struct {
	SchemaVersion int    `json:"schemaVersion"`
	RootDir       string `json:"rootDir"`
	FactorySkill  struct {
		Name                 string `json:"name"`
		RequiredSkillVersion string `json:"requiredSkillVersion"`
		RuntimeContract      int    `json:"runtimeContract"`
		WorkspaceSchema      int    `json:"workspaceSchema"`
		ResourceMode         string `json:"resourceMode"`
	} `json:"factorySkill"`
	Commands struct {
		RegistryPath       string `json:"registryPath"`
		ProjectCommandsDir string `json:"projectCommandsDir"`
	} `json:"commands"`
	Events struct {
		BindingsPath      string `json:"bindingsPath"`
		LogDir            string `json:"logDir"`
		MaxAutoIterations int    `json:"maxAutoIterations"`
	} `json:"events"`
	Migrations struct {
		LogPath    string `json:"logPath"`
		BackupDir  string `json:"backupDir"`
		StagingDir string `json:"stagingDir"`
	} `json:"migrations"`
	Sectors struct {
		RootDir string `json:"rootDir"`
	} `json:"sectors"`
}

type CommandDefinition struct {
	Name         string   `json:"name"`
	Prompt       string   `json:"prompt"`
	Aliases      []string `json:"aliases,omitempty"`
	Description  string   `json:"description,omitempty"`
	InputSchema  string   `json:"inputSchema,omitempty"`
	ResultSchema string   `json:"resultSchema"`
	Emits        []string `json:"emits,omitempty"`
	Consumes     []string `json:"consumes,omitempty"`
}

type CommandRegistry struct {
	SchemaVersion int                          `json:"schemaVersion"`
	Commands      map[string]CommandDefinition `json:"commands"`
}

type EventBinding struct {
	On    string                 `json:"on"`
	Where map[string]string      `json:"where,omitempty"`
	Run   string                 `json:"run"`
	Input map[string]interface{} `json:"input,omitempty"`
	Mode  string                 `json:"mode,omitempty"`
}

type EventBindingsFile struct {
	SchemaVersion int            `json:"schemaVersion"`
	Bindings      []EventBinding `json:"bindings"`
}

type CommandResult struct {
	Status          string                 `json:"status"`
	ResultType      string                 `json:"resultType,omitempty"`
	Summary         string                 `json:"summary,omitempty"`
	Artifacts       []string               `json:"artifacts,omitempty"`
	Events          []FactoryEvent         `json:"events,omitempty"`
	Next            *NextRecommendation    `json:"next,omitempty"`
	NextAction      *NextAction            `json:"nextAction,omitempty"`
	DecisionRequest *DecisionRequest       `json:"decisionRequest,omitempty"`
	MissingEvidence []string               `json:"missingEvidence,omitempty"`
	Errors          []string               `json:"errors,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
	Raw             json.RawMessage        `json:"-"`
}

type NextRecommendation struct {
	Recommendation string                 `json:"recommendation,omitempty"`
	Input          map[string]interface{} `json:"input,omitempty"`
}

type FactoryEvent struct {
	EventID      string                 `json:"eventId"`
	EventType    string                 `json:"eventType"`
	FactoryID    string                 `json:"factoryId,omitempty"`
	Command      string                 `json:"command"`
	Status       string                 `json:"status"`
	ArtifactRefs []string               `json:"artifactRefs,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
}

type RunRecord struct {
	Command string                 `json:"command"`
	Input   map[string]interface{} `json:"input,omitempty"`
	Result  CommandResult          `json:"result"`
	Event   FactoryEvent           `json:"event"`
}

type RunStatus string

const (
	RunStatusIdle             RunStatus = "idle"
	RunStatusRunning          RunStatus = "running"
	RunStatusSucceeded        RunStatus = "succeeded"
	RunStatusBlocked          RunStatus = "blocked"
	RunStatusFailed           RunStatus = "failed"
	RunStatusNeedsHumanAction RunStatus = "needs-human-action"
)

type NextAction struct {
	Kind          string                 `json:"kind"`
	Command       string                 `json:"command,omitempty"`
	Input         map[string]interface{} `json:"input,omitempty"`
	RequiresHuman bool                   `json:"requiresHuman,omitempty"`
}

type DecisionControlOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type DecisionControl struct {
	Type    string                  `json:"type"`
	Name    string                  `json:"name"`
	Label   string                  `json:"label,omitempty"`
	Options []DecisionControlOption `json:"options,omitempty"`
}

type DecisionRequest struct {
	ID           string            `json:"id"`
	Title        string            `json:"title"`
	Description  string            `json:"description,omitempty"`
	Controls     []DecisionControl `json:"controls,omitempty"`
	SubmitAction *NextAction       `json:"submitAction,omitempty"`
}

type RunSession struct {
	RunID           string           `json:"runId"`
	TargetType      string           `json:"targetType"`
	TargetID        string           `json:"targetId"`
	Status          RunStatus        `json:"status"`
	StartedAt       string           `json:"startedAt"`
	FinishedAt      string           `json:"finishedAt,omitempty"`
	Steps           []RunRecord      `json:"steps"`
	NextAction      *NextAction      `json:"nextAction,omitempty"`
	DecisionRequest *DecisionRequest `json:"decisionRequest,omitempty"`
	StopReason      string           `json:"stopReason,omitempty"`
}
