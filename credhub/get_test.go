package credhub_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
)

var _ = Describe("Get", func() {

	Describe("Get()", func() {
		It("requests the credential by name", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))

			ch.Get("/example-password")
			url := dummy.Request.URL.String()
			Expect(url).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-password"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a credential by name", func() {
				responseString := `{
	"data": [
	{
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-password",
      "version_created_at": "2017-01-05T01:01:01Z"
    }
    ]}`
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))

				cred, err := ch.Get("/example-password")
				Expect(err).To(BeNil())
				Expect(cred.Id).To(Equal("some-id"))
				Expect(cred.Name).To(Equal("/example-password"))
				Expect(cred.Type).To(Equal("password"))
				Expect(cred.Value.(string)).To(Equal("some-password"))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})
		})

		Context("when request fails", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Error: errors.New("Network error occurred")}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.Get("/example-password")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}

				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.Get("/example-password")

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the response body contains an empty list", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":[]}`)),
				}}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.Get("/example-password")

				Expect(err).To(MatchError("response did not contain any credentials"))
			})
		})
	})

	Describe("GetPassword()", func() {
		It("requests the credential by name", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))
			ch.GetPassword("/example-password")
			url := dummy.Request.URL
			Expect(url.String()).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-password"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a password credential", func() {
				responseString := `{
  "data": [
    {
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-password",
      "version_created_at": "2017-01-05T01:01:01Z"
    }]}`
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))
				cred, err := ch.GetPassword("/example-password")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value).To(BeEquivalentTo("some-password"))
			})
		})

		Context("when request fails", func() {
			It("returns an error", func() {
				networkError := errors.New("Network error occurred")
				dummy := &DummyAuth{Error: networkError}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.GetPassword("/example-password")

				Expect(err).To(Equal(networkError))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.GetPassword("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetCertificate()", func() {
		It("requests the credential by name", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))
			ch.GetCertificate("/example-certificate")
			url := dummy.Request.URL
			Expect(url.String()).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-certificate"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a certificate credential", func() {
				responseString := `{
				  "data": [{
	"id": "some-id",
	"name": "/example-certificate",
	"type": "certificate",
	"value": {
		"ca": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"certificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
	},
	"version_created_at": "2017-01-01T04:07:18Z"
}]}`
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))

				cred, err := ch.GetCertificate("/example-certificate")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value.Ca).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
				Expect(cred.Value.Certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
				Expect(cred.Value.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))
			})
		})

		Context("when request fails", func() {
			It("returns an error", func() {
				networkError := errors.New("Network error occurred")
				dummy := &DummyAuth{Error: networkError}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.GetCertificate("/example-certificate")
				Expect(err).To(Equal(networkError))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.GetCertificate("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
