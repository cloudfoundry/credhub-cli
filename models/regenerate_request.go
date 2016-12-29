package models

type RegenerateSecretRequest struct {
	Name       string `json:"name"`
	Regenerate bool   `json:"regenerate"`
}
