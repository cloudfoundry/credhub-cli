package models

type SecretBody struct {
	ContentType string           `json:"type" binding:"required"`
	Value       interface{}      `json:"value,omitempty"`
	Parameters  SecretParameters `json:"parameters,omitEmpty"`
	UpdatedAt   string           `json:"updated_at,omitempty"`
}
