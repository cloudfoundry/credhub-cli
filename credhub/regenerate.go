package credhub

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"

// Regenerates a credential using the same parameters as the stored value.
func (ch CredHub) Regenerate(name string) (credentials.Credential, error) {
	panic("Not implemented")
}

// Regenerates a Password using the same parameters as the stored value.
func (ch CredHub) RegeneratePassword(name string) (credentials.Password, error) {
	panic("Not implemented")
}

// Regenerates a User using the same parameters as the stored value.
func (ch CredHub) RegenerateUser(name string) (credentials.User, error) {
	panic("Not implemented")
}

// Regenerates a Certficate using the same parameters as the stored value.
func (ch CredHub) RegenerateCertificate(name string) (credentials.Certificate, error) {
	panic("Not implemented")
}

// Regenerates an RSA using the same parameters as the stored value.
func (ch CredHub) RegenerateRSA(name string) (credentials.RSA, error) {
	panic("Not implemented")
}

// Regenerates an SSH using the same parameters as the stored value.
func (ch CredHub) RegenerateSSH(name string) (credentials.SSH, error) {
	panic("Not implemented")
}
