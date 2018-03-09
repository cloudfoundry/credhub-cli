package commands_test

import (
	"net/http"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Curl", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("curl", "-p", "api/v1/data")
	ItRequiresAnAPIToBeSet("curl", "-p", "api/v1/data")
	ItAutomaticallyLogsIn("GET", "find_response.json", "/api/v1/data", "curl", "-p", "api/v1/data")

	ItBehavesLikeHelp("curl", "curl", func(session *Session) {
		Expect(session.Err).To(Say("Usage"))
		if runtime.GOOS == "windows" {
			Expect(session.Err).To(Say("credhub-cli.exe \\[OPTIONS\\] curl \\[curl-OPTIONS\\]"))
		} else {
			Expect(session.Err).To(Say("credhub-cli \\[OPTIONS\\] curl \\[curl-OPTIONS\\]"))
		}
	})

	It("displays missing required parameter", func() {
		session := runCommand("curl")

		Eventually(session).Should(Exit(1))

		Expect(session.Err).To(Say("A path must be provided. Please update and retry your request."))
	})

	Context("the user provides an invalid path", func() {
		It("receives what the server returns", func() {
			responseJson := `{"error":"An application error occurred. Please contact your CredHub administrator."}`

			server.RouteToHandler("GET", "/api/v1/data/bogus",
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data/bogus"),
					RespondWith(http.StatusOK, responseJson),
				),
			)

			session := runCommand("curl", "-p", "api/v1/data/bogus")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseJson))
		})
	})

	Context("the user provides a valid path", func() {
		It("receives what the server returns", func() {
			responseJson := `{"id":"2993f622-cb1e-4e00-a267-4b23c273bf3d","name":"/example-password","type":"password","value":"6mRPZB3bAfb8lRpacnXsHfDhlPqFcjH2h9YDvLpL","version_created_at":"2017-01-05T01:01:01Z"}`
			server.RouteToHandler("GET", "/api/v1/data/valid-credential-id",
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data/valid-credential-id"),
					RespondWith(http.StatusOK, responseJson),
				),
			)

			session := runCommand("curl", "-p", "api/v1/data/valid-credential-id")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseJson))
		})

		Context("the user does not specify required parameters", func() {
			It("returns a wrapped error", func() {
				responseJson := `{"error":"The query parameter name is required for this request."}`
				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data"),
						RespondWith(http.StatusBadRequest, responseJson),
					),
				)

				session := runCommand("curl", "-p", "api/v1/data")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say(responseJson))
			})
		})

		Context("when parameters are provided by the user", func() {
			It("returns what the server returns", func() {
				responseJson := `{"data":[{"id":"some-id","name":"example-password","type":"password","value":"secret","version_created_at":"time"}]}`
				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data"),
						RespondWith(http.StatusOK, responseJson),
					),
				)

				session := runCommand("curl", "-p", "api/v1/data?name=/example-password&current=true")

				Eventually(session).Should(Exit(0))
				Eventually(string(session.Out.Contents())).Should(Equal(responseJson + "\n"))
			})
		})
	})

})
