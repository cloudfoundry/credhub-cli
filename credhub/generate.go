package credhub

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generator"
)

func (ch CredHub) GeneratePassword(name string, gen generator.Password, overwrite bool) (credentials.Password, error) {
	panic("Not implemented")
}

func (ch CredHub) GenerateUser(name string, gen generator.User, overwrite bool) (credentials.User, error) {
	panic("Not implemented")
}

func (ch CredHub) GenerateCertificate(name string, gen generator.Certificate, overwrite bool) (credentials.Certificate, error) {
	panic("Not implemented")
}

func (ch CredHub) GenerateRSA(name string, gen generator.RSA, overwrite bool) (credentials.RSA, error) {
	panic("Not implemented")
}

func (ch CredHub) GenerateSSH(name string, gen generator.SSH, overwrite bool) (credentials.SSH, error) {
	panic("Not implemented")
}
