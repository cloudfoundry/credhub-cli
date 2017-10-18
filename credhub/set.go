package credhub

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/hashicorp/go-version"
)

// SetValue sets a value credential with a user-provided value.
func (ch *CredHub) SetValue(name string, value values.Value, overwrite Mode) (credentials.Value, error) {
	var cred credentials.Value
	err := ch.setCredential(name, "value", value, overwrite, &cred)

	return cred, err
}

// SetJSON sets a JSON credential with a user-provided value.
func (ch *CredHub) SetJSON(name string, value values.JSON, overwrite Mode) (credentials.JSON, error) {
	var cred credentials.JSON
	err := ch.setCredential(name, "json", value, overwrite, &cred)

	return cred, err
}

// SetPassword sets a password credential with a user-provided value.
func (ch *CredHub) SetPassword(name string, value values.Password, overwrite Mode) (credentials.Password, error) {
	var cred credentials.Password
	err := ch.setCredential(name, "password", value, overwrite, &cred)

	return cred, err
}

// SetUser sets a user credential with a user-provided value.
func (ch *CredHub) SetUser(name string, value values.User, overwrite Mode) (credentials.User, error) {
	var cred credentials.User
	err := ch.setCredential(name, "user", value, overwrite, &cred)

	return cred, err
}

// SetCertificate sets a certificate credential with a user-provided value.
func (ch *CredHub) SetCertificate(name string, value values.Certificate, overwrite Mode) (credentials.Certificate, error) {
	var cred credentials.Certificate
	err := ch.setCredential(name, "certificate", value, overwrite, &cred)

	return cred, err
}

// SetRSA sets an RSA credential with a user-provided value.
func (ch *CredHub) SetRSA(name string, value values.RSA, overwrite Mode) (credentials.RSA, error) {
	var cred credentials.RSA
	err := ch.setCredential(name, "rsa", value, overwrite, &cred)

	return cred, err
}

// SetSSH sets an SSH credential with a user-provided value.
func (ch *CredHub) SetSSH(name string, value values.SSH, overwrite Mode) (credentials.SSH, error) {
	var cred credentials.SSH
	err := ch.setCredential(name, "ssh", value, overwrite, &cred)

	return cred, err
}

// SetCredential sets a credential of any type with a user-provided value. 
func (ch *CredHub) SetCredential(name, credType string, value interface{}, overwrite Mode) (credentials.Credential, error) {
	var cred credentials.Credential
	err := ch.setCredential(name, credType, value, overwrite, &cred)

	return cred, err
}

func (ch *CredHub) setCredential(name, credType string, value interface{}, overwrite Mode, cred interface{}) error {
	isOverwrite := overwrite == Overwrite

	requestBody := map[string]interface{}{}
	requestBody["name"] = name
	requestBody["type"] = credType
	requestBody["value"] = value

	serverVersion, err := ch.ServerVersion()
	if err != nil {
		return err
	}

	constraints, err := version.NewConstraint("< 1.6.0")
	if constraints.Check(serverVersion) {
		if overwrite == Converge {
			return fmt.Errorf("Interaction Mode 'converge' not supported on target server (version: <%s>)", serverVersion.String())
		}
		requestBody["overwrite"] = isOverwrite
	} else {
		requestBody["mode"] = overwrite
	}
	resp, err := ch.Request(http.MethodPut, "/api/v1/data", nil, requestBody)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(cred)
}
