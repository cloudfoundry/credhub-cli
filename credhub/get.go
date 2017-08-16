package credhub

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
)

// Returns the current credential by ID.
func (ch *CredHub) GetById(id string) (credentials.Credential, error) {
	panic("Not implemented")
}

// Returns all historical credential value(s) by name.
func (ch *CredHub) GetAll(name string) ([]credentials.Credential, error) {
	panic("Not implemented")
}

// Returns the current credential by name.
func (ch *CredHub) Get(name string) (credentials.Credential, error) {
	var cred credentials.Credential
	err := ch.getCurrentCredential(name, &cred)
	return cred, err
}

// Returns the Value credential by name.
func (ch *CredHub) GetValue(name string) (credentials.Value, error) {
	panic("Not implemented")
}

// Returns the JSON credential by name.
func (ch *CredHub) GetJSON(name string) (credentials.JSON, error) {
	panic("Not implemented")
}

// Returns the current Password credential by name.
func (ch *CredHub) GetPassword(name string) (credentials.Password, error) {
	var cred credentials.Password
	err := ch.getCurrentCredential(name, &cred)

	return cred, err
}

// Returns the current User credential by name.
func (ch *CredHub) GetUser(name string) (credentials.User, error) {
	var cred credentials.User
	err := ch.getCurrentCredential(name, &cred)

	return cred, err
}

// Returns the current Certificate credential by name.
func (ch *CredHub) GetCertificate(name string) (credentials.Certificate, error) {
	var cred credentials.Certificate
	err := ch.getCurrentCredential(name, &cred)

	return cred, err
}

// Returns the current RSA credential by name.
func (ch *CredHub) GetRSA(name string) (credentials.RSA, error) {
	panic("Not implemented")
}

// Returns the current SSH credential by name.
func (ch *CredHub) GetSSH(name string) (credentials.SSH, error) {
	panic("Not implemented")
}

func (ch *CredHub) getCurrentCredential(name string, cred interface{}) error {
	query := url.Values{}
	query.Set("name", name)
	query.Set("current", "true")

	resp, err := ch.Request(http.MethodGet, "/api/v1/data", query, nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	response := make(map[string][]json.RawMessage)

	if err := dec.Decode(&response); err != nil {
		return err
	}

	var ok bool
	var data []json.RawMessage

	if data, ok = response["data"]; !ok || len(data) == 0 {
		return errors.New("response did not contain any credentials")
	}

	rawMessage := data[0]

	return json.Unmarshal(rawMessage, cred)
}
