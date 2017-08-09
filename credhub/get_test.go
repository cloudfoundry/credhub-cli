package credhub_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get", func() {
	Describe("Get()", func() {
		var dummy DummyAuth
		var cred credentials.Credential
		var err error
		var serv server.Server

		BeforeEach(func() {
			serv = server.Server{
				ApiUrl:             "http://example.com",
				InsecureSkipVerify: true,
			}
		})

		It("requests the credential by name", func() {
			dummy = DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch := CredHub{
				Server: &serv,
				Auth:   &dummy,
			}
			cred, err = ch.Get("/example-password")
			urlPath := dummy.Request.URL.Path
			Expect(urlPath).To(ContainSubstring("http://example.com/api/v1/data?name=/example-password"))
		})

		Context("when successful", func() {
			It("returns a credential by name", func() {
				responseString := `{
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-password",
      "version_created_at": "2017-01-05T01:01:01Z"
    }`
				dummy = DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch := CredHub{
					Server: &serv,
					Auth:   &dummy,
				}
				cred, err = ch.Get("/example-password")
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
				dummy = DummyAuth{Error: errors.New("Network error occurred")}
				ch := CredHub{
					Server: &serv,
					Auth:   &dummy,
				}
				cred, err = ch.Get("/example-password")

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
				cred, err = ch.Get("/example-password")

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
