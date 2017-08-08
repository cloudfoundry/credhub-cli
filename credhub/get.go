package credhub

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
)

// Returns a credential by ID.
func (ch *CredHub) GetById(id string) (credentials.Credential, error) {
	panic("Not implemented")
}

// Returns all historical credential value(s) by name.
func (ch *CredHub) GetAll(name string) ([]credentials.Credential, error) {
	panic("Not implemented")
}

// Returns a credential by name.
func (ch *CredHub) Get(name string) (credentials.Credential, error) {
	panic("Not implemented")
}

// Returns a Value credential by name.
func (ch *CredHub) GetValue(name string) (credentials.Value, error) {
	panic("Not implemented")
}

// Returns a JSON credential by name.
func (ch *CredHub) GetJSON(name string) (credentials.JSON, error) {
	panic("Not implemented")
}

// Returns a Password credential by name.
func (ch *CredHub) GetPassword(name string) (credentials.Password, error) {
	panic("Not implemented yet")
}

// Returns a User credential by name.
func (ch *CredHub) GetUser(name string) (credentials.User, error) {
	panic("Not implemented")
}

// Returns a Certificate credential by name.
func (ch *CredHub) GetCertificate(name string) (credentials.Certificate, error) {
	panic("Not implemented")
}

// Returns an RSA credential by name.
func (ch *CredHub) GetRSA(name string) (credentials.RSA, error) {
	panic("Not implemented")
}

// Returns an SSH credential by name.
func (ch *CredHub) GetSSH(name string) (credentials.SSH, error) {
	panic("Not implemented")
}
