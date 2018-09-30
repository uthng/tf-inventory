package inventory

import (
	"encoding/json"
	//"fmt"
	"strings"

	"github.com/spf13/cast"
)

// TFState is a struct representing terraform tfstate
type TFState struct {
	Modules []Module `json:"modules"`
}

// Module represents module element in Terraform tfstate
type Module struct {
	Path      []string            `json:"path"`
	Outputs   map[string]Output   `json:"outputs"`
	Resources map[string]Resource `json:"resources"`
	DependsOn []string            `json:"depends_on"`
}

// Resource represents a resource element in Terraform tfstate
type Resource struct {
	Type      string   `json:"type"`
	DependsOn []string `json:"depends_on"`
	Primary   Primary  `json:"primary"`
	Deposed   []string `json:"deposed"`
	Provider  string   `json:"provider"`
}

// Primary represents a primary element in Terraform tfstate
type Primary struct {
	ID         string                 `json:"id"`
	Attributes map[string]string      `json:"attributes"`
	Meta       map[string]interface{} `json:"meta"`
	Tainted    bool                   `json:"tainted"`
}

// Output represents output characteristics in Terraform tfstate
type Output struct {
	Sensitive bool        `json:"sensitive"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
}

// HostVars is a map of hosts with their values
type HostVars map[string]interface{}

// Hosts is a list of hosts
type Hosts []string

// GroupHosts is a map of group contains a list of hosts
type GroupHosts map[string]Hosts

// Parse parse and fill up TFState struct
func (t *TFState) Parse(tfs []byte) (HostVars, GroupHosts, error) {
	hostVars := make(HostVars)
	groupHosts := make(GroupHosts)

	err := json.Unmarshal(tfs, t)
	if err != nil {
		return nil, nil, err
	}

	// Loop modules
	for _, module := range t.Modules {
		hvRes, ghRes := parseResources(module.Resources)

		hvOutputs := parseOutputs(module.Outputs)

		// Merge hostvars
		for name, vars := range hvRes {
			newvars := vars.(map[string]interface{})
			for k, v := range hvOutputs {
				newvars[k] = v
			}

			hostVars[name] = vars
		}

		// Merge group hosts
		for group, hosts := range ghRes {
			g, ok := groupHosts[group]
			if ok {
				groupHosts[group] = append(g, hosts...)
			} else {
				groupHosts[group] = hosts
			}
		}
	}

	return hostVars, groupHosts, err
}

//////////////// PRIVATE FUNCTIONS ///////////////////////:
func parseResources(resources map[string]Resource) (HostVars, GroupHosts) {
	hostVars := make(HostVars)
	groupHosts := make(GroupHosts)

	// Loop resources to find supported ones
	for resk, resv := range resources {
		if strings.HasPrefix(resk, "vsphere_virtual_machine.") {
			vars := make(map[string]interface{})
			hostname := ""
			group := ""

			// Loop attributes
			for attrk, attrv := range resv.Primary.Attributes {
				//if attrk == "name" {
				//host = attrv
				//}

				if strings.HasSuffix(attrk, "host_name") || strings.HasSuffix(attrk, "hostname") {
					hostname = attrv
				}

				if strings.HasSuffix(attrk, "ipv4_address") {
					vars["ansible_ssh_host"] = attrv
				}
			}

			hostVars[hostname] = vars

			// Get group name
			arr := strings.Split(resk, ".")
			group = arr[1]

			hosts, ok := groupHosts[group]
			if !ok {
				hosts = []string{}
			}

			hosts = append(hosts, hostname)
			groupHosts[group] = hosts
		}
	}

	return hostVars, groupHosts
}

func parseOutputs(outputs map[string]Output) map[string]interface{} {
	hostvars := make(map[string]interface{})

	// Loop outputs
	for outputk, outputv := range outputs {
		if outputv.Type == "string" {
			hostvars[outputk] = outputv.Value
		}

		if outputv.Type == "list" {
			hostvars[outputk] = cast.ToStringSlice(outputv.Value)
		}

		if outputv.Type == "map" {
			values := cast.ToStringMapString(outputv.Value)
			for k, v := range values {
				hostvars[k] = v
			}
		}
	}

	return hostvars
}
