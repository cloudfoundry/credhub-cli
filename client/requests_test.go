package client_test

import (
	"net/http"

	. "github.com/pivotal-cf/credhub-cli/client"

	"bytes"

	"fmt"

	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/models"
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

	Describe("NewTokenRevocationRequest", func() {
		It("Returns a request to revoke a refresh token", func() {
			cfg.RefreshToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImxlZ2FjeS10b2tlbi1rZXkiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiI1YjljOWZkNTFiYTE0ODM4YWMyZTZiMjIyZDQ4NzEwNi1yIiwic3ViIjoiYzE0ZGJjZGQtNzNkOC00ZDdjLWI5NDctYzM4ODVhODAxYzY2Iiwic2NvcGUiOlsiY3JlZGh1Yi53cml0ZSIsImNyZWRodWIucmVhZCJdLCJpYXQiOjE0NzEzMTAwMTIsImV4cCI6MTQ3MTM5NjQxMiwiY2lkIjoiY3JlZGh1YiIsImNsaWVudF9pZCI6ImNyZWRodWIiLCJpc3MiOiJodHRwczovLzUyLjIwNC40OS4xMDc6ODQ0My9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsInJldm9jYWJsZSI6dHJ1ZSwiZ3JhbnRfdHlwZSI6InBhc3N3b3JkIiwidXNlcl9uYW1lIjoiY3JlZGh1Yl9jbGkiLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX2lkIjoiYzE0ZGJjZGQtNzNkOC00ZDdjLWI5NDctYzM4ODVhODAxYzY2IiwicmV2X3NpZyI6ImQ3MTkyZmUxIiwiYXVkIjpbImNyZWRodWIiXX0.UAp6Ou24f18mdE0XOqG9RLVWZAx3khNHHPeHfuzmcOUYojtILa0_izlGVHhCtNx07f4M9pcRKpo-AijXRw1vSimSTHBeVCDjuuc2nBdznIMhyQSlPpd2stW-WG7Gix82K4gy4oCb1wlTqsK3UKGYoy8JWs6XZqhoZZ6JZM7-Xjj2zag3Q4kgvEBReWC5an_IP6SeCpNt5xWvGdxtTz7ki1WPweUBy0M73ZjRi9_poQT2JmeSIbrePukkfsfCxHG1vM7ApIdzzhdCx6T_KmmMU3xHqhpI_ueLOuvfHjdBinm2atypeTHD83yRRFxhfjRsG1-XguTn-lo_Z2Jis89r5g"
			cfg.AccessToken = "defgh"
			expectedRequest, _ := http.NewRequest(
				"DELETE",
				cfg.AuthURL+"/oauth/token/revoke/5b9c9fd51ba14838ac2e6b222d487106-r",
				nil)
			expectedRequest.Header.Add("Authorization", "Bearer defgh")

			request, _ := NewTokenRevocationRequest(cfg)

			Expect(request).To(Equal(expectedRequest))
		})

		It("Returns error when the token cannot be parsed", func() {
			cfg.RefreshToken = "fake_refresh"
			cfg.AccessToken = "defgh"

			_, err := NewTokenRevocationRequest(cfg)

			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Api Server Requests", func() {
		BeforeEach(func() {
			cfg.AccessToken = "access-token"
		})

		Describe("NewPutValueRequest", func() {
			It("Returns a request for the put-value endpoint", func() {
				requestBody := bytes.NewReader([]byte(`{"type":"value","value":"my-value","parameters":{"overwrite":true}}`))
				expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/data/my-name", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewPutValueRequest(cfg, "my-name", "my-value", true)

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewPutPasswordRequest", func() {
			It("Returns a request for the put-password endpoint", func() {
				requestBody := bytes.NewReader([]byte(`{"type":"password","value":"my-password","parameters":{"overwrite":true}}`))
				expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/data/my-name", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewPutPasswordRequest(cfg, "my-name", "my-password", true)

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewPutCertificateRequest", func() {
			It("Returns a request for the put-certificate endpoint", func() {
				json := fmt.Sprintf(`{"type":"certificate","value":{"root":"%s","certificate":"%s","private_key":"%s"},"parameters":{"overwrite":true}}`,
					"my-ca", "my-cert", "my-priv")
				requestBody := bytes.NewReader([]byte(json))
				expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/data/my-name", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewPutCertificateRequest(cfg, "my-name", "my-ca", "my-cert", "my-priv", true)

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewPutCaRequest", func() {
			It("Returns a request for the put-root-ca endpoint", func() {
				json := fmt.Sprintf(`{"type":"root","value":{"certificate":"%s","private_key":"%s"}}`,
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
					"overwrite": false,
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
			It("returns a request with only overwrite", func() {
				requestBody := bytes.NewReader([]byte(`{"type":"my-type","parameters":{"overwrite":false}}`))
				expectedRequest, _ := http.NewRequest("POST", "http://example.com/api/v1/data/my-name", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				params := models.SecretParameters{}
				request := NewGenerateSecretRequest(cfg, "my-name", params, "my-type")

				Expect(request).To(Equal(expectedRequest))
			})

			It("returns a request with parameters", func() {
				parameters := models.SecretParameters{
					Overwrite:      false,
					ExcludeSpecial: true,
					ExcludeNumber:  true,
					ExcludeUpper:   true,
					ExcludeLower:   true,
					Length:         42,
				}
				expectedRequestBody := `{
				"type":"value",
				"parameters": {
					"overwrite": false,
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
