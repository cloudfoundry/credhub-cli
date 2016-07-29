package util

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/models"
)

type RepositoryStub func(req *http.Request, identifier string) (models.Item, error)

func SequentialStub(stubs ...RepositoryStub) RepositoryStub {
	return func(req *http.Request, identifier string) (models.Item, error) {
		var s RepositoryStub
		s, stubs = stubs[0], stubs[1:]
		return s(req, identifier)
	}
}
