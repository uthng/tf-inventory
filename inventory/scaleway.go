package inventory

import (
	//"encoding/json"
	//"fmt"
	"strings"
)

// parseResourceVsphere loops primary attributes to find out
// hostname & ip of server and by the way, add all these attributes
// to Ansible hostvars list.
func parseResourceScaleway(res Resource, tfsParsed *TFStateParsed) string {
	hostvars := make(map[string]interface{})
	hostname := ""
	publicIP := ""
	privateIP := ""

	// Loop attributes
	for attrk, attrv := range res.Primary.Attributes {
		if strings.HasSuffix(attrk, "name") {
			hostname = attrv
		}

		if strings.HasSuffix(attrk, "public_ip") {
			publicIP = attrv
		}

		if strings.HasSuffix(attrk, "private_ip") {
			privateIP = attrv
		}

		hostvars[attrk] = attrv
	}

	if publicIP != "" {
		hostvars["ansible_ssh_host"] = publicIP
	} else {
		hostvars["ansible_ssh_host"] = privateIP
	}

	if hostname != "" {
		tfsParsed.Hosts[hostname] = &Host{
			Vars: hostvars,
		}
	}

	return hostname
}