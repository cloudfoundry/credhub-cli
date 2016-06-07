package models

import "fmt"

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
	secretBody := secret.SecretBody
	if secretBody.ContentType == "value" {
		return fmt.Sprintf("Type:	value\nName:	%s\nValue:	%s",
			secret.Name,
			secretBody.Value)
	} else {
		return fmt.Sprintf("Type:		certificate\nName:		%s\nCA:		%s\nPublic:		%s\nPrivate:	%s",
			secret.Name,
			secretBody.Certificate.Ca,
			secretBody.Certificate.Public,
			secretBody.Certificate.Private)
	}
}
