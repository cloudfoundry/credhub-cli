package models

import "github.com/pivotal-cf/credhub-cli/util"

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

func (s Secret) String() string {
	return util.Header(s.SecretBody.ContentType, s.Name) + s.SecretBody.String() + util.Footer(s.SecretBody.UpdatedAt)
}
