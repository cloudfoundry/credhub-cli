package models

type GenerateSecretRequest struct {
	Name       string                `json:"name"`
	SecretType string                `json:"type"`
	Overwrite  *bool                 `json:"overwrite"`
	Parameters *GenerationParameters `json:"parameters"`
}

type GenerateCaRequest struct {
	SecretType string                `json:"type"`
	Name       string                `json:"name"`
	Parameters *GenerationParameters `json:"parameters"`
}
