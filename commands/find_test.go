package commands_test

import (
	"net/http"

	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/credhub-cli/commands"
)

var _ = Describe("Find", func() {
	Describe("Help", func() {
		It("displays help", func() {
			session := runCommand("find", "-h")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("Usage"))
			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("credhub-cli.exe \\[OPTIONS\\] find \\[find-OPTIONS\\]"))
			} else {
				Expect(session.Err).To(Say("credhub-cli \\[OPTIONS\\] find \\[find-OPTIONS\\]"))
			}
		})

		It("short flags", func() {
			Expect(commands.FindCommand{}).To(SatisfyAll(
				commands.HaveFlag("name-like", "n"),
			))
		})

		It("displays missing 'n' option as required parameter", func() {
			session := runCommand("find")

			Eventually(session).Should(Exit(1))
			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flag `/n, /name-like' was not specified"))
			} else {
				Expect(session.Err).To(Say("the required flag `-n, --name-like' was not specified"))
			}
		})
	})

	Describe("finds a set of credentials matching a supplied string", func() {
		It("gets a list of string secret names and last-modified dates", func() {
			responseJson := `{
				"credentials": [
						{
							"name": "dan.password",
							"updated_at": "2016-09-06T23:26:58Z"
						},
						{
							"name": "deploy1/dan/id.key",
							"updated_at": "2016-09-06T23:26:58Z"
						}
				]
			}`
			responseTable := "Name                 Updated Date\n" +
                       "dan.password         2016-09-06T23:26:58Z\n" +
                       "deploy1/dan/id.key   2016-09-06T23:26:58Z"


			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "name-like=dan"),
					RespondWith(http.StatusOK, responseJson),
				),
			)

			session := runCommand("find", "-n", "dan")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseTable))
		})
	})

})
