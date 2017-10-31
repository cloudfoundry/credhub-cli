package credhub

import (
	"errors"
	"net/http"

	"encoding/json"

	"io/ioutil"

	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/permissions"
)

type permissionResponse struct {
	CredentialName string                   `json:"credential_name"`
	Permissions    []permissions.Permission `json:"permissions"`
}

// GetPermissions returns the permissions of a credential.
func (ch *CredHub) GetPermissions(credName string) ([]permissions.Permission, error) {
	panic("Not implemented")
}

// AddPermissions adds permissions to a credential.
func (ch *CredHub) AddPermissions(credName string, perms []permissions.Permission) ([]permissions.Permission, error) {
	requestBody := map[string]interface{}{}
	requestBody["credential_name"] = credName
	requestBody["permissions"] = perms

	response, err := ch.Request(http.MethodPost, "/api/v1/permissions", nil, requestBody)
	if err != nil {
		return nil, err
	}

	responseBodyRaw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot read response body in AddPermissions: %s", err))
	}

	responseBody := new(permissionResponse)
	err = json.Unmarshal(responseBodyRaw, responseBody)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot unmarshal JSON in AddPermissions: %s", err))
	}

	return responseBody.Permissions, nil
}

// DeletePermissions deletes permissions on a credential by actor.
func (ch *CredHub) DeletePermissions(credName string, actor string) error {
	panic("Not implemented")
}
