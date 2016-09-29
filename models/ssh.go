package models

import "github.com/pivotal-cf/credhub-cli/util"

type Ssh struct {
	PublicKey  string `json:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

func (ssh Ssh) StringLines() []string {
	lines := []string{}
	if ssh.PublicKey != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Public Key:", ssh.PublicKey))
	}
	if ssh.PrivateKey != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Private Key:", ssh.PrivateKey))
	}
	return lines
}
