// CredHub provides methods to interact with CredHub server
package credhub

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
)

type CredHub struct {
	Server server.Server
	Auth   auth.Auth
}

func (ch CredHub) Request(method string, pathStr string, body interface{}) (http.Response, error) {
	panic("Not implemented")
}

func New(server server.Server, authOption auth.AuthOption) CredHub {
	panic("Not implemented")
}
