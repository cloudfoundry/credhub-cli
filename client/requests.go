package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/models"
)

func NewPutValueRequest(apiTarget, secretIdentifier, secretContent string) *http.Request {
	url := apiTarget + "/api/v1/data/" + secretIdentifier

	secret := models.SecretBody{
		Value:       secretContent,
		ContentType: "value",
	}
	body, _ := json.Marshal(secret)

	request, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	return request
}

func NewPutCertificateRequest(apiTarget, secretIdentifier, ca string, pub string, priv string) *http.Request {
	url := apiTarget + "/api/v1/data/" + secretIdentifier

	secret := models.CertificateRequest{
		ContentType: "certificate",
		Certificate: models.Certificate{
			Ca:      ca,
			Public:  pub,
			Private: priv,
		},
	}
	body, _ := json.Marshal(secret)

	request, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	return request
}

func NewGenerateSecretRequest(apiTarget, secretIdentifier string, parameters models.SecretParameters, contentType string) *http.Request {
	url := apiTarget + "/api/v1/data/" + secretIdentifier

	generateRequest := models.GenerateRequest{
		Parameters:  parameters,
		ContentType: contentType,
	}

	body, _ := json.Marshal(generateRequest)

	request, _ := http.NewRequest("POST", url, bytes.NewReader(body))
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
