package credhub_test

import (
	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"

	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthUrl()", func() {
	Context("Errors", func() {

		Specify("when ApiURL is inaccessible", func() {
			ch, _ := New("http://localhost:1")
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

			ch, _ := New(testServer.URL)
			_, err := ch.AuthUrl()

			Expect(err).ToNot(BeNil())
		})
	})
})
