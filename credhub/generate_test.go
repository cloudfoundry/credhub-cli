package credhub_test

import (
	"errors"
	"net/http"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"

	"bytes"
	"io/ioutil"

	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generate", func() {
	var serv server.Server
	BeforeEach(func() {
		serv = server.Server{
			ApiUrl:             "http://example.com",
			InsecureSkipVerify: true,
		}
	})
	Describe("GenerateCertificate()", func() {
		It("requests to generate the certificate", func() {
			dummy := DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch := CredHub{
				Server: &serv,
				Auth:   &dummy,
			}
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
			Expect(requestBody["parameters"].(map[string]interface{})["ca"]).To(Equal("some-ca"))
		})

		Context("when successful", func() {
			It("returns the generated certificate", func() {
				dummy := DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
  "data": [
    {
      "id": "some-id",
      "name": "/example-certificate",
      "type": "certificate",
      "value": {
        "ca": "some-ca",
        "certificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
        "private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
      },
      "version_created_at": "2017-01-01T04:07:18Z"
    }
  ]
}`)),
				}}

				sserv := server.Server{
					ApiUrl:             "http://example.com",
					InsecureSkipVerify: true,
				}

				ch := CredHub{
					Server: &sserv,
					Auth:   &dummy,
				}
				cert := generate.Certificate{
					Ca: "some-ca",
				}

				generatedCert, _ := ch.GenerateCertificate("/example-certificate", cert, false)
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
				dummy := DummyAuth{Error: networkError}
				ch := CredHub{
					Server: &serv,
					Auth:   &dummy,
				}

				cert := generate.Certificate{
					Ca: "some-ca",
				}

				_, err = ch.GenerateCertificate("/example-certificate", cert, false)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {

				dummy := DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("invalid-response")),
				}}
				ch := CredHub{
					Server: &serv,
					Auth:   &dummy,
				}

				cert := generate.Certificate{
					Ca: "some-ca",
				}

				_, err := ch.GenerateCertificate("/example-certificate", cert, false)

				Expect(err).To(HaveOccurred())
			})

		})

		Context("when response body does not contain a certificate", func() {

			It("returns an error", func() {

				dummy := DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString(`{"data":[]}`)),
				}}
				ch := CredHub{
					Server: &serv,
					Auth:   &dummy,
				}

				cert := generate.Certificate{
					Ca: "some-ca",
				}

				_, err := ch.GenerateCertificate("/example-certificate", cert, false)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
