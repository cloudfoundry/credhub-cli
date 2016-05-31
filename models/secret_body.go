package models

type SecretBody struct {
	ContentType string `json:"type"`
	Value       string `json:"value"`
}

func NewSecretBody(value string) SecretBody {
	return SecretBody{
		ContentType: "value",
		Value:       value,
	}
}
