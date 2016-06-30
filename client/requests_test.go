package client_test

import (
	"net/http"

	. "github.com/pivotal-cf/cm-cli/client"

	"bytes"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	Describe("NewPutSecretValueRequest", func() {
		It("Returns a request for the put-secret endpoint", func() {
			requestBody := bytes.NewReader([]byte(`{"type":"value","value":"my-value"}`))
			expectedRequest, _ := http.NewRequest("PUT", "sample.com/api/v1/data/my-name", requestBody)
			expectedRequest.Header.Set("Content-Type", "application/json")

			request := NewPutValueRequest("sample.com", "my-name", "my-value")

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewPutCertificateRequest", func() {
		It("Returns a request for the put-certificate endpoint", func() {
			json := fmt.Sprintf(`{"type":"certificate","certificate":{"ca":"%s","public":"%s","private":"%s"}}`,
				"my-ca", "my-pub", "my-priv")
			requestBody := bytes.NewReader([]byte(json))
			expectedRequest, _ := http.NewRequest("PUT", "sample.com/api/v1/data/my-name", requestBody)
			expectedRequest.Header.Set("Content-Type", "application/json")

			request := NewPutCertificateRequest("sample.com", "my-name", "my-ca", "my-pub", "my-priv")

			Expect(request).To(Equal(expectedRequest))
		})
	})

	Describe("NewPutCaRequest", func() {
		It("Returns a request for the put-root-ca endpoint", func() {
			json := fmt.Sprintf(`{"type":"root","root":{"public":"%s","private":"%s"}}`,
				"my-pub", "my-priv")
			requestBody := bytes.NewReader([]byte(json))
			expectedRequest, _ := http.NewRequest("PUT", "sample.com/api/v1/ca/my-name", requestBody)
			expectedRequest.Header.Set("Content-Type", "application/json")

			request := NewPutCaRequest("sample.com", "my-name", "root", "my-pub", "my-priv")

			Expect(request).To(Equal(expectedRequest))
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
