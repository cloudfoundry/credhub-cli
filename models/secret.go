package models

import (
	"fmt"
	"strings"
)

type Secret struct {
	Name       string
	SecretBody SecretBody
}

func NewSecret(name string, secretBody SecretBody) Secret {
	return Secret{
		Name:       name,
		SecretBody: secretBody,
	}
}

func (secret Secret) String() string {
	lines := []string{}

	secretBody := secret.SecretBody
	lines = append(lines,
		fmt.Sprintf("Type:		%s", secretBody.ContentType),
		fmt.Sprintf("Name:		%s", secret.Name),
	)

	if secretBody.ContentType == "value" {
		lines = append(lines, fmt.Sprintf("Value:		%s", secretBody.Value))
	} else {
		if secretBody.Certificate.Ca != "" {
			lines = append(lines, fmt.Sprintf("CA:		%s", secretBody.Certificate.Ca))
		}

		if secretBody.Certificate.Public != "" {
			lines = append(lines, fmt.Sprintf("Public:		%s", secretBody.Certificate.Public))
		}

		if secretBody.Certificate.Private != "" {
			lines = append(lines, fmt.Sprintf("Private:	%s", secretBody.Certificate.Private))
		}
	}

	return strings.Join(lines, "\n")
}
