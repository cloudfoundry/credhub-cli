package models

import (
	"fmt"
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
		fmt.Sprintf("Name:		%s", ca.Name),
	)

	if caBody.Ca.Public != "" {
		lines = append(lines, fmt.Sprintf("Public:		%s", caBody.Ca.Public))
	}

	if caBody.Ca.Private != "" {
		lines = append(lines, fmt.Sprintf("Private:	%s", caBody.Ca.Private))
	}

	lines = append(lines, fmt.Sprintf("Updated:	%s", caBody.UpdatedAt))

	return strings.Join(lines, "\n")
}
