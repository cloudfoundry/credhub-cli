package models

type SecretBody struct {
	SecretType       string            `json:"type" binding:"required" yaml:"type"`
	Name             string            `json:"name,omitempty" yaml:"name,omitempty"`
	Value            interface{}       `json:"value,omitempty" yaml:"value,omitempty"`
	Overwrite        *bool             `json:"overwrite,omitempty" yaml:"overwrite,omitempty"`
	Parameters       *SecretParameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	VersionCreatedAt string            `json:"version_created_at,omitempty" yaml:"updated,omitempty"`
}

func NewSecretBody(bodyAsJsonObject map[string]interface{}) SecretBody {
	var secret map[string]interface{}
	if data, ok := bodyAsJsonObject["data"].([]interface{}); ok {
		secret = data[0].(map[string]interface{})
	} else {
		secret = bodyAsJsonObject
	}

	secretBody := SecretBody{
		Name:             secret["name"].(string),
		SecretType:       secret["type"].(string),
		VersionCreatedAt: secret["version_created_at"].(string),
	}

	switch secretBody.SecretType {
	case "value", "password":
		secretBody.Value = secret["value"].(string)
		break
	case "ssh", "rsa":
		value := secret["value"].(map[string]interface{})
		rsaSsh := RsaSsh{}
		rsaSsh.PublicKey, _ = value["public_key"].(string)
		rsaSsh.PrivateKey, _ = value["private_key"].(string)
		secretBody.Value = rsaSsh
		break
	case "certificate":
		value := secret["value"].(map[string]interface{})
		cert := Certificate{}
		cert.Ca, _ = value["ca"].(string)
		cert.Certificate, _ = value["certificate"].(string)
		cert.PrivateKey, _ = value["private_key"].(string)
		secretBody.Value = cert
		break
	}

	return secretBody
}
