package credhub

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
)

// Retrieves a credential by ID. Exactly one value will be returned.
func (ch CredHub) GetById(id string) (credentials.Credential, error) {
	panic("Not implemented")
}

// Returns a credential’s value by name (where the type is unknown)
func (ch CredHub) GetByName(name string) (credentials.Credential, error) {
	panic("Not implemented")
}

// Returns all historical credential value(s) by name
func (ch CredHub) GetAllByName(name string) ([]credentials.Credential, error) {
	panic("Not implemented")
}

// Returns a credential's value by name
func (ch CredHub) GetValue(name string) (credentials.Value, error) {
	panic("Not implemented")
}

// Returns a credential's JSON by name
func (ch CredHub) GetJSON(name string) (credentials.JSON, error) {
	panic("Not implemented")
}

// Returns a credential's Password by name
func (ch CredHub) GetPassword(name string) (credentials.Password, error) {
	panic("Not implemented")
}

// Returns a credential's User by name
func (ch CredHub) GetUser(name string) (credentials.User, error) {
	panic("Not implemented")
}

// Returns a credential's Certificate by name
func (ch CredHub) GetCertificate(name string) (credentials.Certificate, error) {
	panic("Not implemented")
}

// Returns a credential's RSA by name
func (ch CredHub) GetRSA(name string) (credentials.RSA, error) {
	panic("Not implemented")
}

// Returns a credential's SSH by name
func (ch CredHub) GetSSH(name string) (credentials.SSH, error) {
	panic("Not implemented")
}
