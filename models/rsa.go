package models

import (
	"strings"

	"github.com/pivotal-cf/credhub-cli/util"
)

type Rsa struct {
	PublicKey  string `json:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

func (rsa Rsa) String() string {
	lines := []string{}
	if rsa.PublicKey != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Public Key:", rsa.PublicKey))
	}
	if rsa.PrivateKey != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Private Key:", rsa.PrivateKey))
	}
	return strings.Join(lines, "\n") + "\n"
}
