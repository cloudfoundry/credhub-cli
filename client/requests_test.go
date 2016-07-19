package client_test

import (
	"net/http"

	. "github.com/pivotal-cf/cm-cli/client"

	"bytes"

	"fmt"

	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/cm-cli/config"
	"github.com/pivotal-cf/cm-cli/models"
)

var _ = Describe("API", func() {
	Describe("NewInfoRequest", func() {
		It("Returns a request for the info endpoint", func() {
			expectedRequest, _ := http.NewRequest("GET", "fake_target.com/info", nil)

			request := NewInfoRequest("fake_target.com")

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewAuthTokenRequest", func() {
		It("Returns a request for the uaa oauth token endpoint", func() {
			config := config.Config{
				AuthURL: "http://example.com/uaa",
			}
			user := "my-user"
			pass := "my-pass"
			data := url.Values{}
			data.Set("grant_type", "password")
			data.Add("response_type", "token")
			data.Add("username", user)
			data.Add("password", pass)
			expectedRequest, _ := http.NewRequest(
				"POST",
				config.AuthURL+"/oauth/token/",
				bytes.NewBufferString(data.Encode()))
			expectedRequest.SetBasicAuth("credhub", "")
			expectedRequest.Header.Add("Accept", "application/json")
			expectedRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			request := NewAuthTokenRequest(config, user, pass)

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewPutSecretValueRequest", func() {
		It("Returns a request for the put-secret endpoint", func() {
			requestBody := bytes.NewReader([]byte(`{"type":"value","credential":"my-value"}`))
			expectedRequest, _ := http.NewRequest("PUT", "sample.com/api/v1/data/my-name", requestBody)
			expectedRequest.Header.Set("Content-Type", "application/json")

			request := NewPutValueRequest("sample.com", "my-name", "my-value")

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewPutCertificateRequest", func() {
		It("Returns a request for the put-certificate endpoint", func() {
			json := fmt.Sprintf(`{"type":"certificate","credential":{"root":"%s","certificate":"%s","private":"%s"}}`,
				"my-ca", "my-cert", "my-priv")
			requestBody := bytes.NewReader([]byte(json))
			expectedRequest, _ := http.NewRequest("PUT", "sample.com/api/v1/data/my-name", requestBody)
			expectedRequest.Header.Set("Content-Type", "application/json")

			request := NewPutCertificateRequest("sample.com", "my-name", "my-ca", "my-cert", "my-priv")

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewPutCaRequest", func() {
		It("Returns a request for the put-root-ca endpoint", func() {
			json := fmt.Sprintf(`{"type":"root","ca":{"certificate":"%s","private":"%s"}}`,
				"my-cert", "my-priv")
			requestBody := bytes.NewReader([]byte(json))
			expectedRequest, _ := http.NewRequest("PUT", "sample.com/api/v1/ca/my-name", requestBody)
			expectedRequest.Header.Set("Content-Type", "application/json")

			request := NewPutCaRequest("sample.com", "my-name", "root", "my-cert", "my-priv")

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewPostCaRequest", func() {
		It("Returns a request for the post-root-ca endpoint", func() {
			parameters := models.SecretParameters{
				CommonName:       "my-common-name",
				Organization:     "my-organization",
				OrganizationUnit: "my-unit",
				Locality:         "my-locality",
				State:            "my-state",
				Country:          "my-country",
			}
			expectedRequestJson := `{
				"type":"root",
				"parameters": {
					"common_name": "my-common-name",
					"organization": "my-organization",
					"organization_unit": "my-unit",
					"locality": "my-locality",
					"state": "my-state",
					"country": "my-country"
				}
			}`

			expectedRequestBody := bytes.NewReader([]byte(expectedRequestJson))

			expectedRequest, _ := http.NewRequest("POST", "sample.com/api/v1/ca/my-name", expectedRequestBody)
			expectedRequest.Header.Set("Content-Type", "application/json")

			request := NewPostCaRequest("sample.com", "my-name", "root", parameters)

			bodyBuffer := new(bytes.Buffer)
			bodyBuffer.ReadFrom(request.Body)
			Expect(bodyBuffer).To(MatchJSON(expectedRequestJson))
			Expect(request.Method).To(Equal("POST"))
			Expect(request.URL.String()).To(Equal("sample.com/api/v1/ca/my-name"))
			Expect(request.Header.Get("Content-Type")).To(Equal("application/json"))
		})
	})

	Describe("NewGetCaRequest", func() {
		It("Returns a request for the get-root-ca endpoint", func() {
			expectedRequest, _ := http.NewRequest("GET", "sample.com/api/v1/ca/my-name", nil)

			request := NewGetCaRequest("sample.com", "my-name")

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewGenerateSecretRequest", func() {
		It("returns a request with no parameters", func() {
			requestBody := bytes.NewReader([]byte(`{"type":"my-type","parameters":{}}`))
			expectedRequest, _ := http.NewRequest("POST", "sample.com/api/v1/data/my-name", requestBody)
			expectedRequest.Header.Set("Content-Type", "application/json")

			request := NewGenerateSecretRequest("sample.com", "my-name", models.SecretParameters{}, "my-type")

			Expect(request).To(Equal(expectedRequest))
		})

		It("returns a request with parameters", func() {
			parameters := models.SecretParameters{
				ExcludeSpecial: true,
				ExcludeNumber:  true,
				ExcludeUpper:   true,
				ExcludeLower:   true,
				Length:         42,
			}
			expectedRequestBody := `{
				"type":"value",
				"parameters": {
					"exclude_special": true,
					"exclude_number": true,
					"exclude_upper": true,
					"exclude_lower": true,
					"length": 42
				}
			}`

			request := NewGenerateSecretRequest("sample.com", "my-name", parameters, "value")

			bodyBuffer := new(bytes.Buffer)
			bodyBuffer.ReadFrom(request.Body)
			Expect(bodyBuffer).To(MatchJSON(expectedRequestBody))
			Expect(request.Method).To(Equal("POST"))
			Expect(request.URL.String()).To(Equal("sample.com/api/v1/data/my-name"))
			Expect(request.Header.Get("Content-Type")).To(Equal("application/json"))
		})
	})

	Describe("NewGetSecretRequest", func() {
		It("Returns a request for getting secret", func() {
			expectedRequest, _ := http.NewRequest("GET", "sample.com/api/v1/data/my-name", nil)

			request := NewGetSecretRequest("sample.com", "my-name")

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewDeleteSecretRequest", func() {
		It("Returns a request for deleting", func() {
			expectedRequest, _ := http.NewRequest("DELETE", "sample.com/api/v1/data/my-name", nil)

			request := NewDeleteSecretRequest("sample.com", "my-name")

			Expect(request).To(Equal(expectedRequest))
		})
	})
})
