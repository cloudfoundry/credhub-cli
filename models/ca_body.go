package models

type CaBody struct {
	ContentType string        `json:"type" binding:"required"`
	Value       *CaParameters `json:"value,omitempty" binding:"required"`
	UpdatedAt   string        `json:"updated_at,omitempty"`
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
		ContentType: ca["type"].(string),
		UpdatedAt:   ca["updated_at"].(string),
		Value:       &valueBody,
	}

	return caBody
}
