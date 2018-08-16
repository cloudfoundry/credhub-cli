package credhub

import (
	"io"
	"io/ioutil"
	"net/http"

	"encoding/json"

	"code.cloudfoundry.org/credhub-cli/credhub/permissions"
)

type permissionsResponse struct {
	CredentialName string                   `json:"credential_name"`
	Permissions    []permissions.Permission `json:"permissions"`
}

func (ch *CredHub) GetPermission(uuid string) (*permissions.Permission, error) {
	path:= "/api/v2/permissions/" + uuid

	resp, err := ch.Request(http.MethodGet, path, nil, nil, true)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)
	var response permissions.Permission

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (ch *CredHub) AddPermission(path string, actor string, ops []string) (*permissions.Permission, error) {
	requestBody := map[string]interface{}{}
	requestBody["path"] = path
	requestBody["actor"] = actor
	requestBody["operations"] = ops

	resp, err := ch.Request(http.MethodPost, "/api/v2/permissions", nil, requestBody, true)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)
	var response permissions.Permission

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
