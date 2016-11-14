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
			cfg.RefreshToken = "5b9c9fd51ba14838ac2e6b222d487106-r"
			cfg.AccessToken = "defgh"
			expectedRequest, _ := http.NewRequest(
				"DELETE",
				cfg.AuthURL+"/oauth/token/revoke/5b9c9fd51ba14838ac2e6b222d487106-r",
				nil)
			expectedRequest.Header.Add("Authorization", "Bearer defgh")

			request, _ := NewTokenRevocationRequest(cfg)

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("Api Server Requests", func() {
		BeforeEach(func() {
			cfg.AccessToken = "access-token"
		})

		Describe("NewPutValueRequest", func() {
			It("Returns a request for the put-value endpoint", func() {
				requestBody := bytes.NewReader([]byte(`{"type":"value","name":"my-name","value":"my-value","overwrite":true}`))
				expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/data", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewPutValueRequest(cfg, "my-name", "my-value", true)

				Expect(request).To(Equal(expectedRequest))
			})

			It("Returns a request that will not overwrite", func() {
				requestBody := bytes.NewReader([]byte(`{"type":"value","name":"my-name","value":"my-value","overwrite":false}`))
				expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/data", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewPutValueRequest(cfg, "my-name", "my-value", false)

				Expect(request).To(Equal(expectedRequest))

			})
		})

		Describe("NewPutPasswordRequest", func() {
			It("Returns a request for the put-password endpoint", func() {
				requestBody := bytes.NewReader([]byte(`{"type":"password","name":"my-name","value":"my-password","overwrite":true}`))
				expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/data", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewPutPasswordRequest(cfg, "my-name", "my-password", true)

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewPutCertificateRequest", func() {
			It("Returns a request for the put-certificate endpoint", func() {
				json := fmt.Sprintf(`{"type":"certificate","name":"my-name","value":{"ca":"%s","certificate":"%s","private_key":"%s"},"overwrite":true}`,
					"my-ca", "my-cert", "my-priv")
				requestBody := bytes.NewReader([]byte(json))
				expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/data", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewPutCertificateRequest(cfg, "my-name", "my-ca", "my-cert", "my-priv", true)

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewPutRsaSshRequest", func() {
			Describe("of type SSH", func() {
				It("Returns a request for the put-rsa-ssh endpoint", func() {
					json := fmt.Sprintf(`{"type":"%s","name":"my-name","value":{"public_key":"%s","private_key":"%s"},"overwrite":true}`,
						"ssh", "my-pub", "my-priv")
					requestBody := bytes.NewReader([]byte(json))
					expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/data", requestBody)
					expectedRequest.Header.Set("Content-Type", "application/json")
					expectedRequest.Header.Set("Authorization", "Bearer access-token")

					request := NewPutRsaSshRequest(cfg, "my-name", "ssh", "my-pub", "my-priv", true)

					Expect(request).To(Equal(expectedRequest))
				})
			})

			Describe("of type RSA", func() {
				It("Returns a request for the put-rsa-ssh endpoint", func() {
					json := fmt.Sprintf(`{"type":"%s","name":"my-name","value":{"public_key":"%s","private_key":"%s"},"overwrite":true}`,
						"rsa", "my-pub", "my-priv")
					requestBody := bytes.NewReader([]byte(json))
					expectedRequest, _ := http.NewRequest("PUT", "http://example.com/api/v1/data", requestBody)
					expectedRequest.Header.Set("Content-Type", "application/json")
					expectedRequest.Header.Set("Authorization", "Bearer access-token")

					request := NewPutRsaSshRequest(cfg, "my-name", "rsa", "my-pub", "my-priv", true)

					Expect(request).To(Equal(expectedRequest))
				})
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
				expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/ca?name=my-name&current=true", nil)
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewGetCaRequest(cfg, "my-name")

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewGenerateSecretRequest", func() {
			It("returns a request with only overwrite", func() {
				requestBody := bytes.NewReader([]byte(`{"type":"my-type","overwrite":true,"parameters":{}}`))
				expectedRequest, _ := http.NewRequest("POST", "http://example.com/api/v1/data/my-name", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				params := models.SecretParameters{}
				request := NewGenerateSecretRequest(cfg, "my-name", params, "my-type", true)

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
					"type":"password",
					"overwrite":false,
					"parameters": {
						"exclude_special": true,
						"exclude_number": true,
						"exclude_upper": true,
						"exclude_lower": true,
						"length": 42
					}
				}`

				request := NewGenerateSecretRequest(cfg, "my-name", parameters, "password", false)

				bodyBuffer := new(bytes.Buffer)
				bodyBuffer.ReadFrom(request.Body)
				Expect(bodyBuffer).To(MatchJSON(expectedRequestBody))
				Expect(request.Method).To(Equal("POST"))
				Expect(request.URL.String()).To(Equal("http://example.com/api/v1/data/my-name"))
				Expect(request.Header.Get("Content-Type")).To(Equal("application/json"))
			})
		})

		Describe("NewRegenerateSecretRequest", func() {
			It("returns a request with only regenerate", func() {
				requestBody := bytes.NewReader([]byte(`{"regenerate":true}`))
				expectedRequest, _ := http.NewRequest("POST", "http://example.com/api/v1/data/my-name", requestBody)
				expectedRequest.Header.Set("Content-Type", "application/json")
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewRegenerateSecretRequest(cfg, "my-name")

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewGetSecretRequest", func() {
			It("Returns a request for getting secret", func() {
				expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data?name=my-name&current=true", nil)
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewGetSecretRequest(cfg, "my-name")

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewFindCredentialsBySubstringRequest", func() {
			It("Returns a request for getting secret", func() {
				expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data?name-like=my-name", nil)
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewFindCredentialsBySubstringRequest(cfg, "my-name")

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewFindAllCredentialPathsRequest", func() {
			It("Returns a request for getting all credential paths", func() {
				expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data?paths=true", nil)
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewFindAllCredentialPathsRequest(cfg)

				Expect(request).To(Equal(expectedRequest))
			})
		})

		Describe("NewFindCredentialsByPathRequest", func() {
			It("Returns a request for getting secret", func() {
				expectedRequest, _ := http.NewRequest("GET", "http://example.com/api/v1/data?path=my-path", nil)
				expectedRequest.Header.Set("Authorization", "Bearer access-token")

				request := NewFindCredentialsByPathRequest(cfg, "my-path")

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
