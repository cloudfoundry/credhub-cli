package commands_test

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"io/ioutil"
	"net/http"
	"os"
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
			Expect(session.Err).To(Say("Usage:\n(.*)\\[OPTIONS\\] interpolate \\[interpolate-OPTIONS\\]"))
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
			responseValueJSON := fmt.Sprintf(arrayResponseJSON, "value", "relative/value/cred/path", `"should not be interpolated"`, `{}`)

			credentialListJSON, err := credentialsListJSON([]string{"/relative/value/cred/path"})
			Expect(err).Should(BeNil())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, credentialListJSON),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/relative/value/cred/path"),
					RespondWith(http.StatusOK, responseValueJSON),
				),
			)

			session = runCommand("interpolate", "-f", templateFile.Name())
			Eventually(session).Should(gexec.Exit(0))
			Expect(string(session.Out.Contents())).To(MatchYAML(`
value-cred: should not be interpolated
static-value: a normal string
`))
		})

		It("queries for multi-line, multi-part credential types and prints them in the template", func() {
			templateText = `---
full-certificate-cred: ((relative/certificate/cred/path))
cert-only-certificate-cred: ((relative/certificate/cred/path.certificate))
static-value: a normal string`
			templateFile.WriteString(templateText)

			responseCertJSON := fmt.Sprintf(arrayResponseJSON, "certificate", "test-cert", "{\"ca\":\"\",\"certificate\":\"-----BEGIN FAKE CERTIFICATE-----\\n-----END FAKE CERTIFICATE-----\",\"private_key\":\"-----BEGIN FAKE RSA PRIVATE KEY-----\\n-----END FAKE RSA PRIVATE KEY-----\"}", `{}`)

			credentialListJSON, err := credentialsListJSON([]string{"/relative/certificate/cred/path"})
			Expect(err).Should(BeNil())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, credentialListJSON),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/relative/certificate/cred/path"),
					RespondWith(http.StatusOK, responseCertJSON),
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

			responseJSON := fmt.Sprintf(arrayResponseJSON, "json", "test-json", `{"whatthing":"something"}`, `{}`)
			credentialListJSON, err := credentialsListJSON([]string{"/relative/json/cred/path"})
			Expect(err).Should(BeNil())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, credentialListJSON),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/relative/json/cred/path"),
					RespondWith(http.StatusOK, responseJSON),
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

			responseCertJSON := fmt.Sprintf(arrayResponseJSON, "certificate", "test-cert", "{\"ca\":\"\",\"certificate\":\"-----BEGIN FAKE CERTIFICATE-----\\n-----END FAKE CERTIFICATE-----\",\"private_key\":\"-----BEGIN FAKE RSA PRIVATE KEY-----\\n-----END FAKE RSA PRIVATE KEY-----\"}", `{}`)
			credentialListJSON, err := credentialsListJSON([]string{"/relative/certificate/cred/path"})
			Expect(err).Should(BeNil())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, credentialListJSON),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/relative/certificate/cred/path"),
					RespondWith(http.StatusOK, responseCertJSON),
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
		var credentialListJSON string
		var err error
		BeforeEach(func() {
			credentialListJSON, err = credentialsListJSON([]string{"/a/pass1", "/a/myval", "/b/pass2", "/b/pass"})
			Expect(err).Should(BeNil())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, credentialListJSON),
				),
			)
		})
		It("finds credential without path", func() {
			templateFile.WriteString(`---
/a/pass1: ((pass1))`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/a/pass1"),
					RespondWith(http.StatusOK, fmt.Sprintf(arrayResponseJSON, "value", "a/pass1", `"pass1"`, `{}`)),
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
					RespondWith(http.StatusOK, fmt.Sprintf(arrayResponseJSON, "value", "b/pass", `"pass"`, `{}`)),
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
				credentialListJSON, err := credentialsListJSON([]string{""})
				Expect(err).Should(BeNil())
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "path="),
						RespondWith(http.StatusOK, credentialListJSON),
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

				credentialListJSON, err := credentialsListJSON([]string{"/relative/cred/path"})
				Expect(err).Should(BeNil())
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "path="),
						RespondWith(http.StatusOK, credentialListJSON),
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

				credentialListJSON, err := credentialsListJSON([]string{""})
				Expect(err).Should(BeNil())
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "path="),
						RespondWith(http.StatusOK, credentialListJSON),
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

				credentialListJSON, err := credentialsListJSON([]string{""})
				Expect(err).Should(BeNil())
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/api/v1/data", "path="),
						RespondWith(http.StatusOK, credentialListJSON),
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

	Describe("when parameters could not be interpolated", func() {
		It("prints a warning with the filename and line number of the uninterpolated parameters", func() {
			templateText = `---
static-value: a normal string
value-cred: ((relative/value/cred/path))
`
			templateFile.WriteString(templateText)

			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "path="),
					RespondWith(http.StatusOK, `{"credentials":[]}`),
				),
				CombineHandlers(
					VerifyRequest("GET", "/api/v1/data", "current=true&name=/relative/value/cred/path"),
					RespondWith(http.StatusNotFound, ``),
				),
			)

			session = runCommand("interpolate", "-f", templateFile.Name(), "--skip-missing")
			Eventually(session).Should(gexec.Exit(0))
			Expect(string(session.Err.Contents())).To(ContainSubstring(fmt.Sprintf(`Could not find values for:
%s:3 ((relative/value/cred/path))`, templateFile.Name())))
			Expect(string(session.Out.Contents())).To(MatchYAML(`
value-cred: ((relative/value/cred/path))
static-value: a normal string
`))
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
