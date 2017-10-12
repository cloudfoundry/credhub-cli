package credhub

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
)

// GeneratePassword generates a password credential based on the provided parameters.
func (ch *CredHub) GeneratePassword(name string, gen generate.Password, overwrite mode) (credentials.Password, error) {
	var cred credentials.Password
	err := ch.generateCredential(name, "password", gen, overwrite, &cred)
	return cred, err
}

// GenerateUser generates a user credential based on the provided parameters.
func (ch *CredHub) GenerateUser(name string, gen generate.User, overwrite mode) (credentials.User, error) {
	var cred credentials.User
	err := ch.generateCredential(name, "user", gen, overwrite, &cred)
	return cred, err
}

// GenerateCertificate generates a certificate credential based on the provided parameters.
func (ch *CredHub) GenerateCertificate(name string, gen generate.Certificate, overwrite mode) (credentials.Certificate, error) {
	var cred credentials.Certificate
	err := ch.generateCredential(name, "certificate", gen, overwrite, &cred)
	return cred, err
}

// GenerateRSA generates an RSA credential based on the provided parameters.
func (ch *CredHub) GenerateRSA(name string, gen generate.RSA, overwrite mode) (credentials.RSA, error) {
	var cred credentials.RSA
	err := ch.generateCredential(name, "rsa", gen, overwrite, &cred)
	return cred, err
}

// GenerateSSH generates an SSH credential based on the provided parameters.
func (ch *CredHub) GenerateSSH(name string, gen generate.SSH, overwrite mode) (credentials.SSH, error) {
	var cred credentials.SSH
	err := ch.generateCredential(name, "ssh", gen, overwrite, &cred)
	return cred, err
}

// GenerateCredential generates any credential type based on the credType given provided parameters.
func (ch *CredHub) GenerateCredential(name, credType string, gen interface{}, overwrite mode) (credentials.Credential, error) {
	var cred credentials.Credential
	err := ch.generateCredential(name, credType, gen, overwrite, &cred)
	return cred, err
}

func (ch *CredHub) generateCredential(name, credType string, gen interface{}, overwrite mode, cred interface{}) error {
	isOverwrite := true
	if overwrite == NoOverwrite {
		isOverwrite = false
	}
	requestBody := map[string]interface{}{}
	requestBody["name"] = name
	requestBody["type"] = credType
	requestBody["parameters"] = gen
	requestBody["overwrite"] = isOverwrite

	if user, ok := gen.(generate.User); ok {
		requestBody["value"] = map[string]string{"username": user.Username}
	}

	resp, err := ch.Request(http.MethodPost, "/api/v1/data", nil, requestBody)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	if err := ch.checkForServerError(resp); err != nil {
		return err
	}

	return dec.Decode(&cred)
}
