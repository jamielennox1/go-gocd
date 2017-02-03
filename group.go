package gocd

import (
	"strings"
)

type Group struct {
	Name      string `json:"name"`
	Pipelines []struct {
		Name string `json:"name"`
	} `json:"pipelines"`
}

func (p *Group) Exist(pipeline string) bool {
	for _, pp := range p.Pipelines {
		if strings.Compare(pp.Name, pipeline) == 0 {
			return true
		}
	}
	return false
}
