package models

import (
	"strings"

	"github.com/pivotal-cf/credhub-cli/util"
)

type Certificate struct {
	Ca          string `json:"ca,omitempty"`
	Certificate string `json:"certificate,omitempty"`
	PrivateKey  string `json:"private_key,omitempty"`
}

func (cert Certificate) String() string {
	lines := []string{}
	if cert.Ca != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Ca:", cert.Ca))
	}

	if cert.Certificate != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Certificate:", cert.Certificate))
	}

	if cert.PrivateKey != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Private Key:", cert.PrivateKey))
	}
	result := strings.Join(lines, "\n")
	if len(result) > 0 {
		result += "\n"
	}
	return result
}
