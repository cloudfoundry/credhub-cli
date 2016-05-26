package client_test

import (
	"net/http"

	. "github.com/pivotal-cf/cm-cli/client"

	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/cm-cli/models"
)

var _ = Describe("API", func() {
	Describe("NewInfoRequest", func() {
		It("Returns a request for the info endpoint", func() {
			httpRequest, _ := http.NewRequest("GET", "fake_target.com/info", nil)

			request := NewInfoRequest("fake_target.com")

			Expect(request).To(Equal(httpRequest))
		})
	})

	Describe("NewPutSecretRequest", func() {
		It("Returns a request for the put-secret endpoint", func() {
			requestBody := bytes.NewReader([]byte(`{"value":"my-value"}`))
			httpRequest, _ := http.NewRequest("PUT", "sample.com/api/v1/data/my-name", requestBody)
			httpRequest.Header.Set("Content-Type", "application/json")

			request := NewPutSecretRequest("sample.com", "my-name", "my-value")

			Expect(request).To(Equal(httpRequest))
		})
	})

	Describe("NewGenerateSecretRequest", func() {
		It("returns a request with no parameters", func() {
			requestBody := bytes.NewReader([]byte(`{"parameters":{}}`))

			httpRequest, _ := http.NewRequest("POST", "sample.com/api/v1/data/my-name", requestBody)
			httpRequest.Header.Set("Content-Type", "application/json")

			emptyParameters := models.SecretParameters{}
			request := NewGenerateSecretRequest("sample.com", "my-name", emptyParameters)

			Expect(request).To(Equal(httpRequest))
		})

		It("returns a request with length", func() {
			requestBody := bytes.NewReader([]byte(`{"parameters":{"length":42}}`))

			httpRequest, _ := http.NewRequest("POST", "sample.com/api/v1/data/my-name", requestBody)
			httpRequest.Header.Set("Content-Type", "application/json")

			withLengthParameters := models.SecretParameters{
				Length: 42,
			}

			request := NewGenerateSecretRequest("sample.com", "my-name", withLengthParameters)

			Expect(request).To(Equal(httpRequest))
		})
	})

	Describe("NewGetSecretRequest", func() {
		It("Returns a request for getting secret", func() {
			httpRequest, _ := http.NewRequest("GET", "sample.com/api/v1/data/my-name", nil)

			request := NewGetSecretRequest("sample.com", "my-name")

			Expect(request).To(Equal(httpRequest))
		})
	})

	Describe("NewDeleteSecretRequest", func() {
		It("Returns a request for deleting", func() {
			httpRequest, _ := http.NewRequest("DELETE", "sample.com/api/v1/data/my-name", nil)

			request := NewDeleteSecretRequest("sample.com", "my-name")

			Expect(request).To(Equal(httpRequest))
		})
	})
})
