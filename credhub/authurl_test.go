package credhub_test

import (
	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"

	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthUrl()", func() {
	It("should return auth-server url from the /info endpoint", func() {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/info" {
				w.Write([]byte(
					`{
							"auth-server": {
								"url": "https://uaa.example.com:8443"
							},
							"app": {
								"name": "CredHub",
								"version": "0.7.0"
							}
						}`,
				))
			}
		}))
		defer testServer.Close()

		ch := credhubFromConfig(Config{ApiUrl: testServer.URL})

		authUrl, err := ch.AuthUrl()
		Expect(authUrl).To(Equal("https://uaa.example.com:8443"))
		Expect(err).To(BeNil())
	})

	Context("Errors", func() {
		Specify("when ApiUrl is invalid", func() {
			ch := credhubFromConfig(Config{ApiUrl: "://"})

			_, err := ch.AuthUrl()

			Expect(err).ToNot(BeNil())
		})

		Specify("when ApiUrl is inaccessible", func() {
			ch := credhubFromConfig(Config{ApiUrl: "http://localhost:1"})

			_, err := ch.AuthUrl()

			Expect(err).ToNot(BeNil())
		})

		Specify("when /info cannot be parsed", func() {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/info" {
					w.Write([]byte(`INVALID JSON`))
				}
			}))
			defer testServer.Close()

			ch := credhubFromConfig(Config{ApiUrl: testServer.URL})

			_, err := ch.AuthUrl()

			Expect(err).ToNot(BeNil())
		})

		Specify("when auth-server is not returned", func() {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/info" {
					w.Write([]byte(`{}`))
				}
			}))
			defer testServer.Close()

			ch := credhubFromConfig(Config{ApiUrl: testServer.URL})
			_, err := ch.AuthUrl()

			Expect(err).ToNot(BeNil())
		})
	})
})
