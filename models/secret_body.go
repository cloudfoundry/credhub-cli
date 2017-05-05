package models

type SecretBody struct {
	SecretType       string            `json:"type" binding:"required"`
	Name             string            `json:"name,omitempty"`
	Value            interface{}       `json:"value,omitempty"`
	Overwrite        *bool             `json:"overwrite,omitempty"`
	Parameters       *SecretParameters `json:"parameters,omitempty"`
	VersionCreatedAt string            `json:"version_created_at,omitempty"`
}
