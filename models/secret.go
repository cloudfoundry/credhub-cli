package models

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
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

func (secret Secret) ToYaml() string {
	s, _ := yaml.Marshal(secret.SecretBody)
	return string(s)
}

func (secret Secret) ToJson() string {
	s, _ := json.MarshalIndent(secret.SecretBody, "", "\t")
	return string(s)
}
