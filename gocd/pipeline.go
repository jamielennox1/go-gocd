package gocd

type JobStateTransitions struct {
	//StateChangeTime time.Time `json:"state_change_time,omitempty"`
	ID    int    `json:"id,omitempty"`
	State string `json:"state,omitempty"`
}

type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Material struct {
	Description string `json:"description,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Type        string `json:"type,omitempty"`
	ID          int    `json:"id,omitempty"`
}

type Modification struct {
	EmailAddress string `json:"email_address,omitempty"`
	ID           int    `json:"id,omitempty"`
	//ModifiedTime time.Time `json:"modified_time,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Revision string `json:"revision,omitempty"`
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
	//ScheduledDate       time.Time             `json:"scheduled_date,omitempty"`
	OriginalJobID   string `json:"original_job_id,omitempty"`
	PipelineCounter int    `json:"pipeline_counter,omitempty"`
	Rerun           bool   `json:"rerun,omitempty"`
	PipelineName    string `json:"pipeline_name,omitempty"`
	Result          string `json:"result,omitempty"`
	State           string `json:"state,omitempty"`
	ID              int    `json:"id,omitempty"`
	StageCounter    string `json:"stage_counter,omitempty"`
	StageName       string `json:"stage_name,omitempty"`
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
	Counter               string `json:"counter,omitempty"`
	ID                    int    `json:"id,omitempty"`
	RerunOfCounter        int    `json:"rerun_of_counter,omitempty"`
	FetchMaterials        bool   `json:"fetch_materials,omitempty"`
	ArtifactsDeleted      bool   `json:"artifacts_deleted,omitempty"`
}

type PipelineInstance struct {
	Name         string     `json:"name,omitempty"`
	NaturalOrder float64    `json:"natural_order,omitempty"`
	CanRun       bool       `json:"can_run,omitempty"`
	Comment      string     `json:"comment,omitempty"`
	Stages       []Stage    `json:"stages,omitempty"`
	Counter      int        `json:"counter,omitempty"`
	ID           int        `json:"id,omitempty"`
	label        string     `json:"label,omitempty"`
	BuildCause   BuildCause `json:"build_cause,omitempty"`
}

type PipelineInstances struct {
	Instances []PipelineInstance `json:"pipelines"`
}

type PipelineConfig struct {
	LabelTemplate         string                `json:"label_template,omitempty"`
	EnablePipelineLocking bool                  `json:"enable_pipeline_locking"`
	Name                  string                `json:"name"`
	Template              string                `json:"template"`
	Parameters            []Parameter           `json:"parameters"`
	EnvironmentVariables  []EnvironmentVariable `json:"environment_variables"`
	Materials             []Material            `json:"materials"`
	Stages                []Stage               `json:"stages"`
}
