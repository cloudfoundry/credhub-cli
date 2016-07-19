package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"net/url"

	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/models"
)

func NewPutValueRequest(apiTarget, secretIdentifier, secretContent string) *http.Request {
	secret := models.SecretBody{
		Credential:  secretContent,
		ContentType: "value",
	}

	return newSecretRequest("PUT", apiTarget, secretIdentifier, secret)
}

func NewPutCertificateRequest(apiTarget, secretIdentifier, root string, cert string, priv string) *http.Request {
	certificate := models.Certificate{
		Root:        root,
		Certificate: cert,
		Private:     priv,
	}
	secret := models.SecretBody{
		ContentType: "certificate",
		Credential:  &certificate,
	}

	return newSecretRequest("PUT", apiTarget, secretIdentifier, secret)
}

func NewPutCaRequest(apiTarget, caIdentifier, caType, cert, priv string) *http.Request {
	ca := models.CaParameters{
		Certificate: cert,
		Private:     priv,
	}
	caBody := models.CaBody{
		ContentType: caType,
		Ca:          &ca,
	}

	return newCaRequest("PUT", apiTarget, caIdentifier, caBody)
}

func NewPostCaRequest(apiTarget, caIdentifier, caType string, params models.SecretParameters) *http.Request {
	caGenerateRequestBody := models.GenerateRequest{
		ContentType: caType,
		Parameters:  params,
	}

	return newCaRequest("POST", apiTarget, caIdentifier, caGenerateRequestBody)
}

func NewGetCaRequest(apiTarget, caIdentifier string) *http.Request {
	return newCaRequest("GET", apiTarget, caIdentifier, nil)
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

func NewAuthTokenRequest(cfg config.Config, user string, pass string) *http.Request {
	authUrl := cfg.AuthURL + "/oauth/token/"
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Add("response_type", "token")
	data.Add("username", user)
	data.Add("password", pass)
	request, _ := http.NewRequest("POST", authUrl, bytes.NewBufferString(data.Encode()))
	request.SetBasicAuth(config.AuthClient, config.AuthPassword)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func newSecretRequest(requestType, apiTarget, secretIdentifier string, bodyModel interface{}) *http.Request {
	url := apiTarget + "/api/v1/data/" + secretIdentifier

	return newRequest(requestType, url, bodyModel)
}

func newCaRequest(requestType, apiTarget, caIdentifier string, bodyModel interface{}) *http.Request {
	url := apiTarget + "/api/v1/ca/" + caIdentifier

	return newRequest(requestType, url, bodyModel)
}

func newRequest(requestType, url string, bodyModel interface{}) *http.Request {
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
