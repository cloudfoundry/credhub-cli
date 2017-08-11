package auth

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type DummyServerConfig struct{}

func (d *DummyServerConfig) AuthUrl() (string, error) {
	return "http://example.com/auth/url", nil
}

func (d *DummyServerConfig) Client() *http.Client {
	return http.DefaultClient
}

var _ = Describe("Constructors", func() {
	Describe("UaaPasswordGrant()", func() {
		It("constructs a Uaa auth using password grant", func() {
			config := DummyServerConfig{}
			method := UaaPasswordGrant("some-client-id", "some-client-secret", "some-username", "some-password")
			auth := method(&config).(*Uaa)
			Expect(auth.ClientId).To(Equal("some-client-id"))
			Expect(auth.ClientSecret).To(Equal("some-client-secret"))
			Expect(auth.Username).To(Equal("some-username"))
			Expect(auth.Password).To(Equal("some-password"))
			Expect(auth.UaaClient.(*uaa.Client).AuthUrl).To(Equal("http://example.com/auth/url"))
			client := config.Client()
			Expect(auth.UaaClient.(*uaa.Client).Client).To(BeIdenticalTo(client))
			Expect(auth.ApiClient).To(BeIdenticalTo(client))
		})
	})
})
