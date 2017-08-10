package credhub_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Set", func() {
	var serv server.Server
	var dummy DummyAuth
	BeforeEach(func() {
		serv = server.Server{
			ApiUrl:             "http://example.com",
			InsecureSkipVerify: true,
		}
	})

	Describe("Set()", func() {
		It("requests to set the certificate", func() {
			dummy = DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch := CredHub{
				Server: &serv,
				Auth:   &dummy,
			}

			certificate := values.Certificate{
				Ca: "some-ca",
			}
			ch.SetCertificate("/example-certificate", certificate, false)

			urlPath := dummy.Request.URL.Path
			Expect(urlPath).To(Equal("/api/v1/data"))
			Expect(dummy.Request.Method).To(Equal(http.MethodPut))

			var requestBody map[string]interface{}
			body, _ := ioutil.ReadAll(dummy.Request.Body)
			json.Unmarshal(body, &requestBody)

			Expect(requestBody["name"]).To(Equal("/example-certificate"))
			Expect(requestBody["type"]).To(Equal("certificate"))
			Expect(requestBody["value"].(map[string]interface{})["ca"]).To(Equal("some-ca"))
		})

		Context("when successful", func() {
			It("returns the credential that has been set", func() {
				dummy = DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
		  "id": "some-id",
		  "name": "/example-certificate",
		  "type": "certificate",
		  "value": {
		    "ca": "some-ca",
		    "certificate": "some-certificate",
		    "private_key": "some-private-key"
		  },
		  "version_created_at": "2017-01-01T04:07:18Z"
		}`)),
				}}

				ch := CredHub{
					Server: &serv,
					Auth:   &dummy,
				}

				certificate := values.Certificate{
					Certificate: "some-cert",
				}
				cred, _ := ch.SetCertificate("/example-certificate", certificate, false)

				Expect(cred.Name).To(Equal("/example-certificate"))
				Expect(cred.Type).To(Equal("certificate"))
				Expect(cred.Value.Ca).To(Equal("some-ca"))
				Expect(cred.Value.Certificate).To(Equal("some-certificate"))
				Expect(cred.Value.PrivateKey).To(Equal("some-private-key"))
			})
		})
		Context("when request fails", func() {
			It("returns an error", func() {
				dummy = DummyAuth{Error: errors.New("Network error occurred")}
				ch := CredHub{
					Server: &serv,
					Auth:   &dummy,
				}
				certificate := values.Certificate{
					Ca: "some-ca",
				}
				_, err := ch.SetCertificate("/example-certificate", certificate, false)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy = DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}

				ch := CredHub{
					Server: &serv,
					Auth:   &dummy,
				}
				certificate := values.Certificate{
					Ca: "some-ca",
				}
				_, err := ch.SetCertificate("/example-certificate", certificate, false)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
