package repositories

import (
	"net/http"

	"github.com/pivotal-cf/cm-cli/models"
)

type Repository interface {
	SendRequest(request *http.Request, identifier string) (models.Item, error)
}
