package credhub_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("New()", func() {
	It("sets Auth to some default value", func() {
		ch, err := New("http://example.com", ServerVersion("2.2.2"))
		Expect(err).ToNot(HaveOccurred())
		Expect(ch.Auth).ToNot(BeNil())
	})

	Context("when the AuthURL option is used", func() {
		It("returns the url without hitting the api", func() {
			apiHit := false

			credHubServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				apiHit = true
				w.Write([]byte(`{
					"auth-server": {"url": "https://uaa.example.com:8443"},
					"app": {"name": "CredHub", "version": "0.7.0"}
				}`))
			}))

			defer credHubServer.Close()

			var authConfig auth.Config

			authBuilder := func(config auth.Config) (auth.Strategy, error) {
				authConfig = config
				return http.DefaultClient, nil
			}

			ch, err := New(credHubServer.URL, AuthURL("https://some-auth-url.com"), Auth(authBuilder), ServerVersion("2.2.2"))
			Expect(err).ToNot(HaveOccurred())

			Expect(apiHit).To(BeFalse())
			Expect(err).ToNot(HaveOccurred())

			// Ensure builder is invoked with the correct AuthURL
			Expect(authConfig.AuthURL()).To(Equal("https://some-auth-url.com"))
			Expect(ch.AuthURL()).To(Equal("https://some-auth-url.com"))

		})

		Context("when the url is invalid", func() {
			It("returns an error", func() {
				ch, err := New("http://example.com", AuthURL("://some-auth-url.com"), ServerVersion("2.2.2"))
				Expect(ch).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("when the auth builder is used", func() {
		It("invokes the auth builder", func() {
			dummyBuilder := func(config auth.Config) (auth.Strategy, error) {
				return &DummyAuth{Config: config}, nil
			}

			ch, err := New("http://example.com", Auth(dummyBuilder), ServerVersion("2.2.2"))
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

		_, err := New("https://example.com", CaCerts(caCerts...), ServerVersion("2.2.2"))
		Expect(err).To(HaveOccurred())
	})

	Context("With ServerVersion option provided", func() {
		It("stores that server version in the client", func() {
			ch, _ := New("https://example.com", ServerVersion("2.2.2"))

			Expect(ch.ServerVersion).To(Equal("2.2.2"))
		})
	})

	Context("Without ServerVersion option provided", func() {
		It("fetches the server version from the Credhub server", func() {
			var apiHit bool

			credHubServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				apiHit = true
				w.Write([]byte(`{
					"auth-server": {"url": "https://uaa.example.com:8443"},
					"app": {"name": "CredHub", "version": "4.4.4"}
				}`))
			}))

			defer credHubServer.Close()

			var authConfig auth.Config

			authBuilder := func(config auth.Config) (auth.Strategy, error) {
				authConfig = config
				return http.DefaultClient, nil
			}

			ch, err := New(credHubServer.URL, AuthURL("https://some-auth-url.com"), Auth(authBuilder))
			Expect(err).ToNot(HaveOccurred())

			Expect(ch.ServerVersion).To(Equal("4.4.4"))
		})
	})
})
