package models

type GenerateSecretRequest struct {
	ContentType string            `json:"type"`
	Overwrite   bool              `json:"overwrite"`
	Parameters  *SecretParameters `json:"parameters"`
}

type GenerateCaRequest struct {
	ContentType string            `json:"type"`
	Name        string            `json:"name"`
	Parameters  *SecretParameters `json:"parameters"`
}
