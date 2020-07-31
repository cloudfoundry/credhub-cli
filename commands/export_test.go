package commands_test

import (
	"io/ioutil"
	"net/http"
	"os"

	"runtime"

	"code.cloudfoundry.org/credhub-cli/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

func withTemporaryFile(wantingFile func(string)) error {
	f, err := ioutil.TempFile("", "credhub_tests_")

	if err != nil {
		return err
	}

	name := f.Name()

	f.Close()
	wantingFile(name)

	return os.Remove(name)
}

var _ = Describe("Export", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("export")
	ItRequiresAnAPIToBeSet("export")

	testAutoLogIns := []TestAutoLogin{
		{
			method:              "GET",
			responseFixtureFile: "get_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v1/data",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogIns, "export")

	ItBehavesLikeHelp("export", "e", func(session *Session) {
		Expect(session.Err).To(Say("Usage"))
		if runtime.GOOS == "windows" {
			Expect(session.Err).To(Say("credhub-cli.exe \\[OPTIONS\\] export \\[export-OPTIONS\\]"))
		} else {
			Expect(session.Err).To(Say("credhub-cli \\[OPTIONS\\] export \\[export-OPTIONS\\]"))
		}
	})

	Describe("Exporting", func() {
		It("queries for the most recent version of all credentials", func() {
			findJSON := `{
				"credentials": [
					{
						"version_created_at": "idc",
						"name": "/path/to/cred"
					},
					{
						"version_created_at": "idc",
						"name": "/path/to/another/cred"
					}
				]
			}`

			getJSON := `{
				"data": [{
					"type":"value",
					"id":"some_uuid",
					"name":"/path/to/cred",
					"version_created_at":"idc",
					"value": "foo",
					"metadata": {
						"some": "thing"
					}
				}]
			}`

			responseTable := `credentials:
- name: /path/to/cred
  type: value
  value: foo
  metadata:
    some: thing
- name: /path/to/cred
  type: value
  value: foo
  metadata:
    some: thing`

			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, findJSON),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "name=/path/to/cred&current=true"),
					RespondWith(http.StatusOK, getJSON),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "name=/path/to/another/cred&current=true"),
					RespondWith(http.StatusOK, getJSON),
				),
			)

			session := runCommand("export")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseTable))
		})

		It("resolves ca names for certificats", func() {
			findJson := `{
				"credentials": [
					{
						"version_created_at": "idc",
						"name": "/path/to/cert"
					},
					{
						"version_created_at": "idc",
						"name": "/path/to/cert_ca"
					}
				]
			}`

			getCertJson := `{
				"data": [{
					"type": "certificate",
					"certificate_authority": false,
					"expiry_date": "some_expiry_date",
					"generated": true,
					"id": "some_uuid",
					"name": "/path/to/cert",
					"self_signed": false,
					"transitional": false,
					"version_created_at": "idc",
					"value": {
						"ca": "some_ca",
						"certificate": "some_cert",
						"private_key": "private_key"
					}
				}]}`

			getCertCaJson := `{
				"data": [{
					"type": "certificate",
					"certificate_authority": true,
					"expiry_date": "some_expiry_date",
					"generated": true,
					"id": "some_uuid",
					"name": "/path/to/cert_ca",
					"self_signed": true,
					"transitional": false,
					"version_created_at": "idc",
					"value": {
						"ca": "some_ca",
						"certificate": "some_cert",
						"private_key": "private_key"
					}
				}]}`

			getCertMetaJson := `{
				"certificates": [{
					"id": "cert-id",
					"name": "/path/to/cert",
					"signed_by": "/path/to/cert_ca",
					"signs": [],
					"versions": [{
							"certificate_authority": false,
							"expiry_date": "2020-11-28T14:04:40Z",
							"generated": true,
							"id": "cert-version-id",
							"self_signed": false,
							"transitional": false
					}]
				}]}`

			getCertCaMetaJson := `{
				"certificates":[{
					"id": "cert-ca-id",
					"name": "/path/to/cert_ca",
					"signed_by": "/path/to/cert_ca",
					"signs": [
						"/path/to/cert"
					],
					"versions": [{
							"certificate_authority": true,
							"expiry_date": "2020-11-28T14:04:38Z",
							"generated": true,
							"id": "cert-ca-version-id",
							"self_signed": true,
							"transitional": false
						}
					]
				}]}`

			responseTable := `credentials:
- name: /path/to/cert
  type: certificate
  value:
    ca_name: /path/to/cert_ca
    certificate: some_cert
    private_key: private_key
  metadata: {}
- name: /path/to/cert_ca
  type: certificate
  value:
    ca: some_ca
    certificate: some_cert
    private_key: private_key
  metadata: {}`

			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, findJson),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "name=/path/to/cert&current=true"),
					RespondWith(http.StatusOK, getCertJson),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/certificates/", "name=/path/to/cert"),
					RespondWith(http.StatusOK, getCertMetaJson),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "name=/path/to/cert_ca&current=true"),
					RespondWith(http.StatusOK, getCertCaJson),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/certificates/", "name=/path/to/cert_ca"),
					RespondWith(http.StatusOK, getCertCaMetaJson),
				),
			)

			session := runCommand("export")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say(responseTable))
		})

		Context("when given a path", func() {
			It("queries for credentials matching that path", func() {
				noCredsJSON := `{ "credentials" : [] }`

				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "path=some/path"),
						RespondWith(http.StatusOK, noCredsJSON),
					),
				)

				session := runCommand("export", "-p", "some/path")

				Eventually(session).Should(Exit(0))
			})
		})

		Context("when given a file", func() {
			Context("when given output-json flag", func() {
				It("writes the JSON to that file", func() {
					withTemporaryFile(func(filename string) {
						noCredsJSON := `{ "credentials" : [] }`
						jsonOutput := `{"Credentials":[]}`

						server.AppendHandlers(
							CombineHandlers(
								VerifyRequest("GET", "/api/v1/data", "path="),
								RespondWith(http.StatusOK, noCredsJSON),
							),
						)

						session := runCommand("export", "-f", filename, "--output-json")

						Eventually(session).Should(Exit(0))

						Expect(filename).To(BeAnExistingFile())

						fileContents, _ := ioutil.ReadFile(filename)

						Expect(string(fileContents)).To(Equal(jsonOutput))
					})
				})
			})

			Context("when not given output-json flag", func() {
				It("writes the YAML to that file", func() {
					withTemporaryFile(func(filename string) {
						noCredsJSON := `{ "credentials" : [] }`
						noCredsYaml := `credentials: []
`

						server.AppendHandlers(
							CombineHandlers(
								VerifyRequest("GET", "/api/v1/data", "path="),
								RespondWith(http.StatusOK, noCredsJSON),
							),
						)

						session := runCommand("export", "-f", filename)

						Eventually(session).Should(Exit(0))

						Expect(filename).To(BeAnExistingFile())

						fileContents, _ := ioutil.ReadFile(filename)

						Expect(string(fileContents)).To(Equal(noCredsYaml))
					})
				})
			})
		})
	})

	Describe("Errors", func() {
		It("prints an error when the network request fails", func() {
			cfg := config.ReadConfig()
			cfg.ApiURL = "mashed://potatoes"
			config.WriteConfig(cfg)

			session := runCommand("export")

			Eventually(session).Should(Exit(1))
			Eventually(string(session.Err.Contents())).Should(ContainSubstring("unsupported protocol scheme"))
		})

		It("prints an error if the specified output file cannot be opened", func() {
			noCredsJSON := `{ "credentials" : [] }`

			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, noCredsJSON),
				),
			)

			session := runCommand("export", "-f", "this/should/not/exist/anywhere")

			Eventually(session).Should(Exit(1))
			if runtime.GOOS == "windows" {
				Eventually(string(session.Err.Contents())).Should(ContainSubstring("open this/should/not/exist/anywhere: The system cannot find the path specified"))
			} else {
				Eventually(string(session.Err.Contents())).Should(ContainSubstring("open this/should/not/exist/anywhere: no such file or directory"))
			}
		})
	})
})
