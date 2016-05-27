package models

type SecretParameters struct {
	Length       int  `json:"length,omitempty"`
	ExcludeUpper bool `json:"exclude_upper,omitempty"`
	ExcludeLower bool `json:"exclude_lower,omitempty"`
}
