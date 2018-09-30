package inventory

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

	// If a tfstate is specified in commandline
	if tfsFile != "" {
		tfsBytes, err = readFile(tfsFile)
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

	// Parse tfstate
	tfs := &TFState{}

	hostvars, grouphosts, err := tfs.Parse(tfsBytes)
	if err != nil {
		return nil, err
	}

	switch cmd {
	case "list":
		res := cmdList(hostvars, grouphosts)
		return res, nil
	case "host":
		res := cmdHost(arg, hostvars)
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

func cmdList(hv HostVars, gh GroupHosts) map[string]interface{} {
	inventory := make(map[string]interface{})
	groups := []string{}

	// Format group "_meta"
	inventory["_meta"] = map[string]interface{}{
		"hostvars": hv,
	}

	// Format list of groups
	for g, h := range gh {
		inventory[g] = map[string]interface{}{
			"hosts": h,
		}

		groups = append(groups, g)
	}

	// Format group "all"
	groups = append(groups, "ungrouped")
	inventory["all"] = map[string]interface{}{
		"children": groups,
	}

	// Format group "ungrouped"
	inventory["ungrouped"] = map[string]interface{}{}

	return inventory
}

func cmdHost(host string, hv HostVars) map[string]interface{} {
	return hv[host].(map[string]interface{})
}
