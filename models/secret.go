package models

import (
	"encoding/json"

	"github.com/pivotal-cf/credhub-cli/util"
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

func (s Secret) Terminal() string {
	result := ""
	secretBody := s.SecretBody

	switch secretBody.ContentType {
	case "value", "password":
		result = util.BuildLineOfFixedLength("Value:", secretBody.Value.(string)) + "\n"
		break
	case "certificate":
		cert := Certificate{}
		json.Unmarshal(marshalBackIntoJson(secretBody.Value.(map[string]interface{})), &cert)
		result = cert.Terminal()
		break
	case "ssh", "rsa":
		rsaSsh := RsaSsh{}
		json.Unmarshal(marshalBackIntoJson(secretBody.Value.(map[string]interface{})), &rsaSsh)
		result = rsaSsh.Terminal()
		break
	}

	return util.Header(secretBody.ContentType, s.Name) + result + util.Footer(secretBody.UpdatedAt)
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
		Type:      secretBody.ContentType,
		UpdatedAt: secretBody.UpdatedAt,
	}

	switch secretBody.ContentType {
	case "value", "password":
		body.Value = secretBody.Value.(string)
		break
	case "certificate":
		certificate := Certificate{}
		json.Unmarshal(marshalBackIntoJson(secretBody.Value.(map[string]interface{})), &certificate)
		body.Ca = certificate.Ca
		body.Certificate = certificate.Certificate
		body.PrivateKey = certificate.PrivateKey
		break
	case "ssh", "rsa":
		rsaSsh := RsaSsh{}
		json.Unmarshal(marshalBackIntoJson(secretBody.Value.(map[string]interface{})), &rsaSsh)
		body.PublicKey = rsaSsh.PublicKey
		body.PrivateKey = rsaSsh.PrivateKey
		break
	}

	s, _ := json.MarshalIndent(body, "", JSON_PRETTY_PRINT_INDENT_STRING)
	return string(s)
}

func marshalBackIntoJson(value map[string]interface{}) []byte {
	item, _ := json.Marshal(value)
	return item
}
