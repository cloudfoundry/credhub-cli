package client_test

import (
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
	It("returns https client when", func() {
		config := config.Config{
			ApiURL: "https://foo.bar",
		}

		httpsClient := client.NewHttpClient(config.ApiURL)
		Expect(httpsClient.Transport).To(Not(BeNil()))
	})
})
