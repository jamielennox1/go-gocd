package gocd

import (
	"fmt"
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

type Agent struct {
	Links Links  `json:"_links"`
	Uuid  string `json:"uuid"`
}

type ShortPipeline struct {
	Links Links  `json:"_links"`
	Name  string `json:"name"`
}

type EnvironmentVariable struct {
	Secure bool
	Crypt  bool
	Name   string
	Value  string
}

type Environment struct {
	Links                Links                    `json:"_links"`
	Name                 string                   `json:"name"`
	Agents               []Agent                  `json:"agents"`
	EnvironmentVariables []map[string]interface{} `json:"environment_variables"`
	Pipelines            []ShortPipeline          `json:"pipelines"`
}

func NewEnvironment() *Environment {
	return &Environment{Agents: make([]Agent, 0),
		Pipelines:            make([]ShortPipeline, 0),
		EnvironmentVariables: make([]map[string]interface{}, 0)}
}

func (p *Environment) AddPipeline(pipeline string) error {
	for _, p := range p.Pipelines {
		if strings.Compare(p.Name, pipeline) == 0 {
			return fmt.Errorf("Pipeline %s exist", pipeline)
		}
	}
	p.Pipelines = append(p.Pipelines, ShortPipeline{Name: pipeline})
	return nil
}

func (p *Environment) ExistPipeline(pipeline string) bool {
	for _, p := range p.Pipelines {
		if strings.Compare(p.Name, pipeline) == 0 {
			return true
		}
	}
	return false
}

func (p *Environment) DeletePipeline(pipeline string) error {
	for i := 0; i < len(p.Pipelines); i++ {
		if strings.Compare(p.Pipelines[i].Name, pipeline) == 0 {
			p.Pipelines[i] = p.Pipelines[len(p.Pipelines)-1]
			p.Pipelines = p.Pipelines[:len(p.Pipelines)-1]
			return nil
		}
	}
	return fmt.Errorf("Pipeline %s not exist", pipeline)
}

func (p *Environment) AddEnvironmentVariables(env *EnvironmentVariable) error {
	valueField := "value"
	if env.Secure && !env.Crypt {
		valueField = "encrypted_value"
	}
	for _, v := range p.EnvironmentVariables {
		if name, _ := v["name"]; strings.Compare(name.(string), env.Name) == 0 {
			return fmt.Errorf("Env %s exist", env.Name)
		}
	}

	p.EnvironmentVariables = append(p.EnvironmentVariables, map[string]interface{}{
		"name":     env.Name,
		"secure":   env.Secure,
		valueField: env.Value})
	return nil
}

func (p *Environment) GetEnvironmentVariables(name string) (*EnvironmentVariable, error) {
	for _, v := range p.EnvironmentVariables {
		if n, _ := v["name"]; strings.Compare(n.(string), name) == 0 {
			env := EnvironmentVariable{Name: name, Secure: v["secure"].(bool)}
			if val, ok := v["encrypted_value"]; ok {
				env.Value = val.(string)
				env.Crypt = false
			} else {
				env.Value = v["value"].(string)
				env.Crypt = env.Secure
			}
			return &env, nil
		}
	}
	return nil, fmt.Errorf("Env %s not exist", name)
}

func (p *Environment) DeleteEnvironmentVariables(name string) error {
	for i, v := range p.EnvironmentVariables {
		if n, _ := v["name"]; strings.Compare(n.(string), name) == 0 {
			if i != (len(p.EnvironmentVariables) - 1) {
				p.EnvironmentVariables[i] = p.EnvironmentVariables[len(p.EnvironmentVariables)-1]
			}
			p.EnvironmentVariables = p.EnvironmentVariables[:(len(p.EnvironmentVariables) - 1)]
			return nil
		}
	}
	return fmt.Errorf("Env %s not exist", name)
}

type Environments struct {
	Links   Links `json:"_links"`
	Embeded struct {
		Environments []Environment `json:"environments"`
	} `json:"_embedded"`
}

func NewEnvironments() *Environments {
	return &Environments{Embeded: struct {
		Environments []Environment `json:"environments"`
	}{Environments: make([]Environment, 0)}}
}

type Version struct {
	Links struct {
		Self Link `json:"self"`
		Doc  Link `json:"doc"`
	} `json:"_links"`
	Version     string `json:"version"`
	BuildNumber string `json:"build_number"`
	GitSha      string `json:"git_sha"`
	FullVersion string `json:"full_version"`
	CommitUrl   string `json:"commit_url"`
}
