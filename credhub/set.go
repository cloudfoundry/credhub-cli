package credhub

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

func (ch CredHub) SetValue(name string, value values.Value, overwrite bool) (credentials.Value, error) {
	panic("Not implemented")
}

func (ch CredHub) SetJSON(name string, value values.JSON, overwrite bool) (credentials.JSON, error) {
	panic("Not implemented")
}

func (ch CredHub) SetPassword(name string, value values.Password, overwrite bool) (credentials.Password, error) {
	panic("Not implemented")
}

func (ch CredHub) SetUser(name string, value values.User, overwrite bool) (credentials.User, error) {
	panic("Not implemented")
}

func (ch CredHub) SetCertificate(name string, value values.Certificate, overwrite bool) (credentials.Certificate, error) {
	panic("Not implemented")
}

func (ch CredHub) SetRSA(name string, value values.RSA, overwrite bool) (credentials.RSA, error) {
	panic("Not implemented")
}

func (ch CredHub) SetSSH(name string, value values.SSH, overwrite bool) (credentials.SSH, error) {
	panic("Not implemented")
}
