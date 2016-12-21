package models

import (
	"encoding/json"

	"github.com/pivotal-cf/credhub-cli/util"
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

	return util.Header(secretBody.SecretType, secretBody.Name) + result + util.Footer(secretBody.UpdatedAt)
}

func (secret Secret) Json() string {
	secretBody := secret.SecretBody

	body := struct {
		Value       string `json:"value,omitempty"`
		Ca          string `json:"ca,omitempty"`
		Certificate string `json:"certificate,omitempty"`
		PrivateKey  string `json:"private_key,omitempty"`
		PublicKey   string `json:"public_key,omitempty"`
		Type        string `json:"type"`
		UpdatedAt   string `json:"updated_at"`
	}{
		Type:      secretBody.SecretType,
		UpdatedAt: secretBody.UpdatedAt,
	}

	switch secretBody.SecretType {
	case "value", "password":
		body.Value = secretBody.Value.(string)
		break
	case "certificate":
		certificate := secretBody.Value.(Certificate)
		body.Ca = certificate.Ca
		body.Certificate = certificate.Certificate
		body.PrivateKey = certificate.PrivateKey
		break
	case "ssh", "rsa":
		rsaSsh := secretBody.Value.(RsaSsh)
		body.PublicKey = rsaSsh.PublicKey
		body.PrivateKey = rsaSsh.PrivateKey
		break
	}

	s, _ := json.MarshalIndent(body, "", JSON_PRETTY_PRINT_INDENT_STRING)
	return string(s)
}
