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
	return fmt.Sprintf("Type:	%s\nName:	%s\nValue:	%s",
		secret.SecretBody.ContentType, secret.Name, secret.SecretBody.Value)
}
