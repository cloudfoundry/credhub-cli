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
	var cfg config.Config

	BeforeEach(func() {
		cfg = config.Config{
			ApiURL:  "http://example.com",
			AuthURL: "http://example.com/uaa",
		}
	})

	Describe("NewInfoRequest", func() {
		It("Returns a request for the info endpoint", func() {
			expectedRequest, _ := http.NewRequest("GET", "http://example.com/info", nil)

			request := NewInfoRequest(cfg)

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewAuthTokenRequest", func() {
		It("Returns a request for the uaa oauth token endpoint", func() {
			user := "my-user"
			pass := "my-pass"
			data := url.Values{}
			data.Set("grant_type", "password")
			data.Add("response_type", "token")
			data.Add("username", user)
			data.Add("password", pass)
			expectedRequest, _ := http.NewRequest(
				"POST",
				cfg.AuthURL+"/oauth/token/",
				bytes.NewBufferString(data.Encode()))
			expectedRequest.SetBasicAuth("credhub", "")
			expectedRequest.Header.Add("Accept", "application/json")
			expectedRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			request := NewAuthTokenRequest(cfg, user, pass)

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewRefreshTokenRequest", func() {
		It("Returns a request for the uaa oauth token endpoint to get refresh token", func() {
			data := url.Values{}
			data.Set("grant_type", "refresh_token")
			data.Set("refresh_token", cfg.RefreshToken)
			expectedRequest, _ := http.NewRequest(
				"POST",
				cfg.AuthURL+"/oauth/token/",
				bytes.NewBufferString(data.Encode()))
			expectedRequest.SetBasicAuth("credhub", "")
			expectedRequest.Header.Add("Accept", "application/json")
			expectedRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			request := NewRefreshTokenRequest(cfg)

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("Api Server Requests", func() {
		BeforeEach(func() {
			cfg.AccessToken = "access-token"
		})

		Describe("NewPutSecretValueRequest", func() {
			It("Returns a request for the put-secret endpoint", func() {
				requestBody := bytes.NewReader([]byte(`{"type":"value","credential":"my-value"}`))
				expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/data/my-name", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewPutValueRequest(cfg, "my-name", "my-value")

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewPutCertificateRequest", func() {
			It("Returns a request for the put-certificate endpoint", func() {
				json := fmt.Sprintf(`{"type":"certificate","credential":{"root":"%s","certificate":"%s","private":"%s"}}`,
					"my-ca", "my-cert", "my-priv")
				requestBody := bytes.NewReader([]byte(json))
				expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/data/my-name", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewPutCertificateRequest(cfg, "my-name", "my-ca", "my-cert", "my-priv")

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewPutCaRequest", func() {
			It("Returns a request for the put-root-ca endpoint", func() {
				json := fmt.Sprintf(`{"type":"root","ca":{"certificate":"%s","private":"%s"}}`,
					"my-cert", "my-priv")
				requestBody := bytes.NewReader([]byte(json))
				expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/ca/my-name", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewPutCaRequest(cfg, "my-name", "root", "my-cert", "my-priv")

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

				expectedRequest, _ := http.NewRequest("POST", "http://example.com/api/v1/ca/my-name", expectedRequestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewPostCaRequest(cfg, "my-name", "root", parameters)

				bodyBuffer := new(bytes.Buffer)
				bodyBuffer.ReadFrom(request.Body)
				Expect(bodyBuffer).To(MatchJSON(expectedRequestJson))
				Expect(request.Method).To(Equal("POST"))
				Expect(request.URL.String()).To(Equal("http://example.com/api/v1/ca/my-name"))
				Expect(request.Header.Get("Content-Type")).To(Equal("application/json"))
			})
		})

		Describe("NewGetCaRequest", func() {
			It("Returns a request for the get-root-ca endpoint", func() {
				expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/ca/my-name", nil)
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewGetCaRequest(cfg, "my-name")

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewGenerateSecretRequest", func() {
			It("returns a request with no parameters", func() {
				requestBody := bytes.NewReader([]byte(`{"type":"my-type","parameters":{}}`))
				expectedRequest, _ := http.NewRequest("POST", "http://example.com/api/v1/data/my-name", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewGenerateSecretRequest(cfg, "my-name", models.SecretParameters{}, "my-type")

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

				request := NewGenerateSecretRequest(cfg, "my-name", parameters, "value")

				bodyBuffer := new(bytes.Buffer)
				bodyBuffer.ReadFrom(request.Body)
				Expect(bodyBuffer).To(MatchJSON(expectedRequestBody))
				Expect(request.Method).To(Equal("POST"))
				Expect(request.URL.String()).To(Equal("http://example.com/api/v1/data/my-name"))
				Expect(request.Header.Get("Content-Type")).To(Equal("application/json"))
			})
		})

		Describe("NewGetSecretRequest", func() {
			It("Returns a request for getting secret", func() {
				expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data/my-name", nil)
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewGetSecretRequest(cfg, "my-name")

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewDeleteSecretRequest", func() {
			It("Returns a request for deleting", func() {
				expectedRequest, _ := http.NewRequest("DELETE", "http://example.com/api/v1/data/my-name", nil)
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewDeleteSecretRequest(cfg, "my-name")

				Expect(request).To(Equal(expectedRequest))
			})
		})
	})
})
