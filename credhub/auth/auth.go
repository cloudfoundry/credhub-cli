package auth

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
)

type Auth interface {
	Client() http.Client
}

type AuthOption func(server.Server) Auth
