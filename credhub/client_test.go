package credhub_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	"crypto/x509"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

var _ = Describe("Client()", func() {

	Context("Errors", func() {
		Specify("when ApiUrl is invalid", func() {
			ch := credhubFromConfig(credhub.Config{ApiUrl: "://"})
			_, err := ch.Client()
			Expect(err).ToNot(BeNil())
		})

		Specify("when CaCerts are invalid", func() {
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

			ch := credhubFromConfig(credhub.Config{ApiUrl: "https://example.com", CaCerts: caCerts})
			_, err := ch.Client()
			Expect(err).ToNot(BeNil())
		})
	})

	Context("Given HTTP ApiUrl", func() {
		It("should return a simple http.Client", func() {
			ch := credhubFromConfig(credhub.Config{ApiUrl: "http://example.com"})
			client, err := ch.Client()

			Expect(err).To(BeNil())
			Expect(client.Transport).To(BeNil())
			Expect(client.Timeout).To(Equal(45 * time.Second))
		})
	})

	Context("Given HTTPS ApiUrl", func() {

		Context("With CaCerts", func() {
			It("should return a http.Client with tls.Config with RootCAs", func() {
				fixturePath := "./fixtures/"
				caCertFiles := []string{
					"auth-tls-ca.pem",
					"server-tls-ca.pem",
					"extra-ca.pem",
				}
				var caCerts []string
				expectedRootCAs := x509.NewCertPool()
				for _, caCertFile := range caCertFiles {
					caCertBytes, err := ioutil.ReadFile(fixturePath + caCertFile)
					if err != nil {
						Fail("Couldn't read certificate " + caCertFile + ": " + err.Error())
					}

					caCerts = append(caCerts, string(caCertBytes))
					expectedRootCAs.AppendCertsFromPEM(caCertBytes)
				}

				ch := credhubFromConfig(credhub.Config{ApiUrl: "https://example.com", CaCerts: caCerts})
				client, err := ch.Client()

				Expect(err).To(BeNil())

				transport := client.Transport.(*http.Transport)
				tlsConfig := transport.TLSClientConfig

				Expect(client.Timeout).To(Equal(45 * time.Second))

				Expect(tlsConfig.InsecureSkipVerify).To(BeFalse())
				Expect(tlsConfig.PreferServerCipherSuites).To(BeTrue())
				Expect(tlsConfig.RootCAs.Subjects()).To(ConsistOf(expectedRootCAs.Subjects()))
			})
		})

		Context("With InsecureSkipVerify", func() {
			It("should return a http.Client with tls.Config without RootCAs", func() {
				ch := credhubFromConfig(credhub.Config{ApiUrl: "https://example.com", InsecureSkipVerify: true})
				client, err := ch.Client()

				Expect(err).To(BeNil())

				transport := client.Transport.(*http.Transport)
				tlsConfig := transport.TLSClientConfig

				Expect(client.Timeout).To(Equal(45 * time.Second))

				Expect(tlsConfig.InsecureSkipVerify).To(BeTrue())
				Expect(tlsConfig.PreferServerCipherSuites).To(BeTrue())
			})
		})

	})
})

func credhubFromConfig(config credhub.Config) *credhub.CredHub {
	return credhub.New(&config, noopAuth)
}

func noopAuth(auth.ServerConfig) auth.Auth {
	return nil
}
