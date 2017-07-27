package credhub

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"

func (ch CredHub) Regenerate(name string) (credentials.Credential, error) {
	panic("Not implemented")
}

func (ch CredHub) RegeneratePassword(name string) (credentials.Password, error) {
	panic("Not implemented")
}

func (ch CredHub) RegenerateUser(name string) (credentials.User, error) {
	panic("Not implemented")
}

func (ch CredHub) RegenerateCertificate(name string) (credentials.Certificate, error) {
	panic("Not implemented")
}

func (ch CredHub) RegenerateRSA(name string) (credentials.RSA, error) {
	panic("Not implemented")
}

func (ch CredHub) RegenerateSSH(name string) (credentials.SSH, error) {
	panic("Not implemented")
}
