package models

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/util"
)

type Secret struct {
	SecretBody SecretBody
}

func NewSecret(secretBodyMap map[string]interface{}) Secret {
	secretBody := NewSecretBody(secretBodyMap)

	return Secret{
		SecretBody: secretBody,
	}
}

func (s Secret) Terminal() string {
	result := ""
	secretBody := s.SecretBody

	switch secretBody.SecretType {
	case "value", "password":
		result = util.BuildLineOfFixedLength("Value:", secretBody.Value.(string)) + "\n"
		break
	case "certificate":
		result = secretBody.Value.(Certificate).Terminal()
		break
	case "ssh", "rsa":
		result = secretBody.Value.(RsaSsh).Terminal()
		break
	}

	return util.Header(secretBody.SecretType, secretBody.Name) + result + util.Footer(secretBody.VersionCreatedAt)
}

func (secret Secret) Json() string {
	s, _ := json.MarshalIndent(secret.SecretBody, "", JSON_PRETTY_PRINT_INDENT_STRING)
	return string(s)
}
