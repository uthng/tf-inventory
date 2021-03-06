package inventory

import (
	//"encoding/json"
	//"fmt"
	"strings"
)

// parseResourceVsphere loops primary attributes to find out
// hostname & ip of server and by the way, add all these attributes
// to Ansible hostvars list.
func parseResourceVsphere(res Resource) *Host {
	var host *Host

	hostvars := make(map[string]interface{})
	hostname := ""
	ip := ""
	host = nil

	// Loop attributes
	for attrk, attrv := range res.Primary.Attributes {
		if strings.HasSuffix(attrk, "host_name") || strings.HasSuffix(attrk, "hostname") {
			hostname = attrv
		}

		if strings.HasSuffix(attrk, "ipv4_address") {
			ip = attrv
		}

		hostvars[attrk] = attrv
	}

	if ip != "" {
		hostvars["ansible_ssh_host"] = ip
	}

	if hostname != "" {
		host = &Host{
			Name: hostname,
			Vars: hostvars,
		}
	}

	return host
}
