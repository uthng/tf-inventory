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

// Group represents an Ansible inventory group
type Group struct {
	Name  string                 `json:"-"`
	Vars  map[string]interface{} `json:"vars"`
	Hosts []string               `json:"hosts"`
}

// Host represents an Ansible inventory host
type Host struct {
	Name      string                 `json:"-"`
	GroupName string                 `json:"-"`
	Vars      map[string]interface{} `json:"vars"`
}

// Hosts is map representing a list of hosts
type Hosts map[string]*Host

// Groups is map representing a list of groups
type Groups map[string]*Group

// TFStateParsed represents a tfstate under Ansible inventory format
type TFStateParsed struct {
	Hosts  Hosts
	Groups Groups
}

// Parse parse and fill up TFState struct
func (t *TFState) Parse(tfs []byte) (*TFStateParsed, error) {
	tfsParsed := &TFStateParsed{
		Hosts:  make(Hosts),
		Groups: make(Groups),
	}

	//groupName := ""

	err := json.Unmarshal(tfs, t)
	if err != nil {
		return nil, err
	}

	// Loop modules
	for _, module := range t.Modules {
		// Parse Module Resources
		//parseResources(groupName, module.Resources, tfsParsed)

		// Parse Module Outputs
		//parseOutputs(groupName, module.Outputs, tfsParsed)

		hosts, groups := parseModule(module)

		for name, host := range hosts {
			tfsParsed.Hosts[name] = host
		}

		for name, group := range groups {
			tfsParsed.Groups[name] = group
		}
	}

	return tfsParsed, err
}

//////////////// PRIVATE FUNCTIONS ///////////////////////
func parseModule(module Module) (Hosts, Groups) {
	groupvars := make(map[string]interface{})
	groups := make(Groups)
	groupName := ""

	// if Path length > 1, it means that the current
	// module is not a root module => tfstate generated uses module/submodule
	// => 2nd element in Path correspond to the name of module declared
	// in tf files with the following instruction:
	// module "uthng-blog" {
	// ....
	// }
	// Case "Root" module => different ressources
	// Merge group hosts
	if len(module.Path) > 1 {
		groupName = module.Path[1]
	}

	hosts := parseResources(groupName, module.Resources)
	if len(hosts) > 0 {
		groupvars = parseOutputs(module.Outputs)
	}

	for _, host := range hosts {
		group, ok := groups[host.GroupName]
		if !ok {
			group = &Group{
				Name: host.GroupName,
				Vars: groupvars,
			}
		}

		group.Hosts = append(group.Hosts, host.Name)
		groups[host.GroupName] = group
	}

	return hosts, groups
}

func parseResources(groupName string, resources map[string]Resource) Hosts {
	hosts := make(Hosts)

	// Loop resources to find supported ones
	for resk, resv := range resources {
		//hostname := ""
		var host *Host

		host = nil
		//groupname := groupName

		// Parse according to type
		if resv.Type == "vsphere_virtual_machine" {
			host = parseResourceVsphere(resv)
		} else if resv.Type == "scaleway_server" {
			host = parseResourceScaleway(resv)
		}

		// Add to hostvars or grouphost only if resource is
		// supported and return a hostname with its ip
		if host != nil {
			// In case of root path (no module used in Terraform)
			if groupName == "" {
				// Get group name
				arr := strings.Split(resk, ".")
				host.GroupName = arr[1]
			} else {
				host.GroupName = groupName
			}

			hosts[host.Name] = host

			//group, ok := tfsParsed.Groups[groupname]
			//if !ok {
			//group = &Group{
			//Vars:  make(map[string]interface{}),
			//Hosts: []string{},
			//}
			//}

			//group.Hosts = append(group.Hosts, hostname)
			//tfsParsed.Groups[groupname] = group
		}
	}

	return hosts
}

// Parse Module Outputs in tfstate and use the result as
// group vars in ansible inventory
func parseOutputs(outputs map[string]Output) map[string]interface{} {
	groupvars := make(map[string]interface{})

	// Loop outputs
	for outputk, outputv := range outputs {
		if outputv.Type == "string" {
			groupvars[outputk] = outputv.Value
		}

		if outputv.Type == "list" {
			groupvars[outputk] = cast.ToStringSlice(outputv.Value)
		}

		if outputv.Type == "map" {
			values := cast.ToStringMapString(outputv.Value)
			for k, v := range values {
				groupvars[k] = v
			}
		}
	}

	// Case of use of Module/Submodule
	// add groupvars to that specified group
	// Otherwise, no terraform module in use (root)
	// output variables will be for every groups
	//if groupName != "" {
	//group, ok := tfsParsed.Groups[groupName]
	//if ok {
	//fmt.Println("groupvars", groupvars)
	//fmt.Println("groupName", groupName)
	//group.Vars = groupvars
	//}
	//} else {
	//for _, group := range tfsParsed.Groups {
	//group.Vars = groupvars
	//}
	//}

	return groupvars
}
