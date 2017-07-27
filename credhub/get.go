package credhub

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
)

func (ch CredHub) GetByName(name string) (credentials.Credential, error) {
	panic("Not implemented")
}

func (ch CredHub) GetValueByName(name string) (credentials.Value, error) {
	panic("Not implemented")
}

func (ch CredHub) GetJSONByName(name string) (credentials.JSON, error) {
	panic("Not implemented")
}

func (ch CredHub) GetPasswordByName(name string) (credentials.Password, error) {
	panic("Not implemented")
}

func (ch CredHub) GetUserByName(name string) (credentials.User, error) {
	panic("Not implemented")
}

func (ch CredHub) GetCertificateByName(name string) (credentials.Certificate, error) {
	panic("Not implemented")
}

func (ch CredHub) GetRSAByName(name string) (credentials.RSA, error) {
	panic("Not implemented")
}

func (ch CredHub) GetSSHByName(name string) (credentials.SSH, error) {
	panic("Not implemented")
}

func (ch CredHub) GetAllByName(name string) ([]credentials.Credential, error) {
	panic("Not implemented")
}

func (ch CredHub) GetById(id string) (credentials.Credential, error) {
	panic("Not implemented")
}
