package commands_test

import (
	"net/http"

	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/credhub-cli/config"
)

var _ = Describe("Logout", func() {
	AfterEach(func() {
		config.RemoveConfig()
	})

	It("marks the access token and refresh token as revoked if no config exists", func() {
		config.RemoveConfig()
		runLogoutCommand()
	})

	It("leaves the access token and refresh token as revoked if config exists and they were already revoked", func() {
		cfg := config.Config{RefreshToken: "revoked", AccessToken: "revoked"}
		config.WriteConfig(cfg)
		runLogoutCommand()
	})

	It("asks UAA to revoke the refresh token (and UAA succeeds)", func() {
		doRevoke(http.StatusOK)
	})

	It("asks UAA to revoke the refresh token (and reports no error when UAA fails)", func() {
		doRevoke(http.StatusUnauthorized)
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
		RefreshToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6ImxlZ2FjeS10b2tlbi1rZXkiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiI1YjljOWZkNTFiYTE0ODM4YWMyZTZiMjIyZDQ4NzEwNi1yIiwic3ViIjoiYzE0ZGJjZGQtNzNkOC00ZDdjLWI5NDctYzM4ODVhODAxYzY2Iiwic2NvcGUiOlsiY3JlZGh1Yi53cml0ZSIsImNyZWRodWIucmVhZCJdLCJpYXQiOjE0NzEzMTAwMTIsImV4cCI6MTQ3MTM5NjQxMiwiY2lkIjoiY3JlZGh1YiIsImNsaWVudF9pZCI6ImNyZWRodWIiLCJpc3MiOiJodHRwczovLzUyLjIwNC40OS4xMDc6ODQ0My9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsInJldm9jYWJsZSI6dHJ1ZSwiZ3JhbnRfdHlwZSI6InBhc3N3b3JkIiwidXNlcl9uYW1lIjoiY3JlZGh1Yl9jbGkiLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX2lkIjoiYzE0ZGJjZGQtNzNkOC00ZDdjLWI5NDctYzM4ODVhODAxYzY2IiwicmV2X3NpZyI6ImQ3MTkyZmUxIiwiYXVkIjpbImNyZWRodWIiXX0.UAp6Ou24f18mdE0XOqG9RLVWZAx3khNHHPeHfuzmcOUYojtILa0_izlGVHhCtNx07f4M9pcRKpo-AijXRw1vSimSTHBeVCDjuuc2nBdznIMhyQSlPpd2stW-WG7Gix82K4gy4oCb1wlTqsK3UKGYoy8JWs6XZqhoZZ6JZM7-Xjj2zag3Q4kgvEBReWC5an_IP6SeCpNt5xWvGdxtTz7ki1WPweUBy0M73ZjRi9_poQT2JmeSIbrePukkfsfCxHG1vM7ApIdzzhdCx6T_KmmMU3xHqhpI_ueLOuvfHjdBinm2atypeTHD83yRRFxhfjRsG1-XguTn-lo_Z2Jis89r5g",
		AccessToken:  "myAccessToken",
		AuthURL:      authServer.URL(),
	}
	config.WriteConfig(cfg)

	authServer.AppendHandlers(
		CombineHandlers(
			VerifyRequest("DELETE", "/oauth/token/revoke/5b9c9fd51ba14838ac2e6b222d487106-r"),
			RespondWith(uaaResponseStatus, ""),
		),
	)
	runLogoutCommand()
}

func runLogoutCommand() {
	session := runCommand("logout")
	Eventually(session).Should(Exit(0))
	Eventually(session).Should(Say("Logout Successful"))
	cfg := config.ReadConfig()
	Expect(cfg.AccessToken).To(Equal("revoked"))
	Expect(cfg.RefreshToken).To(Equal("revoked"))
}
