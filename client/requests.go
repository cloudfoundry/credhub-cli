package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf/cm-cli/models"
)

func NewPutValueRequest(apiTarget, secretIdentifier, secretContent string) *http.Request {
	secret := models.SecretBody{
		Value:       secretContent,
		ContentType: "value",
	}

	return newSecretRequest("PUT", apiTarget, secretIdentifier, secret)
}

func NewPutCertificateRequest(apiTarget, secretIdentifier, ca string, pub string, priv string) *http.Request {
	certificate := models.Certificate{
		Ca:      ca,
		Public:  pub,
		Private: priv,
	}
	secret := models.SecretBody{
		ContentType: "certificate",
		Certificate: &certificate,
	}

	return newSecretRequest("PUT", apiTarget, secretIdentifier, secret)
}

func NewGenerateSecretRequest(apiTarget, secretIdentifier string, parameters models.SecretParameters, contentType string) *http.Request {
	generateRequest := models.GenerateRequest{
		Parameters:  parameters,
		ContentType: contentType,
	}

	return newSecretRequest("POST", apiTarget, secretIdentifier, generateRequest)
}

func NewGetSecretRequest(apiTarget, secretIdentifier string) *http.Request {
	return newSecretRequest("GET", apiTarget, secretIdentifier, nil)
}

func NewDeleteSecretRequest(apiTarget, secretIdentifier string) *http.Request {
	return newSecretRequest("DELETE", apiTarget, secretIdentifier, nil)
}

func NewInfoRequest(apiTarget string) *http.Request {
	url := apiTarget + "/info"

	request, _ := http.NewRequest("GET", url, nil)

	return request
}

func newSecretRequest(requestType, apiTarget, secretIdentifier string, bodyModel interface{}) *http.Request {
	url := apiTarget + "/api/v1/data/" + secretIdentifier

	var request *http.Request
	if bodyModel == nil {
		request, _ = http.NewRequest(requestType, url, nil)
	} else {
		body, _ := json.Marshal(bodyModel)
		request, _ = http.NewRequest(requestType, url, bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
	}

	return request
}
