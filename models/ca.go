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

func (ca Ca) Terminal() string {
	lines := []string{}

	caBody := ca.CaBody

	if caBody.Value.Certificate != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Certificate:", caBody.Value.Certificate))
	}

	if caBody.Value.PrivateKey != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Private Key:", caBody.Value.PrivateKey))
	}

	return util.Header(caBody.ContentType, ca.Name) + strings.Join(lines, "\n") + "\n" + util.Footer(ca.CaBody.UpdatedAt)
}

func (ca Ca) Json() string {
	return ""
}
