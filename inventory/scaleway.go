package inventory

import (
	//"encoding/json"
	//"fmt"
	"strings"
)

func parseResourceScaleway(res Resource) (string, string) {
	hostname := ""
	publicIP := ""
	privateIP := ""
	ip := ""

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
	}

	if publicIP != "" {
		ip = publicIP
	} else {
		ip = privateIP
	}

	return hostname, ip
}
