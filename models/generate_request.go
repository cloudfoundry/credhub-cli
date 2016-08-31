package models

type GenerateRequest struct {
	ContentType  string           `json:"type"`
	Parameters  *SecretParameters `json:"parameters"`
}
