package credhub

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
)

func (ch CredHub) GeneratePassword(name string, gen generate.Password, overwrite bool) (credentials.Password, error) {
	panic("Not implemented")
}

func (ch CredHub) GenerateUser(name string, gen generate.User, overwrite bool) (credentials.User, error) {
	panic("Not implemented")
}

func (ch CredHub) GenerateCertificate(name string, gen generate.Certificate, overwrite bool) (credentials.Certificate, error) {
	panic("Not implemented")
}

func (ch CredHub) GenerateRSA(name string, gen generate.RSA, overwrite bool) (credentials.RSA, error) {
	panic("Not implemented")
}

func (ch CredHub) GenerateSSH(name string, gen generate.SSH, overwrite bool) (credentials.SSH, error) {
	panic("Not implemented")
}
