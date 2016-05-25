package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SecretBody struct {
	Value string `json:"value"`
}

type Secret struct {
	Name       string
	SecretBody SecretBody
}

func NewSecret(name string, secretBody SecretBody) Secret {
	return Secret{
		Name:       name,
		SecretBody: secretBody,
	}
}

func (secret *Secret) PrintSecret() {
	fmt.Println(fmt.Sprintf("Name:	%s\nValue:	%s", secret.Name, secret.SecretBody.Value))
}

func NewPutSecretRequest(apiTarget, secretIdentifier, secretContent string) *http.Request {
	url := apiTarget + "/api/v1/data/" + secretIdentifier

	secret := SecretBody{Value: secretContent}
	body, _ := json.Marshal(secret)

	request, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	return request
}

func NewGenerateSecretRequest(apiTarget, secretIdentifier string) *http.Request {
	url := apiTarget + "/api/v1/data/" + secretIdentifier

	secret := new(struct{})
	body, _ := json.Marshal(secret)

	request, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	return request
}

func NewGetSecretRequest(apiTarget, secretIdentifier string) *http.Request {
	url := apiTarget + "/api/v1/data/" + secretIdentifier

	request, _ := http.NewRequest("GET", url, nil)

	return request
}

func NewDeleteSecretRequest(apiTarget, secretIdentifier string) *http.Request {
	url := apiTarget + "/api/v1/data/" + secretIdentifier

	request, _ := http.NewRequest("DELETE", url, nil)

	return request
}

func NewInfoRequest(apiTarget string) *http.Request {
	url := apiTarget + "/info"

	request, _ := http.NewRequest("GET", url, nil)

	return request
}
