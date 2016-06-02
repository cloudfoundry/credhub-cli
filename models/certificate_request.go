package models

type CertificateRequest struct {
	ContentType string `json:"type"`

	Certificate Certificate `json:"certificate"`
}
