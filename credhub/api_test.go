package credhub_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type dummyAuth struct {
	Server   *server.Server
	Request  *http.Request
	Response *http.Response
	Error    error
}

func (d *dummyAuth) Do(req *http.Request) (*http.Response, error) {
	d.Request = req

	return d.Response, d.Error
}

var _ = Describe("Api", func() {
	Context("New()", func() {
		It("should assign Server and Auth", func() {
			dummy := func(s *server.Server) auth.Auth {
				return &dummyAuth{Server: s}
			}
			s := &server.Server{ApiUrl: "http://example.com"}

			ch := New(s, dummy)

			Expect(ch.Server).To(BeIdenticalTo(s))

			da, ok := ch.Auth.(*dummyAuth)

			Expect(ok).To(BeTrue())
			Expect(da.Server).To(BeIdenticalTo(s))
		})
	})

	Context("Request()", func() {
		var (
			mockAuth *dummyAuth
			ch       *CredHub
		)

		BeforeEach(func() {
			mockAuth = &dummyAuth{}
			ch = &CredHub{
				Server: &server.Server{ApiUrl: "http://example.com/"},
				Auth:   mockAuth,
			}
		})

		It("should send the requested using the provided auth to the ApiUrl", func() {
			payload := map[string]interface{}{
				"some-field":  1,
				"other-field": "blah",
			}

			mockAuth.Response = &http.Response{}
			mockAuth.Error = errors.New("Some error")

			response, err := ch.Request("PATCH", "/api/v1/some-endpoint", payload)

			Expect(response).To(Equal(mockAuth.Response))
			Expect(err).To(Equal(mockAuth.Error))

			Expect(mockAuth.Request.Method).To(Equal("PATCH"))
			Expect(mockAuth.Request.URL.String()).To(Equal("http://example.com/api/v1/some-endpoint"))

			body, err := ioutil.ReadAll(mockAuth.Request.Body)

			Expect(err).To(BeNil())
			Expect(body).To(MatchJSON(`{"some-field": 1, "other-field": "blah"}`))
		})

		It("fails to send the request when the body cannot be marshalled to JSON", func() {
			_, err := ch.Request("PATCH", "/api/v1/some-endpoint", &NotMarshallable{})
			Expect(err).To(HaveOccurred())
		})

		It("fails to send when the request method is invalid", func() {
			_, err := ch.Request(" ", "/api/v1/some-endpoint", nil)
			Expect(err).To(HaveOccurred())
		})
	})
})

type NotMarshallable struct{}

func (u *NotMarshallable) MarshalJSON() ([]byte, error) {
	return nil, errors.New("I cannot be marshalled")
}

var _ json.Marshaler = new(NotMarshallable)
