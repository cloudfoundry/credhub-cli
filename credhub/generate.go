package credhub

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
)

// Generates a password credential based on the provided parameters.
func (ch *CredHub) GeneratePassword(name string, gen generate.Password, overwrite bool) (credentials.Password, error) {
	panic("Not implemented")
}

// Generates a user credential based on the provided parameters.
func (ch *CredHub) GenerateUser(name string, gen generate.User, overwrite bool) (credentials.User, error) {
	panic("Not implemented")
}

// Generates a user credential based on the provided parameters.
func (ch *CredHub) GenerateCertificate(name string, gen generate.Certificate, overwrite bool) (credentials.Certificate, error) {
	var cred credentials.Certificate

	requestBody := map[string]interface{}{}
	requestBody["name"] = name
	requestBody["type"] = "certificate"
	requestBody["parameters"] = gen
	resp, err := ch.Request(http.MethodPost, "/api/v1/data", requestBody)

	if err != nil {
		return cred, err
	}
	var responseBody map[string]([]credentials.Certificate)

	defer resp.Body.Close()
	bodyResp, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(bodyResp, &responseBody)
	if err != nil {
		return cred, err
	}

	if len(responseBody["data"]) < 1 {
		return cred, errors.New("Could not generate credential")
	}

	cred = responseBody["data"][0]
	return cred, nil
}

// Generates an RSA credential based on the provided parameters.
func (ch *CredHub) GenerateRSA(name string, gen generate.RSA, overwrite bool) (credentials.RSA, error) {
	panic("Not implemented")
}

// Generates an SSH credential based on the provided parameters.
func (ch *CredHub) GenerateSSH(name string, gen generate.SSH, overwrite bool) (credentials.SSH, error) {
	panic("Not implemented")
}
