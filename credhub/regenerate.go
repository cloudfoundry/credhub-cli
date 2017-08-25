package credhub

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"

// Regenerates a credential.
func (ch *CredHub) Regenerate(name string) (credentials.Credential, error) {
	panic("Not implemented")
}

// Regenerates a Password credential.
func (ch *CredHub) RegeneratePassword(name string) (credentials.Password, error) {
	panic("Not implemented")
}

// Regenerates a User credential.
func (ch *CredHub) RegenerateUser(name string) (credentials.User, error) {
	panic("Not implemented")
}

// Regenerates a Certificate credential.
func (ch *CredHub) RegenerateCertificate(name string) (credentials.Certificate, error) {
	panic("Not implemented")
}

// Regenerates an RSA credential.
func (ch *CredHub) RegenerateRSA(name string) (credentials.RSA, error) {
	panic("Not implemented")
}

// Regenerates an SSH credential.
func (ch *CredHub) RegenerateSSH(name string) (credentials.SSH, error) {
	panic("Not implemented")
}
