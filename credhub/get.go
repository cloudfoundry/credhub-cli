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

// Returns a credential's value by name
func (ch CredHub) GetValueByName(name string) (credentials.Value, error) {
	panic("Not implemented")
}

// Returns a credential's JSON by name
func (ch CredHub) GetJSONByName(name string) (credentials.JSON, error) {
	panic("Not implemented")
}

// Returns a credential's Password by name
func (ch CredHub) GetPasswordByName(name string) (credentials.Password, error) {
	panic("Not implemented")
}

// Returns a credential's User by name
func (ch CredHub) GetUserByName(name string) (credentials.User, error) {
	panic("Not implemented")
}

// Returns a credential's Certificate by name
func (ch CredHub) GetCertificateByName(name string) (credentials.Certificate, error) {
	panic("Not implemented")
}

// Returns a credential's RSA by name
func (ch CredHub) GetRSAByName(name string) (credentials.RSA, error) {
	panic("Not implemented")
}

// Returns a credential's SSH by name
func (ch CredHub) GetSSHByName(name string) (credentials.SSH, error) {
	panic("Not implemented")
}

// Returns all historical credential value(s) by name
func (ch CredHub) GetAllByName(name string) ([]credentials.Credential, error) {
	panic("Not implemented")
}
