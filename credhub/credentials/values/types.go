// CredHub credential value types
package values

type Value string

type JSON interface{}

type Password string

type User struct {
	Username string
	Password string
}

type Certificate struct {
	Ca          string
	Certificate string
	PrivateKey  string
}

type RSA struct {
	PublicKey  string
	PrivateKey string
}

type SSH struct {
	PublicKey            string
	PublicKeyFingerprint string
	PrivateKey           string
}
