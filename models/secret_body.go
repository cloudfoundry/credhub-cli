package models

type SecretBody struct {
	ContentType string            `json:"type" binding:"required"`
	Value       interface{}       `json:"value,omitempty"`
	Overwrite   bool              `json:"overwrite"`
	Parameters  *SecretParameters `json:"parameters,omitempty"`
	UpdatedAt   string            `json:"updated_at,omitempty"`
}

func NewSecretBody(m map[string]interface{}) SecretBody {
	secretBody := SecretBody{
		ContentType: m["type"].(string),
	}
	secretBody.UpdatedAt, _ = m["updated_at"].(string)

	switch secretBody.ContentType {
	case "value", "password":
		secretBody.Value = m["value"].(string)
		break
	case "ssh", "rsa":
		value := m["value"].(map[string]interface{})
		rsaSsh := RsaSsh{}
		rsaSsh.PublicKey, _ = value["public_key"].(string)
		rsaSsh.PrivateKey, _ = value["private_key"].(string)
		secretBody.Value = rsaSsh
		break
	case "certificate":
		value := m["value"].(map[string]interface{})
		cert := Certificate{}
		cert.Ca, _ = value["ca"].(string)
		cert.Certificate, _ = value["certificate"].(string)
		cert.PrivateKey, _ = value["private_key"].(string)
		secretBody.Value = cert
		break
	}

	return secretBody
}
