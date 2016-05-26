package models

type SecretParameters struct {
	Length       int  `json:"length,omitempty"`
	ExcludeUpper bool `json:"exclude_upper,omitempty"`
}
