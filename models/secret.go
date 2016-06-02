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
	if secret.SecretBody.ContentType == "value" {
		return fmt.Sprintf("Type:	%s\nName:	%s\nValue:	%s",
			secret.SecretBody.ContentType, secret.Name, secret.SecretBody.Value)
	} else {
		return fmt.Sprintf("Type:	certificate\nName:	%s\nCA:	%s\nPublic:	%s\nPrivate:	%s",
			secret.Name,
			secret.SecretBody.Certificate.Ca,
			secret.SecretBody.Certificate.Public,
			secret.SecretBody.Certificate.Private)
	}
}
