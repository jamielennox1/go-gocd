package gocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Link struct {
	HRef string `json:"href"`
}

type Links struct {
	Self Link `json:"self"`
	Doc  Link `json:"doc"`
	Find Link `json:"find"`
}

type EnvironmentVariable struct {
	Secure          bool   `json:"secure"`
	Name            string `json:"name"`
	Value           string `json:"value, omitempty"`
	Encrypted_value string `json:"encrypted_value, omitempty"`
}

type Environment struct {
	Links                Links                 `json:"_links"`
	Name                 string                `json:"name"`
	Agents               []Agent               `json:"agents"`
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables"`
	Pipelines            []ShortPipeline       `json:"pipelines"`
}

func (p *Environment) AddPipeline(name string) {
	for _, p := range p.Pipelines {
		if strings.Compare(p.Name, name) == 0 {
			return
		}
	}
	p.Pipelines = append(p.Pipelines, ShortPipeline{Name: name})
}

func (p *Environment) DeletePipeline(name string) {
	for i := 0; i < len(p.Pipelines); i++ {
		if strings.Compare(p.Pipelines[i].Name, name) == 0 {
			p.Pipelines[i] = p.Pipelines[len(p.Pipelines)-1]
			p.Pipelines = p.Pipelines[:len(p.Pipelines)-1]
		}
	}
}

type Environments struct {
	Links        Links         `json:"_links"`
	Environments []Environment `json:"_embedded,struct environments"`
}

type Agent struct {
	Links Links  `json:"_links"`
	Uuid  string `json:"uuid"`
}

type ShortPipeline struct {
	Links Links  `json:"_links"`
	Name  string `json:"name"`
}

type JobStateTransitions struct {
	//StateChangeTime time.Time `json:"state_change_time,omitempty"`
	ID    int    `json:"id,omitempty"`
	State string `json:"state,omitempty"`
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

type Pipeline struct {
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

type Client struct {
	host     string
	login    string
	password string
}

func New(host, login, password string) *Client {
	return &Client{host: host, login: login, password: password}
}

func (p *Client) GetPipeline(name string, inst int) (*Pipeline, error) {
	resp, err := p.goCDRequest("GET", fmt.Sprintf("%s/go/api/pipelines/%s/instance/%d", p.host, name, inst), "", map[string]string{})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	pipeline := Pipeline{}
	if err := json.Unmarshal(body, &pipeline); err != nil {
		return nil, err
	} else {
		return &pipeline, nil
	}
}

func (p *Client) goCDRequest(method string, resource string, body string, headers map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest(method, resource, bytes.NewReader([]byte(body)))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(p.login, p.password)
	return http.DefaultClient.Do(req)
}

func (p *Client) CreatePipeline(name string, env string, resource string, body string, headers map[string]string) error {
	if resp, err := p.goCDRequest("POST", fmt.Sprintf("%s/go/api/admin/pipelines", resource), body, headers); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			fmt.Errorf("Operation body: %s", body)
		}
		return fmt.Errorf("Operation error: %s", resp.Status)
	}
	data, tag, err := p.UpdateEnv(resource, env, name, "")
	if err != nil {
		return err
	}

	if resp, err := p.goCDRequest("PUT", fmt.Sprintf("%s/go/api/admin/environments/%s", resource, env), data,
		map[string]string{"If-Match": tag, "Accept": "application/vnd.go.cd.v1+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Operation error: %s", resp.Status)
	}

	return p.UnpausePipeline(fmt.Sprintf("%s/go/api/pipelines/%s", resource, name))
}

func (p *Client) UnpausePipeline(resource string) error {
	if resp, err := p.goCDRequest("POST", resource+"/unpause", "", map[string]string{"Confirm": "true"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Operation error: %s", resp.Status)
	}
	return nil
}

func (p *Client) UpdatePipeline(name string, environment string, resource string, body string, headers map[string]string) error {
	if resp, err := p.goCDRequest("PUT", fmt.Sprintf("%s/go/api/admin/pipelines/%s", resource, name), body, headers); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			fmt.Errorf("Operation body: %s", body)
		}
		return fmt.Errorf("Operation error: %s", resp.Status)
	}

	if cEnv, err := p.FindEnv(resource, name); err == nil {
		if environment != cEnv && cEnv != "" {

			data, tag, err := p.UpdateEnv(resource, cEnv, "", name)
			if err != nil {
				return err
			}
			if resp, err := p.goCDRequest("PUT", fmt.Sprintf("%s/go/api/admin/environments/%s", resource, cEnv), data,
				map[string]string{"If-Match": tag, "Accept": "application/vnd.go.cd.v1+json"}); err != nil {
				return err
			} else if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("Operation error: %s", resp.Status)
			}
		}

		data, tag, err := p.UpdateEnv(resource, environment, name, "")
		if err != nil {
			return err
		}

		if resp, err := p.goCDRequest("PUT", fmt.Sprintf("%s/go/api/admin/environments/%s", resource, environment), data,
			map[string]string{"If-Match": tag, "Accept": "application/vnd.go.cd.v1+json"}); err != nil {
			return err
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Operation error: %s", resp.Status)
		}
	} else {
		return err
	}

	return p.UnpausePipeline(fmt.Sprintf("%s/go/api/pipelines/%s", resource, name))
}

func (p *Client) DeletePipeline(name string, env string, resource string, headers map[string]string) error {
	data, tag, err := p.UpdateEnv(resource, env, "", name)
	if err != nil {
		return err
	}

	if resp, err := p.goCDRequest("PUT", fmt.Sprintf("%s/go/api/admin/environments/%s", resource, env), data,
		map[string]string{"If-Match": tag, "Accept": "application/vnd.go.cd.v1+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Operation error: %s", resp.Status)
	}

	if resp, err := p.goCDRequest("DELETE", fmt.Sprintf("%s/go/api/admin/pipelines/%s", resource, name), "", headers); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Operation error: %s", resp.Status)
	}

	return nil
}

func (p *Client) UpdateEnv(resource string, environment string, addPipeline string, delPipeline string) (string, string, error) {
	resp, err := p.goCDRequest("GET", fmt.Sprintf("%s/go/api/admin/environments/%s", resource, environment), "",
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})
	if err != nil {
		return "", "", err
	} else if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("Operation error: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("read body error: %s", body)
	}

	env := Environment{}
	if err := json.Unmarshal(body, &env); err != nil {
		return "", "", err
	}

	env.AddPipeline(p)
	env.DeletePipeline(p)

	if data, err := json.Marshal(env); err != nil {
		return "", "", err
	} else {
		return data, resp.Header.Get("ETag"), nil
	}
}

func (p *Client) FindEnv(resource string, pipeline string) (string, error) {
	resp, err := p.goCDRequest("GET", fmt.Sprintf("%s/go/api/admin/environments", resource), "",
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})
	if err != nil {
		return "", err
	} else if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Operation error: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	envs := Environments{}

	if err := json.Unmarshal(body, &envs); err != nil {
		return "", err
	}

	for _, env := range envs.Environments {
		for _, p := range env.Pipelines {
			if strings.Compare(pipeline, p.Name) == 0 {
				return env.Name, nil
			}
		}

	}
	return "", fmt.Errorf("not found pipeline %s", pipeline)
}
