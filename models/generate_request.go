package models

type GenerateRequest struct {
	Name           string                `json:"name"`
	CredentialType string                `json:"type"`
	Overwrite      *bool                 `json:"overwrite"`
	Parameters     *GenerationParameters `json:"parameters"`
}
