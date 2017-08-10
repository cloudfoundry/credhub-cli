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

var _ = Describe("Find", func() {
	Describe("FindByPath()", func() {
		var dummy DummyAuth
		var serv server.Server
		var creds []credentials.Base
		var err error

		BeforeEach(func() {
			serv = server.Server{
				ApiUrl:             "http://example.com",
				InsecureSkipVerify: true,
			}
		})

		It("requests credentials for a specified path", func() {
			dummy = DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch := CredHub{
				Server: &serv,
				Auth:   &dummy,
			}
			creds, err = ch.FindByPath("/some/example/path")
			urlPath := dummy.Request.URL.Path
			Expect(urlPath).To(Equal("/api/v1/data?path=/some/example/path"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a list of stored credential names which are within the specified path", func() {
				expectedResponse := `{
  "credentials": [
    {
      "version_created_at": "2017-05-09T21:09:26Z",
      "name": "/some/example/path/example-cred-0"
    },
    {
      "version_created_at": "2017-05-09T21:09:07Z",
      "name": "/some/example/path/example-cred-1"
    }
  ]
}`
				dummy = DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString(expectedResponse)),
				}}

				ch := CredHub{
					Server: &serv,
					Auth:   &dummy,
				}
				creds, err = ch.FindByPath("/some/example/path")

				Expect(creds[0].Name).To(Equal("/some/example/path/example-cred-0"))
				Expect(creds[0].VersionCreatedAt).To(Equal("2017-05-09T21:09:26Z"))
				Expect(creds[1].Name).To(Equal("/some/example/path/example-cred-1"))
				Expect(creds[1].VersionCreatedAt).To(Equal("2017-05-09T21:09:07Z"))

			})
		})

		Context("when request fails", func() {
			It("returns an error", func() {
				dummy = DummyAuth{Error: errors.New("Network error occurred")}

				ch := CredHub{
					Server: &serv,
					Auth:   &dummy,
				}

				creds, err = ch.FindByPath("/some/example/path")

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
				creds, err = ch.FindByPath("/some/example/path")

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
