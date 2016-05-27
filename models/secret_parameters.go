package models

type SecretParameters struct {
	ExcludeSpecial bool `json:"exclude_special,omitempty"`
	ExcludeNumber  bool `json:"exclude_number,omitempty"`
	ExcludeUpper   bool `json:"exclude_upper,omitempty"`
	ExcludeLower   bool `json:"exclude_lower,omitempty"`
	Length         int  `json:"length,omitempty"`
}
