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
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

var _ = Describe("Set", func() {

	Describe("SetCertificate()", func() {
		It("requests to set the certificate", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(authMethod(dummy)))

			certificate := values.Certificate{
				Ca: "some-ca",
			}
			ch.SetCertificate("/example-certificate", certificate, true)

			urlPath := dummy.Request.URL.Path
			Expect(urlPath).To(Equal("/api/v1/data"))
			Expect(dummy.Request.Method).To(Equal(http.MethodPut))

			var requestBody map[string]interface{}
			body, _ := ioutil.ReadAll(dummy.Request.Body)
			json.Unmarshal(body, &requestBody)

			Expect(requestBody["name"]).To(Equal("/example-certificate"))
			Expect(requestBody["type"]).To(Equal("certificate"))
			Expect(requestBody["overwrite"]).To(BeTrue())

			Expect(requestBody["value"].(map[string]interface{})["ca"]).To(Equal("some-ca"))
		})

		Context("when successful", func() {
			It("returns the credential that has been set", func() {
				dummy := &DummyAuth{Response: &http.Response{
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

				ch, _ := New("https://example.com", Auth(authMethod(dummy)))

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
				dummy := &DummyAuth{Error: errors.New("Network error occurred")}
				ch, _ := New("https://example.com", Auth(authMethod(dummy)))
				certificate := values.Certificate{
					Ca: "some-ca",
				}
				_, err := ch.SetCertificate("/example-certificate", certificate, false)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(authMethod(dummy)))
				certificate := values.Certificate{
					Ca: "some-ca",
				}
				_, err := ch.SetCertificate("/example-certificate", certificate, false)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("SetPassword()", func() {
		It("requests to set the password", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(authMethod(dummy)))
			password := values.Password("some-password")

			ch.SetPassword("/example-password", password, false)

			urlPath := dummy.Request.URL.Path
			Expect(urlPath).To(Equal("/api/v1/data"))
			Expect(dummy.Request.Method).To(Equal(http.MethodPut))

			var requestBody map[string]interface{}
			body, _ := ioutil.ReadAll(dummy.Request.Body)
			json.Unmarshal(body, &requestBody)

			Expect(requestBody["name"]).To(Equal("/example-password"))
			Expect(requestBody["type"]).To(Equal("password"))
			Expect(requestBody["value"]).To(BeEquivalentTo("some-password"))
			Expect(requestBody["overwrite"]).To(BeFalse())
		})

		Context("when successful", func() {
			It("returns the credential that has been set", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
		  "id": "some-id",
		  "name": "/example-password",
		  "type": "password",
		  "value": "some-password",
		  "version_created_at": "2017-01-01T04:07:18Z"
		}`)),
				}}

				ch, _ := New("https://example.com", Auth(authMethod(dummy)))

				password := values.Password("some-password")

				cred, _ := ch.SetPassword("/example-password", password, false)

				Expect(cred.Name).To(Equal("/example-password"))
				Expect(cred.Type).To(Equal("password"))

				Expect(cred.Value).To(BeEquivalentTo("some-password"))

			})
		})
		Context("when request fails", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Error: errors.New("Network error occurred")}
				ch, _ := New("https://example.com", Auth(authMethod(dummy)))
				password := values.Password("some-password")

				_, err := ch.SetPassword("/example-password", password, false)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}

				ch, _ := New("https://example.com", Auth(authMethod(dummy)))
				password := values.Password("some-password")

				_, err := ch.SetPassword("/example-password", password, false)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
