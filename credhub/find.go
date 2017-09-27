package credhub

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
)

// FindByPartialName retrieves a list of stored credential names which contain the search.
func (ch *CredHub) FindByPartialName(nameLike string) (credentials.FindByNameResults, error) {
	return ch.findByPathOrNameLike("name-like", nameLike)
}

// FindByPath retrieves a list of stored credential names which are within the specified path.
func (ch *CredHub) FindByPath(path string) (credentials.FindByNameResults, error) {
	return ch.findByPathOrNameLike("path", path)
}

// ShowAllPaths retrieves a list of all paths which contain credentials.
func (ch *CredHub) ShowAllPaths() ([]credentials.Path, error) {
	panic("Not implemented")
}

func (ch *CredHub) findByPathOrNameLike(key, value string) (credentials.FindByNameResults, error) {
	var creds credentials.FindByNameResults

	query := url.Values{}
	query.Set(key, value)

	resp, err := ch.Request(http.MethodGet, "/api/v1/data", query, nil)

	if err != nil {
		return credentials.FindByNameResults{}, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &creds)

	if err != nil {
		return credentials.FindByNameResults{}, err
	}

	return creds, nil
}
