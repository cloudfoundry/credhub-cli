package models

import (
	"strings"

	"github.com/cloudfoundry-incubator/credhub-cli/util"
)

type RsaSsh struct {
	PublicKey  string `json:"public_key,omitempty" yaml:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty" yaml:"private_key,omitempty"`
}

func (rsaSsh RsaSsh) ToYaml() string {
	lines := []string{}
	if rsaSsh.PublicKey != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Public Key:", rsaSsh.PublicKey))
	}
	if rsaSsh.PrivateKey != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Private Key:", rsaSsh.PrivateKey))
	}
	return strings.Join(lines, "\n") + "\n"
}
