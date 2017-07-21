package client_test

import (
	"net/http"

	"io/ioutil"

	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("#NewHttpClient", func() {
	var (
		cfg    config.Config
		testCa string
	)

	BeforeEach(func() {
		ca, _ := ioutil.ReadFile("../test/server-tls-ca.pem")
		testCa = string(ca)
	})

	It("returns http client when a url specifies http scheme", func() {
		cfg = config.Config{
			ApiURL: "http://foo.bar",
		}

		httpClient := client.NewHttpClient(cfg)
		Expect(httpClient.Transport).To(BeNil())
	})

	It("returns https client when url scheme is https", func() {
		cfg = config.Config{
			ApiURL: "https://foo.bar",
		}

		httpsClient := client.NewHttpClient(cfg)
		Expect(httpsClient.Transport).To(Not(BeNil()))
	})

	It("requires tls verification for https client", func() {
		cfg = config.Config{
			ApiURL:             "https://foo.bar",
			InsecureSkipVerify: false,
		}

		httpsClient := client.NewHttpClient(cfg)
		Expect(httpsClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify).To(BeFalse())
	})

	It("can skip tls verification for https client", func() {
		cfg = config.Config{
			ApiURL:             "https://foo.bar",
			InsecureSkipVerify: true,
		}

		httpsClient := client.NewHttpClient(cfg)
		Expect(httpsClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify).To(BeTrue())
	})

	It("prefers server cipher suites for https client", func() {
		cfg = config.Config{
			ApiURL: "https://foo.bar",
		}

		httpsClient := client.NewHttpClient(cfg)
		Expect(httpsClient.Transport.(*http.Transport).TLSClientConfig.PreferServerCipherSuites).To(BeTrue())
	})

	It("uses server ca cert in tls connection if provided", func() {
		cfg = config.Config{
			CaCerts: []string{testCa},
			ApiURL:  "https://test.com",
		}

		httpsClient := client.NewHttpClient(cfg)
		Expect(httpsClient.Transport.(*http.Transport).TLSClientConfig.RootCAs).To(Not(BeNil()))
	})

	It("leaves tls Root CA as nil if no cert provided", func() {
		cfg = config.Config{
			CaCerts: []string{},
			ApiURL:  "https://test.com",
		}

		httpsClient := client.NewHttpClient(cfg)
		Expect(httpsClient.Transport.(*http.Transport).TLSClientConfig.RootCAs).To(BeNil())
	})

})
