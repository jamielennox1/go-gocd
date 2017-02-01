package gocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	host     string
	login    string
	password string
	Etag     string
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
	req.Header.Set("Content-Type", "application/json")
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

	version := Version{}
	if err := p.unmarshal(resp.Body, &version); err != nil {
		return nil, err
	} else {
		return &version, nil
	}
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

	pipeline := PipelineInstance{}
	if err := p.unmarshal(resp.Body, &pipeline); err != nil {
		return nil, err
	} else {
		return &pipeline, nil
	}
}

func (p *Client) GetHistoryPipelineInstance(name string) (*PipelineInstances, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/pipelines/%s/history", p.host, name),
		[]byte{},
		map[string]string{})
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, p.createError(resp)
	}

	pipelines := PipelineInstances{}
	if err := p.unmarshal(resp.Body, &pipelines); err != nil {
		return nil, err
	} else {
		return &pipelines, nil
	}
}

func (p *Client) GetPipelineConfig(name string) (*PipelineConfig, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/admin/pipelines/%s", p.host, name),
		[]byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v2+json"})
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, p.createError(resp)
	}

	pipeline := PipelineConfig{}
	if err := p.unmarshal(resp.Body, &pipeline); err != nil {
		return nil, err
	} else {
		if tag := resp.Header["Etag"]; len(tag) > 0 {
			p.Etag = tag[0]
		}
		return &pipeline, nil
	}
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
	if resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/admin/pipelines", p.host),
		body,
		map[string]string{"Accept": "application/vnd.go.cd.v2+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return p.createError(resp)
	}
	return nil
}

func (p *Client) NewPipelineConfigRaw(data []byte) error {
	if resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/admin/pipelines", p.host),
		data,
		map[string]string{"Accept": "application/vnd.go.cd.v2+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return p.createError(resp)
	}
	return nil
}

func (p *Client) SetPipelineConfig(pipeline *PipelineConfig) error {
	body, err := json.Marshal(pipeline)
	if err != nil {
		return err
	}
	if resp, err := p.goCDRequest("PUT",
		fmt.Sprintf("%s/go/api/admin/pipelines/%s", p.host, pipeline.Name),
		body,
		map[string]string{"If-Match": p.Etag,
			"Accept": "application/vnd.go.cd.v2+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return p.createError(resp)
	}
	return nil
}

func (p *Client) SetPipelineConfigRaw(name string, data []byte) error {
	if resp, err := p.goCDRequest("PUT",
		fmt.Sprintf("%s/go/api/admin/pipelines/%s", p.host, name),
		data,
		map[string]string{"If-Match": p.Etag,
			"Accept": "application/vnd.go.cd.v2+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return p.createError(resp)
	}
	return nil
}

func (p *Client) DeletePipelineConfig(name string) error {
	if resp, err := p.goCDRequest("DELETE",
		fmt.Sprintf("%s/go/api/admin/pipelines/%s", p.host, name),
		[]byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v2+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return p.createError(resp)
	}

	envs, err := p.GetEnvironments()
	if err != nil {
		return err
	}

	for _, env := range envs.Embeded.Environments {
		if env.DeletePipeline(name) {
			if err := p.SetEnvironment(&env); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func (p *Client) GetEnvironments() (*Environments, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/admin/environments", p.host),
		[]byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, p.createError(resp)
	}

	envs := Environments{}
	if err := p.unmarshal(resp.Body, &envs); err != nil {
		return nil, err
	} else {
		if tag := resp.Header["Etag"]; len(tag) > 0 {
			p.Etag = tag[0]
		}
		return &envs, nil
	}
}

func (p *Client) GetEnvironment(name string) (*Environment, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/admin/environments/%s", p.host, name),
		[]byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, p.createError(resp)
	}

	env := Environment{}
	if err := p.unmarshal(resp.Body, &env); err != nil {
		return nil, err
	} else {
		if tag := resp.Header["Etag"]; len(tag) > 0 {
			p.Etag = tag[0]
		}
		return &env, nil
	}
}

func (p *Client) NewEnvironment(env *Environment) error {
	data := struct {
		Name                 string                `json:"name"`
		Pipelines            []map[string]string   `json:"pipelines"`
		Agents               []map[string]string   `json:"agents"`
		EnvironmentVariables []EnvironmentVariable `json:"environment_variables"`
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

	if resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/admin/environments", p.host),
		body,
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return p.createError(resp)
	}
	return nil
}

func (p *Client) SetEnvironment(env *Environment) error {
	data := struct {
		Name                 string                `json:"name"`
		Pipelines            []map[string]string   `json:","`
		Agents               []map[string]string   `json:","`
		EnvironmentVariables []EnvironmentVariable `json:"environment_variables"`
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

	if resp, err := p.goCDRequest("PUT",
		fmt.Sprintf("%s/go/api/admin/environments/%s", p.host, env.Name),
		body,
		map[string]string{
			"If-Match": p.Etag,
			"Accept":   "application/vnd.go.cd.v1+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return p.createError(resp)
	}
	return nil
}

func (p *Client) DeleteEnvironment(name string) error {
	if resp, err := p.goCDRequest("DELETE",
		fmt.Sprintf("%s/go/api/admin/environments/%s", p.host, name),
		[]byte{},
		map[string]string{"If-Match": p.Etag,
			"Accept": "application/vnd.go.cd.v1+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return p.createError(resp)
	}
	return nil
}

func (p *Client) UnpausePipeline(name string) error {
	if resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/pipelines/%s/unpause", p.host, name),
		[]byte{},
		map[string]string{"Confirm": "true"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return p.createError(resp)
	}
	return nil
}

func (p *Client) PausePipeline(name string) error {
	if resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/pipelines/%s/pause", p.host, name),
		[]byte{'p', 'a', 'u', 's', 'e', 'C', 'a', 'u', 's', 'e', '=', 't', 'a', 'k', 'e', ' ', 's', 'o', 'm', 'e', ' ', 'r', 'e', 's', 't'},
		map[string]string{"Confirm": "true"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return p.createError(resp)
	}
	return nil
}

func (p *Client) SchedulePipeline(name string, data []byte) error {
	if resp, err := p.goCDRequest("POST",
		fmt.Sprintf("%s/go/api/pipelines/%s/schedule", p.host, name),
		data,
		map[string]string{"Confirm": "true"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusAccepted {
		return p.createError(resp)
	}
	return nil
}

func (p *Client) GetGroups() (*[]Group, error) {
	resp, err := p.goCDRequest("GET",
		fmt.Sprintf("%s/go/api/config/pipeline_groups", p.host),
		[]byte{},
		map[string]string{})
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, p.createError(resp)
	}

	groups := []Group{}
	if err := p.unmarshal(resp.Body, &groups); err != nil {
		return nil, err
	} else {
		return &groups, nil
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
