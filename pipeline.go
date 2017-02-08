package gocd

import (
	"fmt"

	"github.com/fatih/structs"
)

type JobStateTransitions struct {
	StateChangeTime int    `json:"state_change_time,omitempty"`
	ID              int    `json:"id,omitempty"`
	State           string `json:"state,omitempty"`
}

type Value struct {
	Value string `json:"value"`
}

type Material struct {
	Description string `json:"description,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Type        string `json:"type,omitempty"`
	ID          int    `json:"id,omitempty"`
}

type MaterialGitConfig struct {
	Type       string `json:"type"`
	Attributes struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		Branch      string `json:"branch"`
		Destination string `json:"destination"`
		AutoUpdate  bool   `json:"auto_update"`
		Filter      struct {
			Ignore []string `json:"ignore"`
		} `json:"filter"`
		InvertFilter    bool   `json:"invert_filter"`
		SubmoduleFolder string `json:"submodule_folder"`
		ShallowClone    bool   `json:"shallow_clone"`
	} `json:"attributes"`
}

type Modification struct {
	EmailAddress string `json:"email_address,omitempty"`
	ID           int    `json:"id,omitempty"`
	ModifiedTime int    `json:"modified_time,omitempty"`
	UserName     string `json:"user_name,omitempty"`
	Comment      string `json:"comment,omitempty"`
	Revision     string `json:"revision,omitempty"`
}

type MaterialRevision struct {
	Modifications []Modification `json:"modifications,omitempty"`
	Material      Material       `json:"material,omitempty"`
	Changed       bool           `json:"changed,omitempty"`
}

type BuildCause struct {
	Approver          string             `json:"approver, omitempty"`
	MaterialRevisions []MaterialRevision `json:"material_revisions,omitempty"`
	TriggerForced     bool               `json:"trigger_forced,omitempty"`
	TriggerMessage    string             `json:"trigger_message,omitempty"`
}

type Job struct {
	AgentUUID           string                `json:"agent_uuid,omitempty"`
	Name                string                `json:"name,omitempty"`
	JobStateTransitions []JobStateTransitions `json:"job_state_transitions,omitempty"`
	ScheduledDate       int                   `json:"scheduled_date,omitempty"`
	OriginalJobID       string                `json:"original_job_id,omitempty"`
	PipelineCounter     int                   `json:"pipeline_counter,omitempty"`
	Rerun               bool                  `json:"rerun,omitempty"`
	PipelineName        string                `json:"pipeline_name,omitempty"`
	Result              string                `json:"result,omitempty"`
	State               string                `json:"state,omitempty"`
	ID                  int                   `json:"id,omitempty"`
	StageCounter        string                `json:"stage_counter,omitempty"`
	StageName           string                `json:"stage_name,omitempty"`
}

type JobConfig struct {
	Name                 string                   `json:"name"`
	RunInstanceCount     int                      `json:"run_instance_count"`
	Timeout              int                      `json:"timeout"`
	EnvironmentVariables []map[string]interface{} `json:"environment_variables"`
	Resources            []string                 `json:"resources"`
	Tasks                []map[string]interface{} `json:"tasks"`
}

func (p *JobConfig) AddTask(task interface{}) error {
	switch task.(type) {
	case TaskExecConfig:
		p.Tasks = append(p.Tasks, structs.Map(task))
		return nil
	case TaskAntConfig:
		p.Tasks = append(p.Tasks, structs.Map(task))
		return nil
	case TaskNantConfig:
		p.Tasks = append(p.Tasks, structs.Map(task))
		return nil
	default:
		return fmt.Errorf("Type %T not support", task)
	}
}

type TaskExecConfig struct {
	Type       string
	attributes struct {
		RunIf            []string `json:"run_if"`
		Command          string   `json:"command"`
		Arguments        []string `json:"arguments"`
		WorkingDirectory string   `json:"working_directory"`
	} `json:"attributes"`
}

func NewTaskExecConfig() *TaskExecConfig {
	return &TaskExecConfig{Type: "exec"}
}

type TaskAntConfig struct {
	Type       string `json:"type"`
	attributes struct {
		RunIf            []string `json:"run_if"`
		WorkingDirectory string   `json:"working_directory"`
		BuildFile        string   `json:"build_file"`
		Target           string   `json:"build_file"`
	} `json:"attributes"`
}

func NewTaskAntConfig() *TaskAntConfig {
	return &TaskAntConfig{Type: "ant"}
}

type TaskNantConfig struct {
	Type       string `json:"type"`
	attributes struct {
		RunIf            []string `json:"run_if"`
		WorkingDirectory string   `json:"working_directory"`
		BuildFile        string   `json:"build_file"`
		Target           string   `json:"build_file"`
		NantPath         string   `json:"nant_path"`
	} `json:"attributes"`
}

func NewTaskNantConfig() *TaskNantConfig {
	return &TaskNantConfig{Type: "nant"}
}

type Stage struct {
	Name                  string `json:"name,omitempty"`
	CleanWorkingDirectory bool   `json:"clean_working_directory,omitempty"`
	ApprovedBy            string `json:"approved_by,omitempty"`
	Jobs                  []Job  `json:"jobs,omitempty"`
	PipelineCounter       int    `json:"pipeline_counter,omitempty"`
	PipelineName          string `json:"pipeline_name,omitempty"`
	ApprovalType          string `json:"approval_type,omitempty"`
	Result                string `json:"result,omitempty"`
	Counter               int    `json:"counter,omitempty"`
	ID                    int    `json:"id,omitempty"`
	RerunOfCounter        int    `json:"rerun_of_counter,omitempty"`
	FetchMaterials        bool   `json:"fetch_materials,omitempty"`
	ArtifactsDeleted      bool   `json:"artifacts_deleted,omitempty"`
}

func NewStage() *Stage {
	return &Stage{
		CleanWorkingDirectory: false,
		Jobs:             make([]Job, 0),
		FetchMaterials:   false,
		ArtifactsDeleted: false}
}

type StageConfig struct {
	Name                  string `json:"name"`
	FetchMaterials        bool   `json:"fetch_materials"`
	CleanWorkingDirectory bool   `json:"clean_working_directory"`
	NeverCleanupArtifacts bool   `json:"never_cleanup_artifacts"`
	Approval              struct {
		Type          string `json:"type"`
		Authorization struct {
			Roles []string `json:"roles"`
			Users []string `json:"users"`
		} `json:"authorization"`
	} `json:"approval"`
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables"`
	Jobs                 []JobConfig           `json:"jobs"`
}

type PipelineInstance struct {
	Name         string     `json:"name,omitempty"`
	NaturalOrder float64    `json:"natural_order,omitempty"`
	CanRun       bool       `json:"can_run,omitempty"`
	Comment      string     `json:"comment,omitempty"`
	Stages       []Stage    `json:"stages,omitempty"`
	Counter      int        `json:"counter,omitempty"`
	ID           int        `json:"id,omitempty"`
	Label        string     `json:"label,omitempty"`
	BuildCause   BuildCause `json:"build_cause,omitempty"`
}

func NewPipelineInstance() *PipelineInstance {
	return &PipelineInstance{Stages: make([]Stage, 0)}
}

type PipelineConfig struct {
	LabelTemplate         string                   `json:"label_template,omitempty"`
	EnablePipelineLocking bool                     `json:"enable_pipeline_locking"`
	Name                  string                   `json:"name"`
	Template              string                   `json:"template"`
	Params                []map[string]string      `json:"parameters"`
	EnvironmentVariables  []map[string]interface{} `json:"environment_variables"`
	Materials             []MaterialGitConfig      `json:"materials"`
	Stages                []StageConfig            `json:"stages"`
}

func NewPipelineConfig() *PipelineConfig {
	return &PipelineConfig{
		EnablePipelineLocking: false,
		Params:                make([]map[string]string, 0),
		EnvironmentVariables:  make([]map[string]interface{}, 0),
		Materials:             make([]MaterialGitConfig, 0),
		Stages:                make([]StageConfig, 0)}
}
