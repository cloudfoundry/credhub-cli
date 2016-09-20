package models

import (
	"strings"
)

type Ca struct {
	Name   string
	CaBody CaBody
}

func NewCa(name string, caBody CaBody) Item {
	return Ca{
		Name:   name,
		CaBody: caBody,
	}
}

func (ca Ca) String() string {
	lines := []string{}

	caBody := ca.CaBody
	lines = append(lines,
		buildLineOfFixedLength("Type:", caBody.ContentType),
		buildLineOfFixedLength("Name:", ca.Name),
	)

	if caBody.Value.Certificate != "" {
		lines = append(lines, buildLineOfFixedLength("Certificate:", caBody.Value.Certificate))
	}

	if caBody.Value.PrivateKey != "" {
		lines = append(lines, buildLineOfFixedLength("Private Key:", caBody.Value.PrivateKey))
	}

	if caBody.UpdatedAt != "" {
		lines = append(lines, buildLineOfFixedLength("Updated:", caBody.UpdatedAt))
	}

	return strings.Join(lines, "\n")
}
