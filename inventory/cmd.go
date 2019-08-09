package inventory

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	consul "github.com/hashicorp/consul/api"
)

// ExecCmd executes one of available commands: list, host and inventory.
// It checks if a tfstate is given in command line argument or not.
// If yes, it will use it. Otherwise, it checks 2 environment variables
// TF_BIN (path to terraform binary) and TF_DIR (path to terraform folder).
// If they are set, so `terraform state pull` is executed automatically to
// get Terrafrm tftstate to parse.
func ExecCmd(cmd string, arg string, tfsFile string) (interface{}, error) {
	var tfsBytes []byte
	var err error
	var out bytes.Buffer

	tfBin := os.Getenv("TF_BIN")
	tfDir := os.Getenv("TF_DIR")

	consulURL, consulURLSet := os.LookupEnv("CONSUL_URL")
	consulToken, _ := os.LookupEnv("CONSUL_TOKEN")
	consulKey, consulKeySet := os.LookupEnv("CONSUL_KEY")

	// If a tfstate is specified in commandline
	if tfsFile != "" {
		tfsBytes, err = readFile(tfsFile)
		if err != nil {
			return nil, err
		}
	} else {
		// Check if variables for Consul are set.
		// If yes, we must use it. Otherwise, use TF_BIN & TF_DIR
		// by default
		if consulURLSet && consulKeySet {
			tfsBytes, err = getTfStateFromConsul(consulURL, consulToken, consulKey)
			if err != nil {
				return nil, err
			}
		} else {
			// Check env variables to do a 'terraform state pull'
			if tfBin == "" {
				tfBin = "terraform"
			}

			// Execute cmd
			cmd := exec.Command("terraform", "state", "pull")
			if tfDir != "" {
				cmd.Dir = tfDir
			}

			cmd.Stdout = &out

			err = cmd.Run()
			if err != nil {
				return nil, fmt.Errorf("error `terraform state pull` in directory %s, %s", tfDir, err)
			}

			// read statefile contents
			tfsBytes, err = ioutil.ReadAll(&out)
			if err != nil {
				return nil, err
			}
		}
	}

	// Parse tfstate
	tfs := &TFState{}

	tfsParsed, err := tfs.Parse(tfsBytes)
	if err != nil {
		return nil, err
	}

	switch cmd {
	case "list":
		res := cmdList(tfsParsed)
		return res, nil
	case "host":
		res := cmdHost(arg, tfsParsed.Hosts)
		return res, nil
	case "inventory":
		return nil, nil
	default:
		return nil, fmt.Errorf("command %s unknown", cmd)
	}
}

////////////// PRIVATE FUNCTIONS /////////////////////////////////:
func readFile(filePath string) ([]byte, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(filePath)
}

func cmdList(tfsParsed *TFStateParsed) map[string]interface{} {
	inventory := make(map[string]interface{})
	groupNames := []string{}

	hosts := tfsParsed.Hosts
	groups := tfsParsed.Groups

	// Format group "_meta"
	inventory["_meta"] = map[string]interface{}{
		"hostvars": hosts,
	}

	// Format list of groups
	for name, group := range groups {
		inventory[name] = group
		groupNames = append(groupNames, name)
	}

	// Format group "all"
	groupNames = append(groupNames, "ungrouped")
	inventory["all"] = map[string]interface{}{
		"children": groupNames,
	}

	// Format group "ungrouped"
	inventory["ungrouped"] = map[string]interface{}{}

	return inventory
}

func cmdHost(host string, hosts Hosts) *Host {
	return hosts[host]
}

func getTfStateFromConsul(url, token, key string) ([]byte, error) {

	config := consul.DefaultConfig()
	config.Address = url
	if token != "" {
		config.Token = token
	}

	client, err := consul.NewClient(config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	pair, _, err := client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}

	return pair.Value, nil
}
