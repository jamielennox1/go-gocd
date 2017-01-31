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
	ETag     []string
}

func New(host, login, password string) *Client {
	return &Client{host: host, login: login, password: password}
}

func (p *Client) unmarshal(data io.ReadCloser, v interface{}) error {
	defer data.Close()

	body, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
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

func (p *Client) GetPipelineInstance(name string, inst int) (*PipelineInstance, error) {
	resp, err := p.goCDRequest("GET", fmt.Sprintf("%s/go/api/pipelines/%s/instance/%d", p.host, name, inst), []byte{}, map[string]string{})
	if err != nil {
		return nil, err
	}

	pipeline := PipelineInstance{}

	if err := p.unmarshal(resp.Body, &pipeline); err != nil {
		return nil, err
	} else {
		return &pipeline, nil
	}
}

func (p *Client) GetHistoryPipelineInstance(name string) (*PipelineInstances, error) {
	resp, err := p.goCDRequest("GET", fmt.Sprintf("%s/go/api/pipelines/%s/history", p.host, name), []byte{}, map[string]string{})
	if err != nil {
		return nil, err
	}
	pipelines := PipelineInstances{}
	if err := p.unmarshal(resp.Body, &pipelines); err != nil {
		return nil, err
	} else {
		return &pipelines, nil
	}
}

func (p *Client) GetPipelineConfig(name string) (*PipelineConfig, error) {
	resp, err := p.goCDRequest("GET", fmt.Sprintf("%s/go/api/admin/pipelines/%s", p.host, name), []byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v2+json"})
	if err != nil {
		return nil, err
	}

	pipeline := PipelineConfig{}

	if err := p.unmarshal(resp.Body, &pipeline); err != nil {
		return nil, err
	} else {
		p.ETag = resp.Header["ETag"]
		return &pipeline, nil
	}
}

func (p *Client) NewPipelineConfig(pipeline *PipelineConfig, group string) error {
	data := struct {
		Group    string         `json:"group"`
		Pipeline PipelineConfig `json:"pipeline"`
	}{Group: group, Pipeline: pipeline}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(p.ETag) == 0 {
		return fmt.Errorf("ETag is empty")
	}

	if resp, err := p.goCDRequest("POST", fmt.Sprintf("%s/go/api/admin/pipelines", p.host), body,
		map[string]string{
			"Accept": "application/vnd.go.cd.v2+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Operation error: %s", resp.Status)
	}
	return nil

}

func (p *Client) SetPipelineConfig(pipeline *PipelineConfig) error {
	body, err := json.Marshal(pipeline)
	if err != nil {
		return err
	}

	if len(p.ETag) == 0 {
		return fmt.Errorf("ETag is empty")
	}

	if resp, err := p.goCDRequest("PUT", fmt.Sprintf("%s/go/api/admin/pipelines/%s", p.host, pipeline.Name), body,
		map[string]string{
			"Accept":   "application/vnd.go.cd.v2+json",
			"If-Match": p.ETag[0]}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Operation error: %s", resp.Status)
	}
	return nil

}

func (p *Client) GetEnvironments() (*Environments, error) {
	resp, err := p.goCDRequest("GET", fmt.Sprintf("%s/go/api/admin/environments", p.host), []byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Operation error: %s", resp.Status)
	}

	envs := Environments{}
	if err := p.unmarshal(resp.Body, &envs); err != nil {
		return nil, err
	} else {
		p.ETag = resp.Header["ETag"]
		return &envs, nil
	}
}

func (p *Client) GetEnvironment(name string) (*Environment, error) {
	resp, err := p.goCDRequest("GET", fmt.Sprintf("%s/go/api/admin/environments/%s", p.host, name), []byte{},
		map[string]string{"Accept": "application/vnd.go.cd.v1+json"})
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Operation error: %s", resp.Status)
	}

	env := Environment{}
	if err := p.unmarshal(resp.Body, &env); err != nil {
		return nil, err
	} else {
		p.ETag = resp.Header["ETag"]
		return &env, nil
	}
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

	if len(p.ETag) == 0 {
		return fmt.Errorf("ETag is empty")
	}

	if resp, err := p.goCDRequest("PUT", fmt.Sprintf("%s/go/api/admin/environments/%s", p.host, env.Name), body,
		map[string]string{
			"If-Match": p.ETag[0],
			"Accept":   "application/vnd.go.cd.v1+json"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Operation error: %s", resp.Status)
	}
	return nil
}

func (p *Client) UnpausePipeline(resource string) error {
	if resp, err := p.goCDRequest("POST", fmt.Sprintf("%s/unpause", p.host), []byte{}, map[string]string{"Confirm": "true"}); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Operation error: %s", resp.Status)
	}
	return nil
}
