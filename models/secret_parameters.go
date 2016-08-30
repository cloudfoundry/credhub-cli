package models

type SecretParameters struct {
	Overwrite        bool     `json:"overwrite"`
	ExcludeSpecial   bool     `json:"exclude_special,omitempty"`
	ExcludeNumber    bool     `json:"exclude_number,omitempty"`
	ExcludeUpper     bool     `json:"exclude_upper,omitempty"`
	ExcludeLower     bool     `json:"exclude_lower,omitempty"`
	Length           int      `json:"length,omitempty"`
	CommonName       string   `json:"common_name,omitempty"`
	Organization     string   `json:"organization,omitempty"`
	OrganizationUnit string   `json:"organization_unit,omitempty"`
	Locality         string   `json:"locality,omitempty"`
	State            string   `json:"state,omitempty"`
	Country          string   `json:"country,omitempty"`
	AlternativeName  []string `json:"alternative_names,omitempty"`
	KeyLength        int      `json:"key_length,omitempty"`
	Duration         int      `json:"duration,omitempty"`
	Ca               string   `json:"ca,omitempty"`
}
