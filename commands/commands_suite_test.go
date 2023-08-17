package commands_test

import (
	"crypto/tls"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"

	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"path/filepath"

	"code.cloudfoundry.org/credhub-cli/config"
	test_util "code.cloudfoundry.org/credhub-cli/test"
)

const timestamp = `2016-01-01T12:00:00Z`
const uuid = `5a2edd4f-1686-4c8d-80eb-5daa866f9f86`

const validAccessToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImxlZ2FjeS10b2tlbi1rZXkiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiI3NTY5MTc5OTgzOTY0M2Y4OWI2NGZlNDQ2MWU0OWJlMCIsInN1YiI6IjY3ODdiYjdlLTc4YmItNGJlNi05NTgzLTQyYTc1ZGRiYTNkNSIsInNjb3BlIjpbImNyZWRodWIud3JpdGUiLCJjcmVkaHViLnJlYWQiXSwiY2xpZW50X2lkIjoiY3JlZGh1Yl9jbGkiLCJjaWQiOiJjcmVkaHViX2NsaSIsImF6cCI6ImNyZWRodWJfY2xpIiwicmV2b2NhYmxlIjp0cnVlLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoiNjc4N2JiN2UtNzhiYi00YmU2LTk1ODMtNDJhNzVkZGJhM2Q1Iiwib3JpZ2luIjoidWFhIiwidXNlcl9uYW1lIjoiY3JlZGh1YiIsImVtYWlsIjoiY3JlZGh1YiIsImF1dGhfdGltZSI6MTUwNDgyMTU4NSwicmV2X3NpZyI6ImU0Yjg2ODVlIiwiaWF0IjoxNTA0ODIxNTg1LCJleHAiOjE1MDQ5MDc5ODUsImlzcyI6Imh0dHBzOi8vMzQuMjA2LjIzMy4xOTU6ODQ0My9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjcmVkaHViX2NsaSIsImNyZWRodWIiXX0.Ubi5k3Sy4CkcTqKvKuSkLJFpA5zfwWPlhImuwMW3HyKd6iEPuteXqnSE9r6ndvcKf_B3PS0ZduPg7v81RiZyfTGu3ObWIEdYExlmI97yfg4OQMCfo4jdr2wSzpcwixTK2FeZ2RcDklMfaSp_CTAnNcY4Lj2Jlk2eagWOCXizxsB1SHfegtGWH3FSUT5I3nJVcWAsRCMLqjHzRWYdP3CfpnMhnrjic1Ok_f2HKygiG44uUx2u3yQOV1CiZJwhxPODTuhI8X9kkQ0rLW9jW9ADVFstfXOglr-_k6tJMKMNpbXuCd_XaxOIXsxrSdFwcZw56KjuAA4iMuSfMxCbu1UyFw"
const validAccessTokenJTI = "75691799839643f89b64fe4461e49be0"

const defaultResponseJSON = `{"type":"%s","id":"` + uuid + `","name":"%s","version_created_at":"` + timestamp + `","value":%s,"metadata":%s}`
const redactedResponseJSON = `{"type":"%s","id":"` + uuid + `","name":"%s","version_created_at":"` + timestamp + `","value":"<redacted>", "metadata":%s}`
const arrayResponseJSON = `{"data":[` + defaultResponseJSON + `]}`
const multipleCredentialArrayResponseJSON = `{"data":[` + defaultResponseJSON + `,` + defaultResponseJSON + `]}`

const generateRequestJSON = `{"type":"%s","name":"%s","parameters":%s,"overwrite":%t}`
const generateWithValueRequestJSON = `{"type":"%s","name":"%s","parameters":%s,"overwrite":%t,"value":%s}`
const generateResponseJSON = `{"type":"%s","id":"` + uuid + `","name":"%s","version_created_at":"` + timestamp + `","value":%s,"metadata":{}}`

const addPermissionsRequestJSON = `{"path":"%s","actor":"%s","operations":%s}`
const permissionsResponseJSON = `{"uuid":"` + uuid + `","path":"%s","actor":"%s","operations":%s}`

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var (
	commandPath string
	homeDir     string
	server      *Server
	authServer  *Server
	credhubEnv  map[string]string
)

var _ = BeforeEach(func() {
	var err error
	homeDir, err = ioutil.TempDir("", "cm-test")
	Expect(err).NotTo(HaveOccurred())

	if runtime.GOOS == "windows" {
		os.Setenv("USERPROFILE", homeDir)
	} else {
		os.Setenv("HOME", homeDir)
	}

	server = NewTlsServer("../test/server-tls-cert.pem", "../test/server-tls-key.pem")
	authServer = NewTlsServer("../test/auth-tls-cert.pem", "../test/auth-tls-key.pem")

	SetupServers(server, authServer)

	session := runCommand("api", server.URL(), "--ca-cert", "../test/server-tls-ca.pem", "--ca-cert", "../test/auth-tls-ca.pem")

	server.Reset()
	authServer.Reset()

	Eventually(session).Should(Exit(0))
})

var _ = AfterEach(func() {
	server.Close()
	authServer.Close()
	os.RemoveAll(homeDir)
})

var _ = SynchronizedBeforeSuite(func() []byte {
	executablePath, err := Build("code.cloudfoundry.org/credhub-cli", "-ldflags", "-X code.cloudfoundry.org/credhub-cli/version.Version=test-version")
	Expect(err).NotTo(HaveOccurred())
	return []byte(executablePath)
}, func(data []byte) {
	commandPath = string(data)
	credhubEnv = test_util.UnsetAndCacheCredHubEnvVars()
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
	test_util.RestoreEnv(credhubEnv)
})

func login() {
	authServer.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", "/oauth/token"),
			RespondWith(http.StatusOK, `{
			"access_token":"test-access-token",
			"refresh_token":"test-refresh-token",
			"token_type":"password",
			"expires_in":123456789
			}`),
		),
	)

	server.RouteToHandler("GET", "/info",
		RespondWith(http.StatusOK, `{
				"app":{"name":"CredHub"}
				}`),
	)

	server.RouteToHandler("GET", "/version",
		RespondWith(http.StatusOK, `{"version":"9.9.9"}`),
	)

	runCommand("login", "-u", "test-username", "-p", "test-password")

	authServer.Reset()
}

func resetCachedServerVersion() {
	setCachedServerVersion("")
}

func setCachedServerVersion(version string) {
	cfg := config.ReadConfig()
	cfg.ServerVersion = version
	config.WriteConfig(cfg)
}

func runCommand(args ...string) *Session {
	cmd := exec.Command(commandPath, args...)
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

func runCommandWithEnv(env []string, args ...string) *Session {
	cmd := exec.Command(commandPath, args...)
	existing := os.Environ()
	for _, env_var := range env {
		existing = append(existing, env_var)
	}
	cmd.Env = existing
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

func runCommandWithStdin(stdin io.Reader, args ...string) *Session {
	cmd := exec.Command(commandPath, args...)
	cmd.Stdin = stdin
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

func setupUAAConfig(uaaResponseStatus int) {
	cfg := config.Config{
		ConfigWithoutSecrets: config.ConfigWithoutSecrets{
			RefreshToken: "5b9c9fd51ba14838ac2e6b222d487106-r",
			AccessToken:  "e30K.eyJqdGkiOiIxIn0K.e30K",
			AuthURL:      authServer.URL(),
			ApiURL:       server.URL(),
		},
	}

	Expect(cfg.UpdateTrustedCAs([]string{"../test/auth-tls-ca.pem", "../test/server-tls-ca.pem"})).To(Succeed())
	Expect(config.WriteConfig(cfg)).To(Succeed())

	authServer.RouteToHandler(
		"DELETE", "/oauth/token/revoke/1",
		RespondWith(uaaResponseStatus, ""),
	)
}

func NewTlsServer(certPath, keyPath string) *Server {
	tlsServer := NewUnstartedServer()

	cert, err := ioutil.ReadFile(certPath)
	Expect(err).To(BeNil())
	key, err := ioutil.ReadFile(keyPath)
	Expect(err).To(BeNil())

	tlsCert, err := tls.X509KeyPair(cert, key)
	Expect(err).To(BeNil())

	tlsServer.HTTPTestServer.TLS = &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	tlsServer.HTTPTestServer.StartTLS()

	return tlsServer
}

func SetupServers(chServer, uaaServer *Server) {
	chServer.RouteToHandler("GET", "/info",
		RespondWith(http.StatusOK, `{
				"app":{"name":"CredHub"},
				"auth-server":{"url":"`+uaaServer.URL()+`"}
				}`),
	)

	chServer.RouteToHandler("GET", "/version",
		RespondWith(http.StatusOK, `{"version":"9.9.9"}`),
	)

	uaaServer.RouteToHandler("GET", "/info", RespondWith(http.StatusOK, ""))
}

func ItBehavesLikeHelp(command string, alias string, validate func(*Session)) {
	It("displays help", func() {
		session := runCommand(command, "-h")
		Eventually(session).Should(Exit(1))
		validate(session)
	})

	It("displays help using the alias", func() {
		session := runCommand(alias, "-h")
		Eventually(session).Should(Exit(1))
		validate(session)
	})
}

func ItRequiresAuthentication(args ...string) {
	It("requires authentication", func() {
		setupUAAConfig(http.StatusOK)

		runCommand("logout")

		session := runCommand(args...)

		Eventually(session).Should(Exit(1))
		Expect(session.Err).To(Say("You are not currently authenticated. Please log in to continue."))
	})
}

func ItRequiresAnAPIToBeSet(args ...string) {
	Describe("requires an API endpoint", func() {
		BeforeEach(func() {
			cfg := config.ReadConfig()
			cfg.ApiURL = ""
			config.WriteConfig(cfg)
		})

		Context("when using password grant", func() {
			It("requires an API endpoint", func() {
				session := runCommand(args...)

				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("An API target is not set. Please target the location of your server with `credhub api --server api.example.com` to continue."))
			})
		})

		Context("when using client_credentials", func() {
			It("requires an API endpoint", func() {
				session := runCommandWithEnv([]string{"CREDHUB_CLIENT=test_client", "CREDHUB_SECRET=test_secret"}, args...)

				Eventually(session).Should(Exit(1))
				Expect(session.Err).To(Say("An API target is not set. Please target the location of your server with `credhub api --server api.example.com` to continue."))
			})
		})
	})
}

type TestAutoLogin struct {
	method              string
	responseFixtureFile string
	responseStatus      int
	endpoint            string
}

func ItAutomaticallyLogsIn(autoLogins []TestAutoLogin, args ...string) {
	var serverResponse = make([]string, len(autoLogins))
	Describe("automatic authentication", func() {
		BeforeEach(func() {
			for i, autoLogin := range autoLogins {
				buf, _ := ioutil.ReadFile(filepath.Join("testdata", autoLogin.responseFixtureFile))
				serverResponse[i] = string(buf)
			}
		})
		AfterEach(func() {
			server.Reset()
		})

		Context("with correct environment and unauthenticated", func() {
			BeforeEach(func() {
				for i, autoLogin := range autoLogins {
					server.AppendHandlers(
						CombineHandlers(
							VerifyRequest(autoLogin.method, autoLogin.endpoint),
							VerifyHeader(http.Header{
								"Authorization": []string{"Bearer 2YotnFZFEjr1zCsicMWpAA"},
							}),
							RespondWith(autoLogin.responseStatus, serverResponse[i]),
						),
					)
				}
			})

			It("automatically authenticates", func() {
				setupUAAConfig(http.StatusOK)

				authServer.AppendHandlers(
					CombineHandlers(
						VerifyRequest("POST", "/oauth/token"),
						VerifyBody([]byte(`client_id=test_client&client_secret=test_secret&grant_type=client_credentials&response_type=token`)),
						RespondWith(http.StatusOK, `{
								"access_token":"2YotnFZFEjr1zCsicMWpAA",
								"token_type":"bearer",
								"expires_in":3600}`),
					),
				)

				runCommand("logout")

				session := runCommandWithEnv([]string{"CREDHUB_CLIENT=test_client", "CREDHUB_SECRET=test_secret"}, args...)

				Eventually(session).Should(Exit(0))
			})
		})

		Context("with correct environment and expired token", func() {
			BeforeEach(func() {

				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(autoLogins[0].method, autoLogins[0].endpoint),
						VerifyHeader(http.Header{
							"Authorization": []string{"Bearer test-access-token"},
						}),
						RespondWith(http.StatusUnauthorized, `{
						"error":"access_token_expired",
						"error_description":"error description"}`),
					),
				)

				authServer.AppendHandlers(
					CombineHandlers(
						VerifyRequest("POST", "/oauth/token"),
						VerifyBody([]byte(`client_id=test_client&client_secret=test_secret&grant_type=client_credentials&response_type=token`)),
						RespondWith(http.StatusOK, `{
								"access_token":"new-token",
								"token_type":"bearer",
								"expires_in":3600}`),
					),
				)

				for i, autoLogin := range autoLogins {

					server.AppendHandlers(
						CombineHandlers(
							VerifyRequest(autoLogin.method, autoLogin.endpoint),
							VerifyHeader(http.Header{
								"Authorization": []string{"Bearer new-token"},
							}),
							RespondWith(autoLogin.responseStatus, serverResponse[i]),
						),
					)
				}
			})

			It("automatically authenticates", func() {
				session := runCommandWithEnv([]string{"CREDHUB_CLIENT=test_client", "CREDHUB_SECRET=test_secret"}, args...)
				Eventually(session).Should(Exit(0))
			})
		})
	})
}
