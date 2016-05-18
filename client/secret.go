package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
)

type secretPutRequest struct {
	Value string `json:"value"`
}

type Secret struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

func (secret *Secret) PrintSecret(){
	fmt.Println(fmt.Sprintf("Name:	%s\nValue:	%s", secret.Name, secret.Value))
}

func NewPutSecretRequest(apiTarget, secretIdentifier, secretContent string) *http.Request {
	url := apiTarget + "/api/v1/secret/" + secretIdentifier

	secret := secretPutRequest{Value: secretContent}
	body, _ := json.Marshal(secret)

	request, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	return request
}

func NewGetSecretRequest(apiTarget, secretIdentifier string) *http.Request {
	url := apiTarget + "/api/v1/secret/" + secretIdentifier

	request, _ := http.NewRequest("GET", url, nil)

	return request
}

func NewDeleteSecretRequest(apiTarget, secretIdentifier string) *http.Request {
	url := apiTarget + "/api/v1/secret/" + secretIdentifier

	request, _ := http.NewRequest("DELETE", url, nil)

	return request
}

func NewInfoRequest(apiTarget string) *http.Request {
	url := apiTarget + "/info"

	request, _ := http.NewRequest("GET", url, nil)

	return request
}
