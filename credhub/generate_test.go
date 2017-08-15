package credhub_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
)

var _ = Describe("Generate", func() {

	Describe("GenerateCertificate()", func() {
		It("requests to generate the certificate", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))

			cert := generate.Certificate{
				Ca: "some-ca",
			}
			ch.GenerateCertificate("/example-certificate", cert, false)
			urlPath := dummy.Request.URL.Path
			Expect(urlPath).To(Equal("/api/v1/data"))
			Expect(dummy.Request.Method).To(Equal(http.MethodPost))

			var requestBody map[string]interface{}
			body, _ := ioutil.ReadAll(dummy.Request.Body)
			json.Unmarshal(body, &requestBody)

			Expect(requestBody["name"]).To(Equal("/example-certificate"))
			Expect(requestBody["type"]).To(Equal("certificate"))
			Expect(requestBody["overwrite"]).To(BeFalse())
			Expect(requestBody["parameters"].(map[string]interface{})["ca"]).To(Equal("some-ca"))
		})

		Context("when successful", func() {
			It("returns the generated certificate", func() {
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: 200,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
      "id": "some-id",
      "name": "/example-certificate",
      "type": "certificate",
      "value": {
        "ca": "some-ca",
        "certificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
        "private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
      },
      "version_created_at": "2017-01-01T04:07:18Z"
}`)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))

				cert := generate.Certificate{
					Ca: "some-ca",
				}

				generatedCert, err := ch.GenerateCertificate("/example-certificate", cert, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(generatedCert.Id).To(Equal("some-id"))
				Expect(generatedCert.Name).To(Equal("/example-certificate"))
				Expect(generatedCert.Value.Ca).To(Equal("some-ca"))
				Expect(generatedCert.Value.Certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
				Expect(generatedCert.Value.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))

			})
		})

		Context("when request fails", func() {
			var err error
			It("returns an error", func() {
				networkError := errors.New("Network error occurred")
				dummy := &DummyAuth{Error: networkError}
				ch, _ := New("https://example.com", Auth(dummy))

				cert := generate.Certificate{
					Ca: "some-ca",
				}

				_, err = ch.GenerateCertificate("/example-certificate", cert, false)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {

				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("invalid-response")),
				}}
				ch, _ := New("https://example.com", Auth(dummy))

				cert := generate.Certificate{
					Ca: "some-ca",
				}

				_, err := ch.GenerateCertificate("/example-certificate", cert, false)

				Expect(err).To(HaveOccurred())
			})

		})
	})

	Describe("GeneratePassword()", func() {
		It("requests to generate the password", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))

			cert := generate.Password{
				Length: 12,
			}
			ch.GeneratePassword("/example-password", cert, true)
			urlPath := dummy.Request.URL.Path
			Expect(urlPath).To(Equal("/api/v1/data"))

			Expect(dummy.Request.Method).To(Equal(http.MethodPost))

			var requestBody map[string]interface{}
			body, _ := ioutil.ReadAll(dummy.Request.Body)
			json.Unmarshal(body, &requestBody)

			Expect(requestBody["name"]).To(Equal("/example-password"))
			Expect(requestBody["type"]).To(Equal("password"))
			Expect(requestBody["overwrite"]).To(BeTrue())
			Expect(requestBody["parameters"].(map[string]interface{})["length"]).To(BeEquivalentTo(12))
		})

		Context("when successful", func() {
			It("returns the generated password", func() {
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: 200,
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
	      "id": "some-id",
	      "name": "/example-password",
	      "type": "password",
	      "value": "some-password",
	      "version_created_at": "2017-01-01T04:07:18Z"
	}`)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))

				p := generate.Password{
					Length: 12,
				}

				generatedCert, err := ch.GeneratePassword("/example-password", p, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(generatedCert.Id).To(Equal("some-id"))
				Expect(generatedCert.Type).To(Equal("password"))
				Expect(generatedCert.Name).To(Equal("/example-password"))
				Expect(generatedCert.Value).To(BeEquivalentTo("some-password"))
			})
		})

		Context("when request fails to complete", func() {
			var err error
			It("returns an error", func() {
				networkError := errors.New("Network error occurred")
				dummy := &DummyAuth{Error: networkError}
				ch, _ := New("https://example.com", Auth(dummy))

				_, err = ch.GeneratePassword("/example-password", generate.Password{}, false)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {

				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("invalid-response")),
				}}
				ch, _ := New("https://example.com", Auth(dummy))

				_, err := ch.GeneratePassword("/example-password", generate.Password{}, false)

				Expect(err).To(HaveOccurred())
			})

		})
	})
})
