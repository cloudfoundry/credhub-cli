package models

type SecretBody struct {
	ContentType string      `json:"type" binding:"required"`
	Credential  interface{} `json:"credential,omitempty"`
	UpdatedAt   string      `json:"updated_at,omitempty"`
}
