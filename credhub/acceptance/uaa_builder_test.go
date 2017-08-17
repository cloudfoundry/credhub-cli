package acceptance_test

import (
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UaaBuilder", func() {
	Describe("AuthBuilder", func() {
		It("builds an OAuthStrategy using existing tokens", func() {
			oauth := credhubClient.Auth.(*auth.OAuthStrategy)
			err := oauth.Login()
			Expect(err).ToNot(HaveOccurred())

			accessToken := oauth.AccessToken()
			refreshToken := oauth.RefreshToken()

			builder := uaa.AuthBuilder("credhub_cli", "", "credhub", "password", accessToken, refreshToken)
			ch, err := credhub.New("https://localhost:9000", credhub.SkipTLSValidation(), credhub.AuthBuilder(builder))
			Expect(err).ToNot(HaveOccurred())

			oauth = ch.Auth.(*auth.OAuthStrategy)
			_, err = ch.FindByPath("/something")
			Expect(err).ToNot(HaveOccurred())

			Expect(oauth.AccessToken()).To(Equal(accessToken))
			Expect(oauth.RefreshToken()).To(Equal(refreshToken))
		})
	})

	Describe("ClientCredentialsBuilder", func() {
		It("builds an OAuthStrategy using client credentials", func() {
			builder := uaa.ClientCredentialsGrantBuilder("credhub_client", "secret")
			ch, err := credhub.New("https://localhost:9000", credhub.SkipTLSValidation(), credhub.AuthBuilder(builder))
			Expect(err).ToNot(HaveOccurred())

			_, err = ch.FindByPath("/something")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
