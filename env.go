package gocd

import (
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

func (p *Environment) AddPipeline(pipeline string) bool {
	for _, p := range p.Pipelines {
		if strings.Compare(p.Name, pipeline) == 0 {
			return false
		}
	}
	p.Pipelines = append(p.Pipelines, ShortPipeline{Name: pipeline})
	return true
}

func (p *Environment) ExistPipeline(pipeline string) bool {
	for _, p := range p.Pipelines {
		if strings.Compare(p.Name, pipeline) == 0 {
			return true
		}
	}
	return false
}

func (p *Environment) DeletePipeline(pipeline string) bool {
	for i := 0; i < len(p.Pipelines); i++ {
		if strings.Compare(p.Pipelines[i].Name, pipeline) == 0 {
			p.Pipelines[i] = p.Pipelines[len(p.Pipelines)-1]
			p.Pipelines = p.Pipelines[:len(p.Pipelines)-1]
			return true
		}
	}
	return false
}

type Environments struct {
	Links   Links `json:"_links"`
	Embeded struct {
		Environments []Environment `json:"environments"`
	} `json:"_embedded"`
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
