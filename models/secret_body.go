package models

type SecretBody struct {
	ContentType string       `json:"type" binding:"required"`
	Value       string       `json:"value,omitempty"`
	Certificate *Certificate `json:"certificate,omitempty"`
}
