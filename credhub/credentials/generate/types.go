// CredHub credential types for generating credentials
package generate

type Password struct {
	Length         int  `json:"length,omitempty"`
	IncludeSpecial bool `json:"include_special,omitempty"`
	ExcludeNumber  bool `json:"exclude_number,omitempty"`
	ExcludeUpper   bool `json:"exclude_upper,omitempty"`
	ExcludeLower   bool `json:"exclude_lower,omitempty"`
}

type User struct {
	Username       string
	Length         int
	IncludeSpecial bool
	ExcludeNumber  bool
	ExcludeUpper   bool
	ExcludeLower   bool
}

type Certificate struct {
	KeyLength        int
	Duration         int
	CommonName       string
	Organization     string
	OrganizationUnit string
	Locality         string
	State            string
	Country          string
	AlternativeName  []string
	KeyUsage         []string
	ExtendedKeyUsage []string
	Ca               string `json:"ca"`
	IsCA             bool
	SelfSign         bool
}

type RSA struct {
	KeyLength int
}

type SSH struct {
	SshComment string
	KeyLength  int
}
