package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type secretRequest struct {
	Value string `json:"value"`
}

func NewPutSecretRequest(apiTarget, secretIdentifier, secretContent string) *http.Request {
	url := apiTarget + "/api/v1/secret/" + secretIdentifier

	secret := secretRequest{Value: secretContent}
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
