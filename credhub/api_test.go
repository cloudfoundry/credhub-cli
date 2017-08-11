package credhub_test

import (
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Api", func() {
	Context("New()", func() {
		dummy := func(config auth.ServerConfig) auth.Auth {
			return &DummyAuth{Config: config}
		}

		It("should assign Config and Auth", func() {
			config := &Config{ApiUrl: "http://example.com"}

			ch, err := New(config, dummy)

			Expect(err).ToNot(HaveOccurred())
			Expect(ch.Config).To(BeIdenticalTo(config))

			da, ok := ch.Auth.(*DummyAuth)

			Expect(ok).To(BeTrue())
			Expect(da.Config).To(BeIdenticalTo(ch))
		})

		It("returns an error when the ApiUrl is invalid", func() {
			config := &Config{ApiUrl: "://example.com"}

			ch, err := New(config, dummy)
			Expect(err).To(HaveOccurred())
			Expect(ch).To(BeNil())

		})

	})

	Context("Request()", func() {
		var (
			mockAuth *DummyAuth
			ch       *CredHub
		)

		BeforeEach(func() {
			mockAuth = &DummyAuth{}
			ch = credhubWithAuth(Config{ApiUrl: "http://example.com/"}, mockAuth)
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

func credhubFromConfig(config Config) *CredHub {
	c, err := New(&config, noopAuth)

	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	return c
}

func credhubWithAuth(config Config, authentication auth.Auth) *CredHub {
	c, err := New(&config, func(auth.ServerConfig) auth.Auth {
		return authentication
	})

	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	return c
}

func noopAuth(auth.ServerConfig) auth.Auth {
	return nil
}
