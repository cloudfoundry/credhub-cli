package auth

import (
	"net/http"

	"errors"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type DummyServerConfig struct {
	Error error
}

func (d *DummyServerConfig) AuthUrl() (string, error) {
	return "http://example.com/auth/url", d.Error
}

func (d *DummyServerConfig) Client() *http.Client {
	return http.DefaultClient
}

var _ = Describe("Constructors", func() {
	Describe("UaaPasswordGrant()", func() {
		It("constructs a OAuthStrategy auth using password grant", func() {
			config := DummyServerConfig{}
			builder := UaaPasswordGrant("some-client-id", "some-client-secret", "some-username", "some-password")
			strategy, _ := builder(&config)
			auth := strategy.(*OAuthStrategy)
			Expect(auth.ClientId).To(Equal("some-client-id"))
			Expect(auth.ClientSecret).To(Equal("some-client-secret"))
			Expect(auth.Username).To(Equal("some-username"))
			Expect(auth.Password).To(Equal("some-password"))
			Expect(auth.OAuthClient.(*uaa.Client).AuthUrl).To(Equal("http://example.com/auth/url"))
			client := config.Client()
			Expect(auth.OAuthClient.(*uaa.Client).Client).To(BeIdenticalTo(client))
			Expect(auth.ApiClient).To(BeIdenticalTo(client))
		})

		Context("when fetching an Auth URL fails", func() {
			It("returns an error", func() {
				config := DummyServerConfig{
					Error: errors.New("Failed to fetch Auth URL"),
				}
				builder := UaaPasswordGrant("some-client-id", "some-client-secret", "some-username", "some-password")
				_, err := builder(&config)

				Expect(err).To(MatchError("Failed to fetch Auth URL"))
			})

		})
	})
})
