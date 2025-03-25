package credhub

import (
	"encoding/json"
	"io"
	"net/http"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
)

type RegenerateOption func(options *RegenerateOptions) error

type regenerateRequest struct {
	Name       string      `json:"name"`
	Regenerate bool        `json:"regenerate"`
	Parameters interface{} `json:"parameters,omitempty"`
	RegenerateOptions
}

type RegenerateOptions struct {
	Metadata credentials.Metadata `json:"metadata,omitempty"`
}

// Regenerate generates and returns a new credential version using the same parameters as the existing credential. The returned credential may be of any type.
func (ch *CredHub) Regenerate(name string, gen interface{}, options ...RegenerateOption) (credentials.Credential, error) {
	var cred credentials.Credential

	request := regenerateRequest{
		Name:       name,
		Regenerate: true,
		Parameters: gen,
	}

	for _, option := range options {
		if err := option(&request.RegenerateOptions); err != nil {
			return cred, err
		}
	}

	serverVersion, err := ch.ServerVersion()
	if err != nil {
		return cred, err
	}

	if request.Metadata != nil && !supportsMetadata(serverVersion) {
		return cred, ServerDoesNotSupportMetadataError
	}

	resp, err := ch.Request(http.MethodPost, "/api/v1/data", nil, request, true)

	if err != nil {
		return credentials.Credential{}, err
	}

	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&cred)

	return cred, err
}
