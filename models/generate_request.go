package models

type GenerateSecretRequest struct {
	ContentType string            `json:"type"`
	Overwrite   bool              `json:"overwrite"`
	Parameters  *SecretParameters `json:"parameters"`
}

type GenerateCaRequest struct {
	ContentType string            `json:"type"`
	Parameters  *SecretParameters `json:"parameters"`
}
