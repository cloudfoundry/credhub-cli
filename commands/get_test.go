package commands_test

import (
	"bytes"
	"net/http"

	"runtime"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Get", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("get", "-n", "test-credential")
	ItRequiresAnAPIToBeSet("get", "-n", "test-credential")
	testAutoLogin := []TestAutoLogin{
		{
			method:              "GET",
			responseFixtureFile: "get_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v1/data",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogin, "get", "-n", "test-credential")

	ItBehavesLikeHelp("get", "g", func(session *Session) {
		Expect(session.Err).To(Say("Usage"))
		if runtime.GOOS == "windows" {
			Expect(session.Err).To(Say("credhub-cli.exe \\[OPTIONS\\] get \\[get-OPTIONS\\]"))
		} else {
			Expect(session.Err).To(Say("credhub-cli \\[OPTIONS\\] get \\[get-OPTIONS\\]"))
		}
	})

	It("displays missing required parameter", func() {
		session := runCommand("get")

		Eventually(session).Should(Exit(1))

		if runtime.GOOS == "windows" {
			Expect(session.Err).To(Say("A name or ID must be provided. Please update and retry your request."))
		} else {
			Expect(session.Err).To(Say("A name or ID must be provided. Please update and retry your request."))
		}
	})

	Describe("value type", func() {
		It("gets a value secret", func() {
			responseJSON := fmt.Sprintf(arrayResponseJSON, "value", "my-value", `"potatoes"`, `{}`)

			server.RouteToHandler("GET", "/api/v1/data",
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=my-value"),
					RespondWith(http.StatusOK, responseJSON),
				),
			)

			session := runCommand("get", "-n", "my-value")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: my-value"))
			Eventually(session.Out).Should(Say("type: value"))
			Eventually(session.Out).Should(Say("value: potatoes"))
		})

		Context("with --quiet flag", func() {
			It("returns only the value", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "value", "my-value", `"potatoes"`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=my-value"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-value", "-q")

				Eventually(session).Should(Exit(0))
				contents := string(bytes.TrimSpace(session.Out.Contents()))
				Eventually(contents).Should(Equal("potatoes"))
			})
		})

		Context("multiple versions with --quiet flag", func() {
			It("returns array of values", func() {
				responseJSON := fmt.Sprintf(multipleCredentialArrayResponseJSON, "value", "my-cred", `"potatoes"`, `{}`, "value", "my-cred", `"tomatoes"`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "name=my-cred&versions=2"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-cred", "-q", "--versions", "2")

				Eventually(session).Should(Exit(0))
				contents := string(bytes.TrimSpace(session.Out.Contents()))
				Eventually(contents).Should(Equal(`versions:
- potatoes
- tomatoes`))
			})
		})

		Context("--quiet flag with multi-line value", func() {
			It("should not return the value with yaml formatting", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "value", "my-value", "\"potatoes\\nand\\ntomatoes\"", `{}`)
				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=my-value"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-value", "-q")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say(`potatoes
and
tomatoes`))
			})
		})

	})

	Describe("password type", func() {
		It("gets a password secret", func() {
			responseJSON := fmt.Sprintf(arrayResponseJSON, "password", "my-password", `"potatoes"`, `{}`)

			server.RouteToHandler("GET", "/api/v1/data",
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=my-password"),
					RespondWith(http.StatusOK, responseJSON),
				),
			)

			session := runCommand("get", "-n", "my-password")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: my-password"))
			Eventually(session.Out).Should(Say("type: password"))
			Eventually(session.Out).Should(Say("value: potatoes"))
		})

		It("gets a secret by ID", func() {
			responseJSON := fmt.Sprintf(defaultResponseJSON, "password", "my-password", `"potatoes"`, `{}`)

			server.RouteToHandler("GET", "/api/v1/data/"+uuid,
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data/"+uuid),
					RespondWith(http.StatusOK, responseJSON),
				),
			)

			session := runCommand("get", "--id", uuid)

			Eventually(session).Should(Exit(0))
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: my-password"))
			Eventually(session.Out).Should(Say("type: password"))
			Eventually(session.Out).Should(Say("value: potatoes"))

		})

		Context("with key and version", func() {
			It("returns an error", func() {
				session := runCommand("get", "-n", "my-password", "--versions", "2", "-k", "someflag")
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("The --version flag and --key flag are incompatible"))
			})
		})

		Context("with --quiet flag", func() {
			It("can quiet output for password", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "password", "my-password", `"potatoes"`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=my-password"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-password", "-q")

				Eventually(session).Should(Exit(0))
				contents := string(bytes.TrimSpace(session.Out.Contents()))
				Eventually(contents).Should(Equal("potatoes"))
			})
		})

		Context("with --versions flag", func() {
			It("gets the specified number of versions of a secret", func() {
				responseJSON := `{"data":[{"type":"password","id":"` + uuid + `","name":"my-password","version_created_at":"` + timestamp + `","value":"old-password"},{"type":"password","id":"` + uuid + `","name":"my-password","version_created_at":"` + timestamp + `","value":"new-password"}]}`

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "name=my-password&versions=2"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-password", "--versions", "2")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say(`versions:
- id: ` + uuid + `
  name: my-password
  type: password
  value: old-password
  version_created_at: "` + timestamp + `"
- id: ` + uuid + `
  name: my-password
  type: password
  value: new-password
  version_created_at: "` + timestamp + `"
`))

			})
		})

		Context("multiple versions with --quiet flag", func() {
			It("returns an error", func() {
				responseJSON := `{"data":[{"type":"password","id":"` + uuid + `","name":"my-password","version_created_at":"` + timestamp + `","value":"new-password"},{"type":"password","id":"` + uuid + `","name":"my-password","version_created_at":"` + timestamp + `","value":"old-password"}]}`

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "name=my-password&versions=2"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-password", "--versions", "2", "-q")

				Eventually(session).Should(Exit(0))
				contents := string(bytes.TrimSpace(session.Out.Contents()))
				Eventually(contents).Should(Equal(`versions:
- new-password
- old-password`))
			})
		})

	})

	Describe("json type", func() {
		It("gets a json secret", func() {
			serverResponse := fmt.Sprintf(arrayResponseJSON, "json", "json-secret", `{"foo":"bar","nested":{"a":1},"an":["array"]}`, `{}`)

			server.RouteToHandler("GET", "/api/v1/data",
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=json-secret"),
					RespondWith(http.StatusOK, serverResponse),
				),
			)

			session := runCommand("get", "-n", "json-secret")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: json-secret"))
			Eventually(session.Out).Should(Say("type: json"))
			Eventually(session.Out).Should(Say(`value:
  an:
  - array
  foo: bar
  nested:
    a: 1`))

		})

		Context("with --output-json flag", func() {
			It("can output json", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "password", "my-password", `"potatoes"`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=my-password"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-password", "--output-json")

				Eventually(session).Should(Exit(0))
				Expect(string(session.Out.Contents())).To(MatchJSON(fmt.Sprintf(defaultResponseJSON, "password", "my-password", `"potatoes"`, `{}`)))
			})
		})

		Context("with --output-json and --quiet flags", func() {
			It("should return an error", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "password", "my-password", `"potatoes"`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=my-password"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-password", "--output-json", "-q")

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("The --output-json flag and --quiet flag are incompatible"))
			})
		})

		Context("with --key flag", func() {
			It("returns only the requested JSON field from the value object", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "json", "json-secret", `{"foo":"bar","nested":{"a":1},"an":["array"]}`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=json-secret"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "json-secret", "-k", "nested")

				Eventually(session).Should(Exit(0))
				Eventually(string(session.Out.Contents())).Should(Equal(`a: 1

`))
			})
		})

		Context("with --quiet flag", func() {
			It("only return the value", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "json", "json-secret", `{"foo":"bar","nested":{"a":1},"an":["array"]}`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=json-secret"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "json-secret", "-q")

				Eventually(session).Should(Exit(0))
				contents := string(bytes.TrimSpace(session.Out.Contents()))
				Eventually(contents).Should(Equal(`an:
- array
foo: bar
nested:
  a: 1`))
			})
		})

		Context("multiple versions with --quiet flag", func() {
			It("returns an array of values", func() {
				responseJSON := `{"data":[{"type":"json","id":"` + uuid + `","name":"my-json","version_created_at":"` + timestamp + `","value":{"secret":"newSecret"}},{"type":"json","id":"` + uuid + `","name":"my-json","version_created_at":"` + timestamp + `","value":{"secret":"oldSecret"}}]}`
				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "name=my-json&versions=2"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-json", "-q", "--versions", "2")

				Eventually(session).Should(Exit(0))
				contents := string(bytes.TrimSpace(session.Out.Contents()))
				Eventually(contents).Should(Equal(`versions:
- secret: newSecret
- secret: oldSecret`))
			})
		})
	})

	Describe("certificate type", func() {
		It("gets a certificate secret", func() {
			responseJSON := fmt.Sprintf(arrayResponseJSON, "certificate", "my-secret", `{"ca":"my-ca","certificate":"my-cert","private_key":"my-priv"}`, `{}`)

			server.RouteToHandler("GET", "/api/v1/data",
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=my-secret"),
					RespondWith(http.StatusOK, responseJSON),
				),
			)

			session := runCommand("get", "-n", "my-secret")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: my-secret"))
			Eventually(session.Out).Should(Say("type: certificate"))
			Eventually(session.Out).Should(Say("ca: my-ca"))
			Eventually(session.Out).Should(Say("certificate: my-cert"))
			Eventually(session.Out).Should(Say("private_key: my-priv"))
		})

		Context("with --key flag", func() {
			It("only returns the request field from the value object", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "certificate", "my-secret", `{"ca":"----begin----my-ca-----end------","certificate":"my-cert","private_key":"my-priv"}`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=my-secret"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-secret", "-k", "ca")

				Eventually(session).Should(Exit(0))
				Eventually(string(session.Out.Contents())).Should(Equal("----begin----my-ca-----end------\n"))
			})
		})

		Context("with invalid key", func() {
			It("returns nothing", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "certificate", "my-secret", `{"ca":"my-ca","certificate":"my-cert","private_key":"my-priv"}`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=my-secret"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-secret", "-k", "invalidkey")

				Eventually(session).Should(Exit(0))
				Eventually(string(session.Out.Contents())).Should(Equal(``))

			})
		})

		Context("with --quiet flag", func() {
			It("only returns the value", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "certificate", "my-secret", `{"ca":"----begin----my-ca-----end------","certificate":"----begin----my-cert-----end------","private_key":"----begin----my-priv-----end------"}`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=my-secret"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-secret", "-q")

				Eventually(session).Should(Exit(0))
				Eventually(string(session.Out.Contents())).Should(ContainSubstring(`ca: '----begin----my-ca-----end------'
certificate: '----begin----my-cert-----end------'
private_key: '----begin----my-priv-----end------'`))
			})
		})

		Context("multiple versions with --quiet flag", func() {
			It("only returns the value", func() {
				responseJSON := fmt.Sprintf(multipleCredentialArrayResponseJSON,
					"certificate",
					"my-secret",
					`{"ca":"----begin----my-new-ca-----end------","certificate":"----begin----my-new-cert-----end------","private_key":"----begin----my-new-priv-----end------"}`,
					"{}",
					"certificate",
					"my-secret",
					`{"ca":"----begin----my-old-ca-----end------","certificate":"----begin----my-old-cert-----end------","private_key":"----begin----my-old-priv-----end------"}`,
					"{}")
				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "name=my-secret&versions=2"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-secret", "-q", "--versions", "2")

				Eventually(session).Should(Exit(0))
				Eventually(string(bytes.TrimSpace(session.Out.Contents()))).Should(Equal(`versions:
- ca: '----begin----my-new-ca-----end------'
  certificate: '----begin----my-new-cert-----end------'
  private_key: '----begin----my-new-priv-----end------'
- ca: '----begin----my-old-ca-----end------'
  certificate: '----begin----my-old-cert-----end------'
  private_key: '----begin----my-old-priv-----end------'`))
			})
		})

		Context("--quiet flag with key", func() {
			It("should not only return the value", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "certificate", "my-secret", `{"ca":"----begin----my-ca-----end------","certificate":"----begin----my-cert-----end------","private_key":"----begin----my-priv-----end------"}`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=my-secret"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-secret", "-q", "-k", "ca")

				Eventually(session).Should(Exit(0))
				Eventually(string(session.Out.Contents())).ShouldNot(Equal(`ca: |-
  ----begin----my-ca-----end------
certificate: |-
  ----begin----my-cert-----end------
private_key: |-
  ----begin----my-priv-----end------

`))
			})
		})

	})

	Describe("rsa type", func() {
		It("gets an rsa secret", func() {
			responseJSON := fmt.Sprintf(arrayResponseJSON, "rsa", "foo-rsa-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{}`)

			server.RouteToHandler("GET", "/api/v1/data",
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=foo-rsa-key"),
					RespondWith(http.StatusOK, responseJSON),
				),
			)

			session := runCommand("get", "-n", "foo-rsa-key")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: foo-rsa-key"))
			Eventually(session.Out).Should(Say("type: rsa"))
			Eventually(session.Out).Should(Say("private_key: some-private-key"))
			Eventually(session.Out).Should(Say("public_key: some-public-key"))
		})

		Context("with --quiet flag", func() {
			It("gets only the value", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "rsa", "foo-rsa-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=foo-rsa-key"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "foo-rsa-key", "-q")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).ShouldNot(Say("name: foo-rsa-key"))
				Eventually(session.Out).ShouldNot(Say("type: rsa"))
				Eventually(session.Out).Should(Say("private_key: some-private-key"))
				Eventually(session.Out).Should(Say("public_key: some-public-key"))
			})
		})

		Context("multiple versions with --quiet flag", func() {
			It("only returns the value", func() {
				responseJSON := fmt.Sprintf(multipleCredentialArrayResponseJSON,
					"rsa",
					"foo-rsa-key",
					`{"public_key":"new-public-key","private_key":"new-private-key"}`,
					`{}`,
					"rsa",
					"foo-rsa-key",
					`{"public_key":"old-public-key","private_key":"old-private-key"}`,
					`{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "versions=2&name=foo-rsa-key"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "foo-rsa-key", "-q", "--versions", "2")

				Eventually(session).Should(Exit(0))
				contents := string(bytes.TrimSpace(session.Out.Contents()))
				Eventually(contents).Should(Equal(`versions:
- private_key: new-private-key
  public_key: new-public-key
- private_key: old-private-key
  public_key: old-public-key`))
			})
		})

		Context("--quiet flag with key", func() {
			It("should not only return the value", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "rsa", "foo-rsa-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=foo-rsa-key"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "foo-rsa-key", "-q", "-k", "public_key")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).ShouldNot(Say("name: foo-rsa-key"))
				Eventually(session.Out).ShouldNot(Say("type: rsa"))
				Eventually(session.Out).Should(Say("some-public-key"))
			})
		})
	})

	Describe("user type", func() {
		It("gets a user secret", func() {
			responseJSON := fmt.Sprintf(arrayResponseJSON, "user", "my-username-credential", `{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4$h"}`, `{}`)

			server.RouteToHandler("GET", "/api/v1/data",
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=my-username-credential"),
					RespondWith(http.StatusOK, responseJSON),
				),
			)

			session := runCommand("get", "-n", "my-username-credential")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: my-username-credential"))
			Eventually(session.Out).Should(Say("type: user"))
			Eventually(session.Out).Should(Say("password: test-password"))
			Eventually(session.Out).Should(Say(`password_hash: passw0rd-H4\$h`))
			Eventually(session.Out).Should(Say("username: my-username"))
		})

		Context("with --quiet flag", func() {
			It("gets only the value", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "user", "my-username-credential", `{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4$h"}`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=my-username"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-username", "-q")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).ShouldNot(Say("name: my-username-credential"))
				Eventually(session.Out).ShouldNot(Say("type: user"))
				Eventually(session.Out).Should(Say("password: test-password"))
				Eventually(session.Out).Should(Say(`password_hash: passw0rd-H4\$h`))
				Eventually(session.Out).Should(Say("username: my-username"))
			})
		})

		Context("multiple versions with the --quiet flag", func() {
			It("returns an array of values", func() {
				responseJSON := fmt.Sprintf(multipleCredentialArrayResponseJSON,
					"user",
					"my-username-credential",
					`{"username":"new-username", "password":"new-password", "password_hash":"new-passw0rd-H4$h"}`,
					`{}`,
					"user",
					"my-username-credential",
					`{"username":"old-username", "password":"old-password", "password_hash":"old-passw0rd-H4$h"}`,
					`{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "name=my-username-credential&versions=2"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-username-credential", "-q", "--versions", "2")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).ShouldNot(Say("name: my-username-credential"))
				Eventually(session.Out).ShouldNot(Say("type: user"))
				Eventually(session.Out).Should(Say("versions:"))
				Eventually(session.Out).Should(Say("- password: new-password"))
				Eventually(session.Out).Should(Say(`  password_hash: new-passw0rd-H4\$h`))
				Eventually(session.Out).Should(Say("  username: new-username"))
				Eventually(session.Out).Should(Say("- password: old-password"))
				Eventually(session.Out).Should(Say(`  password_hash: old-passw0rd-H4\$h`))
				Eventually(session.Out).Should(Say("  username: old-username"))
			})
		})

		Context("--quiet flag with key", func() {
			It("ignores the quiet flag", func() {
				responseJSON := fmt.Sprintf(arrayResponseJSON, "user", "my-username-credential", `{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4$h"}`, `{}`)

				server.RouteToHandler("GET", "/api/v1/data",
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=my-username"),
						RespondWith(http.StatusOK, responseJSON),
					),
				)

				session := runCommand("get", "-n", "my-username", "-q", "-k", "password_hash")

				Eventually(session).Should(Exit(0))
				contents := string(bytes.TrimSpace(session.Out.Contents()))
				Eventually(contents).Should(Equal(`passw0rd-H4$h`))
			})

		})
	})

	It("does not use Printf on user-supplied data", func() {
		responseJSON := fmt.Sprintf(arrayResponseJSON, "password", "injected", `"et''%/7(V&|?m|Ckih$"`, `{}`)

		server.RouteToHandler("GET", "/api/v1/data",
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "current=true&name=injected"),
				RespondWith(http.StatusOK, responseJSON),
			),
		)

		session := runCommand("get", "-n", "injected")

		Eventually(session).Should(Exit(0))
		outStr := "et''%/7\\(V&|\\?m\\|Ckih\\$"
		Eventually(session.Out).Should(Say(outStr + timestamp))
	})
})
