package client_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
)

var _ = Describe("#NewHttpClient", func() {
	It("returns http client when a url specifies http scheme", func() {
		config := config.Config{
			ApiURL: "http://foo.bar",
		}

		httpClient := client.NewHttpClient(config.ApiURL)
		Expect(httpClient.Transport).To(BeNil())
	})

	It("returns https client when url scheme is https", func() {
		config := config.Config{
			ApiURL: "https://foo.bar",
		}

		httpsClient := client.NewHttpClient(config.ApiURL)
		Expect(httpsClient.Transport).To(Not(BeNil()))
	})

	It("requires tls verification for https client", func() {
		config := config.Config{
			ApiURL: "https://foo.bar",
		}

		httpsClient := client.NewHttpClient(config.ApiURL)
		Expect(httpsClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify).To(BeFalse())
	})

	It("prefers server cipher suites for https client", func() {
		config := config.Config{
			ApiURL: "https://foo.bar",
		}

		httpsClient := client.NewHttpClient(config.ApiURL)
		Expect(httpsClient.Transport.(*http.Transport).TLSClientConfig.PreferServerCipherSuites).To(BeTrue())
	})
})
