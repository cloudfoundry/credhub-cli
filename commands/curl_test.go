package commands_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"net/http"
	"net/url"
)

var _ = Describe("Curl", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("curl", "-p", "api/v1/data")
	ItRequiresAnAPIToBeSet("curl", "-p", "api/v1/data")

	testAutoLogIns := []TestAutoLogin{
		{
			method:              "GET",
			responseFixtureFile: "find_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v1/data",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogIns, "curl", "-p", "api/v1/data")

	ItBehavesLikeHelp("curl", "curl", func(session *Session) {
		Expect(session.Err).To(Say("Usage:\n(.*)\\[OPTIONS\\] curl \\[curl-OPTIONS\\]"))
	})

	It("displays missing required parameter", func() {
		session := runCommand("curl")

		Eventually(session).Should(Exit(1))

		Expect(session.Err).To(Say("A path must be provided. Please update and retry your request."))
	})

	Context("the user provides an invalid path", func() {
		It("receives what the server returns", func() {
			//language=json
			responseJSON := `
        {
          "error": "An application error occurred. Please contact your CredHub administrator."
        }
      `
			server.RouteToHandler("GET", "/api/v1/data/bogus",
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data/bogus"),
					RespondWith(http.StatusOK, responseJSON),
				),
			)

			session := runCommand("curl", "-p", "api/v1/data/bogus")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out.Contents()).Should(MatchJSON(`
				{
  				"error": "An application error occurred. Please contact your CredHub administrator."
				}
			`))
		})
	})

	Context("When the api returns an array", func() {
		It("can be parsed successfully", func() {
			//language=json
			responseJSON := `
        [
          {
            "id": "2993f622-cb1e-4e00-a267-4b23c273bf3d",
            "name": "/example-password",
            "type": "password",
            "value": "6mRPZB3bAfb8lRpacnXsHfDhlPqFcjH2h9YDvLpL",
            "version_created_at": "2017-01-05T01:01:01Z"
          }
        ]
      `
			server.RouteToHandler("GET", "/api/v1/data/valid-credential-id",
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data/valid-credential-id"),
					RespondWith(http.StatusOK, responseJSON),
				),
			)

			session := runCommand("curl", "-p", "api/v1/data/valid-credential-id")

			Eventually(session).Should(Exit(0))

			Eventually(session.Out.Contents()).Should(MatchJSON(`
				[
					{
						"id": "2993f622-cb1e-4e00-a267-4b23c273bf3d",
						"name": "/example-password",
						"type": "password",
						"value": "6mRPZB3bAfb8lRpacnXsHfDhlPqFcjH2h9YDvLpL",
						"version_created_at": "2017-01-05T01:01:01Z"
					}
				]
			`))
		})
	})

	Context("the user provides a valid path", func() {
		It("receives what the server returns", func() {
			//language=json
			responseJSON := `
        {
          "id": "2993f622-cb1e-4e00-a267-4b23c273bf3d",
          "name": "/example-password",
          "type": "password",
          "value": "6mRPZB3bAfb8lRpacnXsHfDhlPqFcjH2h9YDvLpL",
          "version_created_at": "2017-01-05T01:01:01Z"
        }
      `
			server.RouteToHandler("GET", "/api/v1/data/valid-credential-id",
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data/valid-credential-id"),
					RespondWith(http.StatusOK, responseJSON),
				),
			)

			session := runCommand("curl", "-p", "api/v1/data/valid-credential-id")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out.Contents()).Should(MatchJSON(`
				{
					"id": "2993f622-cb1e-4e00-a267-4b23c273bf3d",
					"name": "/example-password",
					"type": "password",
					"value": "6mRPZB3bAfb8lRpacnXsHfDhlPqFcjH2h9YDvLpL",
					"version_created_at": "2017-01-05T01:01:01Z"
				}
			`))
		})

		Context("the user does not specify required parameters", func() {
			It("returns a wrapped error", func() {
				//language=json
				responseJSON := `
          {
            "error": "The query parameter name is required for this request."
          }
        `
				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data"),
						RespondWith(http.StatusBadRequest, responseJSON),
					),
				)

				session := runCommand("curl", "-p", "api/v1/data")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out.Contents()).Should(MatchJSON(`
					{
						"error": "The query parameter name is required for this request."
					}
				`))
			})
		})

		Context("when parameters are provided by the user", func() {
			It("returns what the server returns", func() {
				//language=json
				responseJSON := `
          {
            "data": [
              {
                "id": "some-id",
                "name": "example-password",
                "type": "password",
                "value": "secret",
                "version_created_at": "time"
              }
            ]
          }
        `
				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data"),
						VerifyForm(url.Values{"name": []string{"/example-password"},
							"current": []string{"true"}}),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("curl", "-p", "api/v1/data?name=/example-password&current=true")

				Eventually(session).Should(Exit(0))
				Eventually(string(session.Out.Contents())).Should(MatchJSON(`
					{
						"data": [
							{
								"id": "some-id",
								"name": "example-password",
								"type": "password",
								"value": "secret",
								"version_created_at": "time"
							}
						]
					}
				`))
			})
		})

		Context("the user specifies a method with -X", func() {
			It("returns what the server returns", func() {
				//language=json
				responseJSON := `
          {
            "error": "The request could not be fulfilled because the request path or body did not meet expectation. Please check the documentation for required formatting and retry your request."
          }
        `
				server.RouteToHandler("PUT", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("PUT", "/api/v1/data"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("curl", "-X", "PUT", "-p", "api/v1/data")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out.Contents()).Should(MatchJSON(`
					{
						"error": "The request could not be fulfilled because the request path or body did not meet expectation. Please check the documentation for required formatting and retry your request."
					}
				`))
			})
		})

		Context("the user provides a request body", func() {
			It("receives what the server returns", func() {
				//language=json
				responseJSON := `
          {
            "type": "password",
            "version_created_at": "2018-03-06T09:10:18Z",
            "id": "93959091-0fcd-4a2a-bedb-97d3ee0d0e87",
            "name": "/some/cred",
            "value": "XbD5KGiLB4pBi24WEYq857psfvMMww"
          }
        `
				body := `{"name":"/some/cred","type":"password"}`

				server.RouteToHandler("PUT", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("PUT", "/api/v1/data"),
						VerifyBody([]byte(body)),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("curl", "-p", "api/v1/data", "-X", "PUT", "-d", body)

				Eventually(session).Should(Exit(0))
				Eventually(session.Out.Contents()).Should(MatchJSON(`
					{
						"id": "93959091-0fcd-4a2a-bedb-97d3ee0d0e87",
						"name": "/some/cred",
						"type": "password",
						"value": "XbD5KGiLB4pBi24WEYq857psfvMMww",
						"version_created_at": "2018-03-06T09:10:18Z"
					}
				`))
			})
		})

		Context("the user provides a -i flag", func() {
			It("displays the response headers to the user", func() {
				responseJSON := `{"id":"2993f622-cb1e-4e00-a267-4b23c273bf3d","name":"/example-password","type":"password","value":"6mRPZB3bAfb8lRpacnXsHfDhlPqFcjH2h9YDvLpL","version_created_at":"2017-01-05T01:01:01Z"}`
				server.RouteToHandler("GET", "/api/v1/data/valid-credential-id",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data/valid-credential-id"),
						RespondWith(http.StatusOK, responseJSON, http.Header{"Test1": []string{"test1"}, "Test2": []string{"test2"}}),
					),
				)

				session := runCommand("curl", "-p", "api/v1/data/valid-credential-id", "-i")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("HTTP/1.1 200"))
				Eventually(session.Out).Should(Say("Test1: test1"))
				Eventually(session.Out).Should(Say("Test2: test2"))
				Eventually(session.Out).Should(Say(`{
  "id": "2993f622-cb1e-4e00-a267-4b23c273bf3d",
  "name": "/example-password",
  "type": "password",
  "value": "6mRPZB3bAfb8lRpacnXsHfDhlPqFcjH2h9YDvLpL",
  "version_created_at": "2017-01-05T01:01:01Z"
}`))
			})
		})

		Context("and the user provides a -f flag", func() {
			Context("when the request fails", func() {
				It("the request fails silently", func() {
					//language=json
					responseJSON := `
       				{
					"error": "The request could not be completed because the credential does not exist or you do not have sufficient authorization."
        			}
      				`
					server.RouteToHandler("GET", "/api/v1/data",
						CombineHandlers(
							VerifyRequest("GET", "/api/v1/data"),
							VerifyForm(url.Values{"name": []string{"/doesNotExist"}}),
							RespondWith(http.StatusNotFound, responseJSON),
						),
					)

					session := runCommand("curl", "-f", "-p", "api/v1/data?name=/doesNotExist")

					Eventually(session).Should(Exit(22))
					Eventually(session.Out.Contents()).Should(BeEmpty())
				})
				Context("when request includes -i flag", func() {
					It("displays response headers to the user and exits with proper code", func() {
						//language=json
						responseJSON := `
						{
						"error": "The request could not be completed because the credential does not exist or you do not have sufficient authorization."
						}
						`
						server.RouteToHandler("GET", "/api/v1/data",
							CombineHandlers(
								VerifyRequest("GET", "/api/v1/data"),
								VerifyForm(url.Values{"name": []string{"/doesNotExist"}}),
								RespondWith(http.StatusNotFound, responseJSON),
							),
						)

						session := runCommand("curl", "-f", "-p", "api/v1/data?name=/doesNotExist", "-i")

						Eventually(session).Should(Exit(22))
						Eventually(session.Out).Should(Say("HTTP/1.1 404"))
					})
				})
			})

			It("the request continues to succeed", func() {
				responseJSON := `
				  {
					"type": "password",
					"version_created_at": "2018-03-06T09:10:18Z",
					"id": "93959091-0fcd-4a2a-bedb-97d3ee0d0e87",
					"name": "/some/cred",
					"value": "XbD5KGiLB4pBi24WEYq857psfvMMww"
				  }
				`
				body := `{"name":"/some/cred","type":"password"}`
				server.RouteToHandler("PUT", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("PUT", "/api/v1/data"),
						VerifyBody([]byte(body)),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("curl", "-f", "-p", "api/v1/data", "-X", "PUT", "-d", body)

				Eventually(session).Should(Exit(0))
			})
		})
	})
})
