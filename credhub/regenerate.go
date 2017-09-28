package credhub

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	version "github.com/hashicorp/go-version"
)

// Regenerate generates and returns a new credential version using the same parameters existing credential. The returned credential may be of any type.
func (ch *CredHub) Regenerate(name string) (credentials.Credential, error) {
	var cred credentials.Credential

	regenerateEndpoint := "/api/v1/regenerate"

	requestBody := map[string]interface{}{}
	requestBody["name"] = name

	if ch.ServerVersion != "" {
		serverVersion, err := version.NewVersion(ch.ServerVersion)
		if err != nil {
			return credentials.Credential{}, err
		}

		constraints, err := version.NewConstraint("< 1.4.0")
		if constraints.Check(serverVersion) {
			regenerateEndpoint = "/api/v1/data"
			requestBody["regenerate"] = true
		}
	}

	resp, err := ch.Request(http.MethodPost, regenerateEndpoint, nil, requestBody)

	if err != nil {
		return credentials.Credential{}, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&cred)

	return cred, err
}
