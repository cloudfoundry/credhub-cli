// CredHub credential types for regenerating credentials
package regenerate

type Password struct{}

type User struct{}

type Certificate struct {
	KeyLength int `json:"key_length"`
}

type RSA struct{}

type SSH struct{}
