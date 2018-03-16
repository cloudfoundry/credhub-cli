package commands_test

import (
	"fmt"
	"net/http"

	"runtime"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = FDescribe("Logout", func() {
	AfterEach(func() {
		config.RemoveConfig()
	})

	// TODO fails
	It("marks the access token and refresh token as revoked if no config exists", func() {
		config.RemoveConfig()
		runLogoutCommand()
	})

	It("leaves the access token and refresh token as revoked if config exists and they were already revoked", func() {
		cfg := config.Config{RefreshToken: "revoked", AccessToken: "revoked"}
		config.WriteConfig(cfg)
		runLogoutCommand()
	})

	// TODO fails
	FIt("asks UAA to revoke the token (and UAA succeeds)", func() {
		doRevoke(http.StatusOK)
		runLogoutCommand()
	})

	It("asks UAA to revoke the token (and reports error when UAA fails)", func() {
		doRevoke(http.StatusUnauthorized)

		session := runCommand("logout")
		Eventually(session).Should(Exit(1))
		Eventually(session).Should(Say("Logout Failed"))
	})

	ItBehavesLikeHelp("logout", "o", func(session *Session) {
		Expect(session.Err).To(Say("Usage:"))
		if runtime.GOOS == "windows" {
			Expect(session.Err).To(Say("credhub-cli.exe \\[OPTIONS\\] logout"))
		} else {
			Expect(session.Err).To(Say("credhub-cli \\[OPTIONS\\] logout"))
		}
	})
})

func doRevoke(uaaResponseStatus int) {
	cfg := config.Config{
		RefreshToken: "5b9c9fd51ba14838ac2e6b222d487106-r",
		AccessToken:  "e30K.eyJqdGkiOiIxIn0K.e30K",
		AuthURL:      authServer.URL(),
	}

	cfg.UpdateTrustedCAs([]string{"../test/auth-tls-ca.pem"})
	Expect(config.WriteConfig(cfg)).To(Succeed())
	fmt.Println(authServer.URL())

	authServer.RouteToHandler(
		"DELETE", "/oauth/token/revoke/5b9c9fd51ba14838ac2e6b222d487106-r",
		RespondWith(uaaResponseStatus, ""),
	)
}

func runLogoutCommand() {
	fmt.Println("runLogoutCommand")
	session := runCommand("logout")
	fmt.Println("ranLogoutCommand")
	Eventually(session).Should(Exit(0))
	Eventually(session).Should(Say("Logout Successful"))
	cfg := config.ReadConfig()
	Expect(cfg.AccessToken).To(Equal("revoked"))
	Expect(cfg.RefreshToken).To(Equal("revoked"))
}
