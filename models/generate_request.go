package models

type GenerateSecretRequest struct {
	SecretType string            `json:"type"`
	Overwrite  bool              `json:"overwrite"`
	Parameters *SecretParameters `json:"parameters"`
}

type GenerateCaRequest struct {
	SecretType string            `json:"type"`
	Name       string            `json:"name"`
	Parameters *SecretParameters `json:"parameters"`
}
