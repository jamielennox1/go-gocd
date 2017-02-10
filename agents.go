package gocd

import (
	"reflect"
	"strings"
)

type Agent struct {
	Uuid            string `json:"uuid"`
	HostName        string `json:"hostname"`
	IpAddress       string `json:"ip_address"`
	Sandbox         string `json:"sandbox"`
	OperatingSystem string `json:"operating_system"`
	//FreeSpace        int   `json:"free_space,omitempty"`
	AgentConfigState string   `json:"agent_config_state"`
	AgentState       string   `json:"agent_state"`
	BuildState       string   `json:"build_state"`
	Resources        []string `json:"resources"`
	Environments     []string `json:"environments"`
}

func NewAgent() *Agent {
	return &Agent{
		Resources:    make([]string, 0),
		Environments: make([]string, 0)}
}

func (p Agent) Diff(agent Agent) map[string]interface{} {
	result := make(map[string]interface{})
	if strings.Compare(p.Uuid, agent.Uuid) != 0 {
		result["uuid"] = agent.Uuid
	}
	if strings.Compare(p.HostName, agent.HostName) != 0 {
		result["hostname"] = agent.HostName
	}
	if strings.Compare(p.IpAddress, agent.IpAddress) != 0 {
		result["ip_address"] = agent.IpAddress
	}
	if strings.Compare(p.Sandbox, agent.Sandbox) != 0 {
		result["sandbox"] = agent.Sandbox
	}
	if strings.Compare(p.OperatingSystem, agent.OperatingSystem) != 0 {
		result["operating_system"] = agent.OperatingSystem
	}
	if strings.Compare(p.AgentConfigState, agent.AgentConfigState) != 0 {
		result["agent_config_state"] = agent.AgentConfigState
	}
	if strings.Compare(p.AgentState, agent.AgentState) != 0 {
		result["agent_state"] = agent.AgentState
	}
	if strings.Compare(p.BuildState, agent.BuildState) != 0 {
		result["build_state"] = agent.BuildState
	}
	if reflect.DeepEqual(p.Resources, agent.Resources) {
		result["resources"] = agent.Resources
	}
	if reflect.DeepEqual(p.Environments, agent.Environments) {
		result["environments"] = agent.Environments
	}
	return result
}
