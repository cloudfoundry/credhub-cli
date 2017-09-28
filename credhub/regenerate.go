package credhub

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
)

// Regenerate generates and returns a new credential version using the same parameters existing credential. The returned credential may be of any type.
func (ch *CredHub) Regenerate(name string) (credentials.Credential, error) {
	var cred credentials.Credential

	requestBody := map[string]interface{}{}
	requestBody["name"] = name
	resp, err := ch.Request(http.MethodPost, "/api/v1/regenerate", nil, requestBody)

	if err != nil {
		return credentials.Credential{}, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&cred)

	return cred, err
}
