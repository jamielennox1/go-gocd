package gocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

const VERSION = "0.1.0"

type Client struct {
	host     string
	login    string
	password string
	Etag     string
	EtagEnv  string
}

func New(host, login, password string) *Client {
	return &Client{host: host, login: login, password: password}
}

func (p *Client) unmarshal(data io.ReadCloser, v interface{}) error {
	defer data.Close()
	if body, err := ioutil.ReadAll(data); err != nil {
		return err
	} else {

		//fmt.Println(string(body))

		return json.Unmarshal(body, v)
	}
}

func (p *Client) createError(resp *http.Response) error {
	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		return fmt.Errorf("Operation error: %s (%s)", resp.Status, body)
	}
	return fmt.Errorf("Operation error: %s", resp.Status)
}

func (p *Client) goCDRequest(method string, resource string, body []byte, headers map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest(method, resource, bytes.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.SetBasicAuth(p.login, p.password)
	return http.DefaultClient.Do(req)
}

func (p *Client) Version() (*Version, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/version", p.host),
		[]byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, p.createError(resp)
	}

	version := Version{ClientVersion: VERSION}
	return &version, p.unmarshal(resp.Body, &version)
}

func (p *Client) GetPipelineInstance(name string, inst int) (*PipelineInstance, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/pipelines/%s/instance/%d", p.host, name, inst),
		[]byte{},
		map[string]string{})
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, p.createError(resp)
	}

	pipeline := NewPipelineInstance()
	return pipeline, p.unmarshal(resp.Body, pipeline)
}

func (p *Client) GetHistoryPipelineInstance(name string) ([]*PipelineInstance, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/pipelines/%s/history", p.host, name),
		[]byte{},
		map[string]string{})

	switch true {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, p.createError(resp)
	}

	pipelines := struct {
		Instances []*PipelineInstance `json:"pipelines"`
	}{make([]*PipelineInstance, 0)}

	return pipelines.Instances, p.unmarshal(resp.Body, pipelines)
}

func (p *Client) GetPipelineConfig(name string) (*PipelineConfig, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/admin/pipelines/%s", p.host, name),
		[]byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v2+json"})

	switch true {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, p.createError(resp)
	default:
		if tag := resp.Header["Etag"]; len(tag) > 0 {
			p.Etag = tag[0]
		}
	}

	pipeline := NewPipelineConfig()
	return pipeline, p.unmarshal(resp.Body, pipeline)
}

func (p *Client) NewPipelineConfig(pipeline *PipelineConfig, group string) error {
	data := struct {
		Group    string         `json:"group"`
		Pipeline PipelineConfig `json:"pipeline"`
	}{Group: group, Pipeline: *pipeline}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/admin/pipelines", p.host),
		body,
		map[string]string{"Content-Type": "application/json",
			"Accept": "application/vnd.go.cd.v2+json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) NewPipelineConfigRaw(data []byte) error {
	resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/admin/pipelines", p.host),
		data,
		map[string]string{"Content-Type": "application/json",
			"Accept": "application/vnd.go.cd.v2+json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) SetPipelineConfig(pipeline *PipelineConfig) error {
	body, err := json.Marshal(pipeline)
	if err != nil {
		return err
	}
	resp, err := p.goCDRequest("PUT",
		fmt.Sprintf("%s/go/api/admin/pipelines/%s", p.host, pipeline.Name),
		body,
		map[string]string{"If-Match": p.Etag,
			"Content-Type": "application/json",
			"Accept":       "application/vnd.go.cd.v2+json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		if tag := resp.Header["Etag"]; len(tag) > 0 {
			p.Etag = tag[0]
		}
		return nil
	}
}

func (p *Client) SetPipelineConfigRaw(name string, data []byte) error {
	resp, err := p.goCDRequest("PUT",
		fmt.Sprintf("%s/go/api/admin/pipelines/%s", p.host, name),
		data,
		map[string]string{"If-Match": p.Etag,
			"Content-Type": "application/json",
			"Accept":       "application/vnd.go.cd.v2+json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		if tag := resp.Header["Etag"]; len(tag) > 0 {
			p.Etag = tag[0]
		}
		return nil
	}
}

func (p *Client) DeletePipelineConfig(name string) error {
	pipeline, env, err := p.FindPipelineConfig(name)
	if pipeline == nil {
		return fmt.Errorf("%s not found", name)
	}
	if env != nil {
		env.DeletePipeline(name)
		if err := p.SetEnvironment(env); err != nil {
			return err
		}
	}

	resp, err := p.goCDRequest("DELETE",
		fmt.Sprintf("%s/go/api/admin/pipelines/%s", p.host, name),
		[]byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v2+json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) GetEnvironments() (*Environments, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/admin/environments", p.host),
		[]byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})

	switch true {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, p.createError(resp)
	default:
		if tag := resp.Header["Etag"]; len(tag) > 0 {
			p.EtagEnv = tag[0]
		}
	}

	envs := NewEnvironments()
	return envs, p.unmarshal(resp.Body, envs)
}

func (p *Client) GetEnvironment(name string) (*Environment, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/admin/environments/%s", p.host, name),
		[]byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})

	switch true {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, p.createError(resp)
	default:
		if tag := resp.Header["Etag"]; len(tag) > 0 {
			p.EtagEnv = tag[0]
		}
	}

	env := NewEnvironment()
	return env, p.unmarshal(resp.Body, env)
}

func (p *Client) NewEnvironment(env *Environment) error {
	data := struct {
		Name                 string                   `json:"name"`
		Pipelines            []map[string]string      `json:"pipelines"`
		Agents               []map[string]string      `json:"agents"`
		EnvironmentVariables []map[string]interface{} `json:"environment_variables"`
	}{Name: env.Name, EnvironmentVariables: env.EnvironmentVariables}

	for _, p := range env.Pipelines {
		data.Pipelines = append(data.Pipelines, map[string]string{"name": p.Name})
	}
	for _, a := range env.Agents {
		data.Agents = append(data.Agents, map[string]string{"uuid": a.Uuid})
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/admin/environments", p.host),
		body,
		map[string]string{"Content-Type": "application/json",
			"Accept": "application/vnd.go.cd.v1+json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		if tag := resp.Header["Etag"]; len(tag) > 0 {
			p.EtagEnv = tag[0]
		}
		return nil
	}
}

func (p *Client) SetEnvironment(env *Environment) error {
	data := struct {
		Name                 string                   `json:"name"`
		Pipelines            []map[string]string      `json:"pipelines"`
		Agents               []map[string]string      `json:"agents"`
		EnvironmentVariables []map[string]interface{} `json:"environment_variables"`
	}{Name: env.Name}

	for _, p := range env.Pipelines {
		data.Pipelines = append(data.Pipelines, map[string]string{"name": p.Name})
	}
	for _, a := range env.Agents {
		data.Agents = append(data.Agents, map[string]string{"uuid": a.Uuid})
	}
	data.EnvironmentVariables = env.EnvironmentVariables

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	p.GetEnvironment(env.Name)

	resp, err := p.goCDRequest("PUT",
		fmt.Sprintf("%s/go/api/admin/environments/%s", p.host, env.Name),
		body,
		map[string]string{"If-Match": p.EtagEnv,
			"Content-Type": "application/json",
			"Accept":       "application/vnd.go.cd.v1+json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		if tag := resp.Header["Etag"]; len(tag) > 0 {
			p.EtagEnv = tag[0]
		}
		return nil
	}
}

func (p *Client) DeleteEnvironment(name string) error {
	resp, err := p.goCDRequest("DELETE",
		fmt.Sprintf("%s/go/api/admin/environments/%s", p.host, name),
		[]byte{},
		map[string]string{"If-Match": p.EtagEnv,
			"Accept": "application/vnd.go.cd.v1+json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) UnpausePipeline(name string) error {
	resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/pipelines/%s/unpause", p.host, name),
		[]byte{},
		map[string]string{"Confirm": "true"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) PausePipeline(name string) error {
	resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/pipelines/%s/pause", p.host, name),
		[]byte{'p', 'a', 'u', 's', 'e', 'C', 'a', 'u', 's', 'e', '=', 't', 'a', 'k', 'e', ' ', 's', 'o', 'm', 'e', ' ', 'r', 'e', 's', 't'},
		map[string]string{"Confirm": "true"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) SchedulePipeline(name string, data []byte) error {
	resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/pipelines/%s/schedule", p.host, name),
		data,
		map[string]string{"Confirm": "true"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusAccepted:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) GetGroups() (*[]*Group, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/config/pipeline_groups", p.host),
		[]byte{},
		map[string]string{})

	switch true {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, p.createError(resp)
	}

	groups := make([]*Group, 0)
	if err := p.unmarshal(resp.Body, &groups); err != nil {
		return nil, err
	} else {
		return &groups, nil
	}
}

func (p *Client) StageCancel(pipeline string, stage string) error {
	resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/stages/%s/%s/cancel", p.host, pipeline, stage),
		make([]byte, 0),
		map[string]string{"Confirm": "true"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) GetStageInstance(pipeline string, pInst int, stage string, sInst int) (*Stage, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/stages/%s/%s/instance/%d/%d", p.host, pipeline, stage, pInst, sInst),
		make([]byte, 0),
		map[string]string{})

	switch true {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, p.createError(resp)
	}

	s := NewStage()
	return s, p.unmarshal(resp.Body, s)
}

func (p *Client) GetStageInstanceHystory(pipeline string, stage string) ([]*Stage, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/stages/%s/%s/history", p.host, pipeline, stage),
		make([]byte, 0),
		map[string]string{})

	switch true {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, p.createError(resp)
	}

	stages := struct {
		Stages []*Stage `json:"stages"`
	}{Stages: make([]*Stage, 0)}
	return stages.Stages, p.unmarshal(resp.Body, &stages)
}

func (p *Client) GetAllAgents() ([]*Agent, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/agents", p.host),
		make([]byte, 0),
		map[string]string{"Accept": "application/vnd.go.cd.v2+json"})

	switch true {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, p.createError(resp)
	}

	data := struct {
		Embeded struct {
			Agents []*Agent `json:"agents"`
		} `json:"_embedded"`
	}{Embeded: struct {
		Agents []*Agent `json:"agents"`
	}{Agents: make([]*Agent, 0)}}

	return data.Embeded.Agents, p.unmarshal(resp.Body, &data)
}

func (p *Client) GetAgent(uuid string) (*Agent, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/agents/%s", p.host, uuid),
		make([]byte, 0),
		map[string]string{"Accept": "application/vnd.go.cd.v2+json"})

	switch true {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, p.createError(resp)
	}

	agent := NewAgent()
	return agent, p.unmarshal(resp.Body, agent)
}

func (p *Client) SetAgent(agent Agent) error {
	old_agent, err := p.GetAgent(agent.Uuid)
	if err != nil {
		return err
	}
	diff := old_agent.Diff(agent)
	if reflect.DeepEqual(diff, make(map[string]interface{})) {
		return nil
	}
	body, err := json.Marshal(diff)
	if err != nil {
		return err
	}

	resp, err := p.goCDRequest("PATCH",
		fmt.Sprintf("%s/go/api/agents/%s", p.host, agent.Uuid),
		body,
		map[string]string{"Accept": "application/vnd.go.cd.v2+json",
			"Content-Type": "application/json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) DeleteAgent(uuid string) error {
	resp, err := p.goCDRequest("DELETE",
		fmt.Sprintf("%s/go/api/agents/%s", p.host, uuid),
		make([]byte, 0),
		map[string]string{"Accept": "application/vnd.go.cd.v2+json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) GetAllUsers() ([]*User, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/users", p.host),
		make([]byte, 0),
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})

	switch true {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, p.createError(resp)
	}

	data := struct {
		Embeded struct {
			Users []*User `json:"users"`
		} `json:"_embedded"`
	}{Embeded: struct {
		Users []*User `json:"users"`
	}{Users: make([]*User, 0)}}

	return data.Embeded.Users, p.unmarshal(resp.Body, &data)
}

func (p *Client) GetUser(login string) (*User, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/users/%s", p.host, login),
		make([]byte, 0),
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})

	switch true {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, p.createError(resp)
	}

	user := NewUser()
	return user, p.unmarshal(resp.Body, user)
}

func (p *Client) NewUser(user *User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/users", p.host),
		body,
		map[string]string{"Accept": "application/vnd.go.cd.v1+json",
			"Content-Type": "application/json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) SetUser(user *User) error {
	old_user, err := p.GetUser(user.LoginName)
	if err != nil {
		return err
	}

	diff := old_user.Diff(user)

	if reflect.DeepEqual(diff, make(map[string]interface{})) {
		return nil
	}

	body, err := json.Marshal(diff)
	if err != nil {
		return err
	}

	resp, err := p.goCDRequest("PATCH",
		fmt.Sprintf("%s/go/api/users/%s", p.host, user.LoginName),
		body,
		map[string]string{"Accept": "application/vnd.go.cd.v1+json",
			"Content-Type": "application/json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) DeleteUser(login string) error {
	resp, err := p.goCDRequest("DELETE",
		fmt.Sprintf("%s/go/api/users/%s", p.host, login),
		make([]byte, 0),
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})

	switch true {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return p.createError(resp)
	default:
		return nil
	}
}

func (p *Client) FindPipelineConfig(name string) (*PipelineConfig, *Environment, error) {
	pipeline, err := p.GetPipelineConfig(name)
	if err != nil {
		return nil, nil, err
	}
	envs, err := p.GetEnvironments()
	if err != nil {
		return pipeline, nil, err
	}
	for _, env := range envs.Embeded.Environments {
		if env.ExistPipeline(name) {
			return pipeline, &env, nil
		}
	}
	return pipeline, nil, nil
}
