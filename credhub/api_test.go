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
		It("sets Auth to some default value", func() {
			ch, err := New("http://example.com")
			Expect(err).ToNot(HaveOccurred())
			Expect(ch.Auth).ToNot(BeNil())
		})

		Context("when the Auth option is used", func() {
			It("sets the Auth", func() {
				expectedAuth := &DummyAuth{}
				ch, err := New("http://example.com", Auth(expectedAuth))

				Expect(err).ToNot(HaveOccurred())

				auth, ok := ch.Auth.(*DummyAuth)

				Expect(ok).To(BeTrue())
				Expect(auth).To(BeIdenticalTo(expectedAuth))
			})

			It("ignores the auth builder option", func() {
				builderCalled := false
				builder := func(config auth.ServerConfig) auth.Auth {
					builderCalled = true
					return nil
				}

				_, err := New("http://example.com", Auth(&DummyAuth{}), AuthBuilder(builder))

				Expect(err).ToNot(HaveOccurred())

				Expect(builderCalled).To(BeFalse())
			})
		})

		Context("when the auth builder is used", func() {
			It("invokes the auth builder", func() {
				dummyBuilder := func(config auth.ServerConfig) auth.Auth {
					return &DummyAuth{Config: config}
				}

				ch, err := New("http://example.com", AuthBuilder(dummyBuilder))
				Expect(err).ToNot(HaveOccurred())

				da, ok := ch.Auth.(*DummyAuth)

				Expect(ok).To(BeTrue())
				Expect(da.Config).To(BeIdenticalTo(ch))
			})
		})

		It("returns an error when the ApiURL is invalid", func() {
			ch, err := New("://example.com")
			Expect(err).To(HaveOccurred())
			Expect(ch).To(BeNil())

		})

		It("returns an error when CaCerts are invalid", func() {
			fixturePath := "./fixtures/"
			caCertFiles := []string{
				"auth-tls-ca.pem",
				"server-tls-ca.pem",
				"extra-ca.pem",
			}
			var caCerts []string
			for _, caCertFile := range caCertFiles {
				caCertBytes, err := ioutil.ReadFile(fixturePath + caCertFile)
				if err != nil {
					Fail("Couldn't read certificate " + caCertFile + ": " + err.Error())
				}

				caCerts = append(caCerts, string(caCertBytes))
			}
			caCerts = append(caCerts, "invalid certificate")

			_, err := New("https://example.com", CACerts(caCerts))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Request()", func() {
		var (
			mockAuth *DummyAuth
			ch       *CredHub
		)

		BeforeEach(func() {
			mockAuth = &DummyAuth{}
			ch, _ = New("http://example.com/", Auth(mockAuth))
		})

		It("should send the requested using the provided auth to the ApiURL", func() {
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
