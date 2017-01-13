package models

type CaBody struct {
	SecretType       string        `json:"type" binding:"required"`
	Name             string        `json:"name"`
	Value            *CaParameters `json:"value,omitempty" binding:"required"`
	VersionCreatedAt string        `json:"version_created_at,omitempty"`
}

func NewCaBody(bodyAsJsonObject map[string]interface{}) CaBody {
	var ca map[string]interface{}
	if data, ok := bodyAsJsonObject["data"].([]interface{}); ok {
		ca = data[0].(map[string]interface{})
	} else {
		ca = bodyAsJsonObject
	}
	value := ca["value"].(map[string]interface{})

	valueBody := CaParameters{
		Certificate: value["certificate"].(string),
		PrivateKey:  value["private_key"].(string),
	}
	caBody := CaBody{
		SecretType:       ca["type"].(string),
		VersionCreatedAt: ca["version_created_at"].(string),
		Value:            &valueBody,
	}
	if data, ok := ca["name"].(string); ok {
		caBody.Name = data
	}

	return caBody
}
