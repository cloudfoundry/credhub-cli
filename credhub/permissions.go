package credhub

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/permissions"
)

// GetPermissions returns the permissions of a credential.
func (ch *CredHub) GetPermissions(credName string) ([]permissions.Permission, error) {
	panic("Not implemented")
}

// AddPermissions adds permissions to a credential.
func (ch *CredHub) AddPermissions(credName string, perms []permissions.Permission) ([]permissions.Permission, error) {
	requestBody := map[string]interface{}{}
	requestBody["credential_name"] = credName
	requestBody["permissions"] = perms

	_, err := ch.Request(http.MethodPost, "/api/v1/permissions", nil, requestBody)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// DeletePermissions deletes permissions on a credential by actor.
func (ch *CredHub) DeletePermissions(credName string, actor string) error {
	panic("Not implemented")
}
