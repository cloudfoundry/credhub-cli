package commands_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var (
	session      *gexec.Session
	templateFile *os.File
	templateText string
	err          error
)

var _ = Describe("interpolate", func() {
	BeforeEach(func() {
		login()
		templateFile, err = ioutil.TempFile("", "credhub_test_interpolate_template_")
	})

	Describe("behavior shared with other commands", func() {
		templateFile, err = ioutil.TempFile("", "credhub_test_interpolate_template_")
		templateFile.WriteString("---")
		testAutoLogin := []TestAutoLogin{
			{
				method:              "GET",
				responseFixtureFile: "get_response.json",
				responseStatus:      http.StatusOK,
				endpoint:            "/api/v1/data",
			},
		}
		ItAutomaticallyLogsIn(testAutoLogin, "interpolate", "-f", templateFile.Name())

		ItBehavesLikeHelp("interpolate", "interpolate", func(session *gexec.Session) {
			Expect(session.Err).To(Say("Usage"))
			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("credhub-cli.exe \\[OPTIONS\\] interpolate \\[interpolate-OPTIONS\\]"))
			} else {
				Expect(session.Err).To(Say("credhub-cli \\[OPTIONS\\] interpolate \\[interpolate-OPTIONS\\]"))
			}
		})
		ItRequiresAuthentication("interpolate", "-f", "testinterpolationtemplate.yml")
		ItRequiresAnAPIToBeSet("interpolate", "-f", "testinterpolationtemplate.yml")
	})

	Describe("interpolating various types of credentials", func() {
		It("queries for string creds and prints them in the template as strings", func() {
			templateText = `---
value-cred: ((relative/value/cred/path))
static-value: a normal string`
			templateFile.WriteString(templateText)
			responseValueJson := fmt.Sprintf(STRING_CREDENTIAL_ARRAY_RESPONSE_JSON, "value", "relative/value/cred/path", `{\"value\": \"should not be interpolated\"}`)

			credentialListJson, err := credentialsListJSON([]string{"/relative/value/cred/path"})
			Expect(err).Should(BeNil())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, credentialListJson),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/relative/value/cred/path"),
					RespondWith(http.StatusOK, responseValueJson),
				),
			)

			session = runCommand("interpolate", "-f", templateFile.Name())
			Eventually(session).Should(gexec.Exit(0))
			Expect(string(session.Out.Contents())).To(MatchYAML(`
value-cred: "{\"value\": \"should not be interpolated\"}"
static-value: a normal string
`))
		})

		It("queries for multi-line, multi-part credential types and prints them in the template", func() {
			templateText = `---
full-certificate-cred: ((relative/certificate/cred/path))
cert-only-certificate-cred: ((relative/certificate/cred/path.certificate))
static-value: a normal string`
			templateFile.WriteString(templateText)

			responseCertJson := fmt.Sprintf(CERTIFICATE_CREDENTIAL_ARRAY_RESPONSE_JSON, "test-cert", "", "-----BEGIN FAKE CERTIFICATE-----\\n-----END FAKE CERTIFICATE-----", "-----BEGIN FAKE RSA PRIVATE KEY-----\\n-----END FAKE RSA PRIVATE KEY-----")

			credentialListJson, err := credentialsListJSON([]string{"/relative/certificate/cred/path"})
			Expect(err).Should(BeNil())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, credentialListJson),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/relative/certificate/cred/path"),
					RespondWith(http.StatusOK, responseCertJson),
				),
			)

			session = runCommand("interpolate", "-f", templateFile.Name())
			Eventually(session).Should(gexec.Exit(0))
			Expect(string(session.Out.Contents())).To(MatchYAML(`
full-certificate-cred:
  ca: ""
  certificate: |-
    -----BEGIN FAKE CERTIFICATE-----
    -----END FAKE CERTIFICATE-----
  private_key: |-
    -----BEGIN FAKE RSA PRIVATE KEY-----
    -----END FAKE RSA PRIVATE KEY-----
cert-only-certificate-cred: |-
  -----BEGIN FAKE CERTIFICATE-----
  -----END FAKE CERTIFICATE-----
static-value: a normal string
`))
		})

		It("queries for json creds and prints them in the template rendered as yaml", func() {
			templateText = `json-cred: ((relative/json/cred/path))`
			templateFile.WriteString(templateText)

			responseJson := fmt.Sprintf(JSON_CREDENTIAL_ARRAY_RESPONSE_JSON, "test-json", `{"whatthing":"something"}`)
			credentialListJson, err := credentialsListJSON([]string{"/relative/json/cred/path"})
			Expect(err).Should(BeNil())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, credentialListJson),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/relative/json/cred/path"),
					RespondWith(http.StatusOK, responseJson),
				),
			)

			session = runCommand("interpolate", "-f", templateFile.Name())
			Eventually(session).Should(gexec.Exit(0))
			Expect(string(session.Out.Contents())).To(MatchYAML(`json-cred: {"whatthing":"something"}`))
		})
	})

	Describe("the optional --prefix flag", func() {
		BeforeEach(func() {
			templateText = `---
full-certificate-cred: ((certificate/cred/path))
cert-only-certificate-cred: ((/relative/certificate/cred/path.certificate))
static-value: a normal string`
			templateFile.WriteString(templateText)

			responseCertJson := fmt.Sprintf(CERTIFICATE_CREDENTIAL_ARRAY_RESPONSE_JSON, "test-cert", "", "-----BEGIN FAKE CERTIFICATE-----\\n-----END FAKE CERTIFICATE-----", "-----BEGIN FAKE RSA PRIVATE KEY-----\\n-----END FAKE RSA PRIVATE KEY-----")
			credentialListJson, err := credentialsListJSON([]string{"/relative/certificate/cred/path"})
			Expect(err).Should(BeNil())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, credentialListJson),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/relative/certificate/cred/path"),
					RespondWith(http.StatusOK, responseCertJson),
				),
			)
		})
		It("prints the values of credential names derived from the prefix, unless the cred paths start with /", func() {
			session = runCommand("interpolate", "-f", templateFile.Name(), "-p=/relative")
			Eventually(session).Should(gexec.Exit(0))
			Expect(string(session.Out.Contents())).To(MatchYAML(`
full-certificate-cred:
  ca: ""
  certificate: |-
    -----BEGIN FAKE CERTIFICATE-----
    -----END FAKE CERTIFICATE-----
  private_key: |-
    -----BEGIN FAKE RSA PRIVATE KEY-----
    -----END FAKE RSA PRIVATE KEY-----
cert-only-certificate-cred: |-
  -----BEGIN FAKE CERTIFICATE-----
  -----END FAKE CERTIFICATE-----
static-value: a normal string
`))
		})
	})

	Describe("when template has different paths than prefix", func() {
		var credentialListJson string
		var err error
		BeforeEach(func() {
			credentialListJson, err = credentialsListJSON([]string{"/a/pass1", "/a/myval", "/b/pass2", "/b/pass"})
			Expect(err).Should(BeNil())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, credentialListJson),
				),
			)
		})
		It("finds credential without path", func() {
			templateFile.WriteString(`---
/a/pass1: ((pass1))`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/a/pass1"),
					RespondWith(http.StatusOK, fmt.Sprintf(STRING_CREDENTIAL_ARRAY_RESPONSE_JSON, "value", "a/pass1", "pass1")),
				),
			)
			session = runCommand("interpolate", "-f", templateFile.Name(), "-p=/a")
			Eventually(session).Should(gexec.Exit(0))
			Expect(string(session.Out.Contents())).To(MatchYAML(`
/a/pass1: pass1
`))
		})

		It("finds credential in different path", func() {
			templateFile.WriteString(`---
/b/pass: ((/b/pass))`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/b/pass"),
					RespondWith(http.StatusOK, fmt.Sprintf(STRING_CREDENTIAL_ARRAY_RESPONSE_JSON, "value", "b/pass", "pass")),
				),
			)
			session = runCommand("interpolate", "-f", templateFile.Name(), "-p=/a")
			Eventually(session).Should(gexec.Exit(0))
			Expect(string(session.Out.Contents())).To(MatchYAML(`
/b/pass: pass
`))
		})
	})
	Describe("Errors", func() {
		Context("when no template file is provided", func() {
			BeforeEach(func() {
				session = runCommand("interpolate")
			})
			It("prints missing required parameter", func() {
				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err).To(Say("A file to interpolate must be provided. Please add a file flag and try again."))
			})
		})

		Context("when the --file doesn't exist", func() {
			var invalidFile = "does-not-exist.yml"

			It("prints an error that includes the filepath and the filesystem error", func() {
				session := runCommand("interpolate", "-f", invalidFile)
				Eventually(session).Should(gexec.Exit(1), "interpolate should have failed")
				Expect(session.Err).To(Say(invalidFile))
			})
		})

		Context("when the template file contains no credentials to resolve", func() {
			BeforeEach(func() {
				templateText = `---
yaml-key-with-static-value: a normal string`
				templateFile.WriteString(templateText)
				credentialListJson, err := credentialsListJSON([]string{""})
				Expect(err).Should(BeNil())
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "path="),
						RespondWith(http.StatusOK, credentialListJson),
					),
				)
			})
			It("succeeds and prints the template to stdout", func() {
				session := runCommand("interpolate", "-f", templateFile.Name())
				Eventually(session).Should(gexec.Exit(0), "command should succeed")
				Expect(session.Out).To(Say("yaml-key-with-static-value: a normal string"))
			})
		})

		Context("when a path in the --file can't be found", func() {
			BeforeEach(func() {
				templateText = `---
yaml-key-with-template-value: ((relative/cred/path))
yaml-key-with-static-value: a normal string`
				templateFile.WriteString(templateText)

				credentialListJson, err := credentialsListJSON([]string{"/relative/cred/path"})
				Expect(err).Should(BeNil())
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "path="),
						RespondWith(http.StatusOK, credentialListJson),
					),
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "current=true&name=/relative/cred/path"),
						RespondWith(http.StatusOK, `{"data":[]}`),
					),
				)
			})

			It("prints an error that includes the credential path and the underlying error", func() {
				session = runCommand("interpolate", "-f", templateFile.Name())
				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err).To(Say("Finding variable 'relative/cred/path': response did not contain any credentials"))
			})
		})

		Context("when skip is specified", func() {
			BeforeEach(func() {
				templateText = `---
yaml-key-with-template-value: ((not_a_cred))
yaml-key-with-static-value: a normal string`
				templateFile.WriteString(templateText)

				credentialListJson, err := credentialsListJSON([]string{""})
				Expect(err).Should(BeNil())
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "path="),
						RespondWith(http.StatusOK, credentialListJson),
					),
				)
			})

			It("succeeds and prints out same yaml", func() {
				session = runCommand("interpolate", "-f", templateFile.Name(), "-s")
				Eventually(session).Should(gexec.Exit(0))
				Expect(string(session.Out.Contents())).To(MatchYAML(`
yaml-key-with-template-value: ((not_a_cred))
yaml-key-with-static-value: a normal string`))
			})
		})

		Context("when the file has invalid yaml", func() {
			It("prints an error that says there is invalid yaml", func() {
				templateText = `key: - value`
				templateFile.WriteString(templateText)

				credentialListJson, err := credentialsListJSON([]string{""})
				Expect(err).Should(BeNil())
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "path="),
						RespondWith(http.StatusOK, credentialListJson),
					),
				)

				session = runCommand("interpolate", "-f", templateFile.Name())
				Eventually(session).Should(gexec.Exit(1))
			})
		})
	})

	Describe("Empty file", func() {
		Context("when the template file is empty", func() {
			It("does not throw an error", func() {
				session := runCommand("interpolate", "-f", templateFile.Name())
				Eventually(session).Should(gexec.Exit(0))
			})
		})
	})
})

func credentialsListJSON(names []string) (string, error) {
	type creds struct {
		Name string `json:"name"`
	}
	type credholder struct {
		Credentials []creds `json:"credentials"`
	}
	holder := &credholder{}
	for _, name := range names {
		holder.Credentials = append(holder.Credentials, creds{Name: name})
	}
	bytes, err := json.Marshal(holder)
	return string(bytes), err
}
