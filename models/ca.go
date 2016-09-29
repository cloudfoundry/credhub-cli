package models

import (
	"strings"
	"github.com/pivotal-cf/credhub-cli/util"
)

type Ca struct {
	Name   string
	CaBody CaBody
}

func NewCa(name string, caBody CaBody) Ca {
	return Ca{
		Name:   name,
		CaBody: caBody,
	}
}

func (ca Ca) String() string {
	lines := []string{}

	caBody := ca.CaBody
	lines = append(lines,
		util.BuildLineOfFixedLength("Type:", caBody.ContentType),
		util.BuildLineOfFixedLength("Name:", ca.Name),
	)

	if caBody.Value.Certificate != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Certificate:", caBody.Value.Certificate))
	}

	if caBody.Value.PrivateKey != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Private Key:", caBody.Value.PrivateKey))
	}

	if caBody.UpdatedAt != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Updated:", caBody.UpdatedAt))
	}

	return strings.Join(lines, "\n")
}
