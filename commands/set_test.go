package commands_test

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"code.cloudfoundry.org/credhub-cli/commands"
	"code.cloudfoundry.org/credhub-cli/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = FDescribe("Set", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("set", "-n", "test-credential", "-t", "password", "-w", "value")
	ItRequiresAnAPIToBeSet("set", "-n", "test-credential", "-t", "password", "-w", "value")
	testAutoLogin := []TestAutoLogin{
		{
			method:              "PUT",
			responseFixtureFile: "set_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v1/data",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogin, "set", "-n", "test-credential", "-t", "password", "-w", "test-value")

	Describe("not specifying type", func() {
		It("returns an error", func() {
			session := runCommand("set", "-n", "my-password", "-w", "potatoes")

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("A type must be specified when setting a credential. Valid types include 'value', 'json', 'password', 'user', 'certificate', 'ssh' and 'rsa'."))
		})
	})

	Describe("setting value secrets", func() {
		It("puts a secret", func() {
			setupSetServer("my-value", "value", `"potatoes"`)

			session := runCommand("set", "-n", "my-value", "-v", "potatoes", "-t", "value")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-value"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: value"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("escapes special characters in the value", func() {
			setupSetServer("my-character-test", "value", `"{\"password\":\"some-still-bad-password\"}"`)

			session := runCommand("set", "-t", "value", "-n", "my-character-test", "-v", `{"password":"some-still-bad-password"}`)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-character-test"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: value"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`value: <redacted>`))
		})

		It("puts a secret and returns in json format", func() {
			setupSetServer("my-value", "value", `"potatoes"`)

			session := runCommand("set", "-n", "my-value", "-v", "potatoes", "-t", "value", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(responseSetMyValuePotatoesJson))
		})

		It("accepts case-insensitive type", func() {
			setupSetServer("my-value", "value", `"potatoes"`)

			session := runCommand("set", "-n", "my-value", "-v", "potatoes", "-t", "VALUE", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(responseSetMyValuePotatoesJson))
		})
	})

	Describe("setting json secrets", func() {
		It("puts a secret", func() {
			jsonValue := `{"foo":"bar","nested":{"a":1},"an":["array"]}`
			setupSetServer("json-secret", "json", jsonValue)

			session := runCommand("set", "-n", "json-secret", "-v", jsonValue, "-t", "json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: json-secret"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: json"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("escapes special characters in the json", func() {
			setupSetServer("my-character-test", "json", `{"foo":"b\"ar"}`)

			session := runCommand("set", "-t", "json", "-n", "my-character-test", "-v", `{"foo":"b\"ar"}`)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-character-test"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: json"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret using explicit json type and returns in json format", func() {
			jsonValue := `{"foo":"bar","nested":{"a":1},"an":["array"]}`
			setupSetServer("json-secret", "json", jsonValue)

			session := runCommand("set", "-n", "json-secret", "-v", jsonValue, "-t", "json", "--output-json")

			Eventually(session).Should(Exit(0))
			responseJson := `{
			"id": "5a2edd4f-1686-4c8d-80eb-5daa866f9f86",
			"name": "json-secret",
			"type": "json",
			"value": "<redacted>",
			"version_created_at": "2016-01-01T12:00:00Z"
			}`
			Eventually(string(session.Out.Contents())).Should(MatchJSON(responseJson))
		})

		It("accepts case-insensitive type", func() {
			jsonValue := `{"foo":"bar","nested":{"a":1},"an":["array"]}`
			setupSetServer("json-secret", "json", jsonValue)

			session := runCommand("set", "-n", "json-secret", "-v", jsonValue, "-t", "JSON")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: json-secret"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: json"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})
	})

	Describe("setting SSH secrets", func() {
		It("puts a secret using explicit ssh type", func() {
			setupSetServer("foo-ssh-key", "ssh", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			session := runCommand("set", "-n", "foo-ssh-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "ssh")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: foo-ssh-key"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: ssh"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret using values read from files", func() {
			setupSetServer("foo-ssh-key", "ssh", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			tempDir := test.CreateTempDir("sshFilesForTesting")
			publicFileName := test.CreateCredentialFile(tempDir, "rsa.pub", "some-public-key")
			privateFilename := test.CreateCredentialFile(tempDir, "rsa.key", "some-private-key")

			session := runCommand("set", "-n", "foo-ssh-key",
				"-t", "ssh",
				"-u", publicFileName,
				"-p", privateFilename)

			os.RemoveAll(tempDir)
			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: foo-ssh-key"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: ssh"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret using explicit ssh type and returns in json format", func() {
			setupSetServer("foo-ssh-key", "ssh", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			session := runCommand("set", "-n", "foo-ssh-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "ssh", "--output-json")

			Eventually(session).Should(Exit(0))
			responseJson := `{
        	"id": "5a2edd4f-1686-4c8d-80eb-5daa866f9f86",
        	"name": "foo-ssh-key",
        	"type": "ssh",
        	"value": "<redacted>",
        	"version_created_at": "2016-01-01T12:00:00Z"
        	}`
			Eventually(string(session.Out.Contents())).Should(MatchJSON(responseJson))
		})

		It("accepts case-insensitive type", func() {
			setupSetServer("foo-ssh-key", "ssh", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			session := runCommand("set", "-n", "foo-ssh-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "SSH")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: foo-ssh-key"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: ssh"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("handles newline characters", func() {
			setupSetServer("foo-ssh-key", "ssh", `{"public_key":"some\npublic\nkey","private_key":"some\nprivate\nkey"}`)
			session := runCommand("set", "-n", "foo-ssh-key", "-u", `some\npublic\nkey`, "-p", `some\nprivate\nkey`, "-t", "ssh", "--output-json")

			responseJson := `{
        	"id": "5a2edd4f-1686-4c8d-80eb-5daa866f9f86",
        	"name": "foo-ssh-key",
        	"type": "ssh",
        	"value": "<redacted>",
        	"version_created_at": "2016-01-01T12:00:00Z"
        	}`

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(responseJson))
		})
	})

	Describe("setting RSA secrets", func() {
		It("puts a secret using explicit rsa type", func() {
			setupSetServer("foo-rsa-key", "rsa", `{"public_key":"some-public-key","private_key":"some-private-key"}`)
			session := runCommand("set", "-n", "foo-rsa-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "rsa")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: foo-rsa-key"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: rsa"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret using values read from files", func() {
			setupSetServer("foo-rsa-key", "rsa", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			tempDir := test.CreateTempDir("rsaFilesForTesting")
			publicFileName := test.CreateCredentialFile(tempDir, "rsa.pub", "some-public-key")
			privateFilename := test.CreateCredentialFile(tempDir, "rsa.key", "some-private-key")

			session := runCommand("set", "-n", "foo-rsa-key",
				"-t", "rsa",
				"-u", publicFileName,
				"-p", privateFilename)

			os.RemoveAll(tempDir)
			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: foo-rsa-key"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: rsa"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret using explicit rsa type and returns in json format", func() {
			setupSetServer("foo-rsa-key", "rsa", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			session := runCommand("set", "-n", "foo-rsa-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "rsa", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"name": "foo-rsa-key"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"type": "rsa"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"value": "<redacted>"`))

		})

		It("accepts case-insensitive type", func() {
			setupSetServer("foo-rsa-key", "rsa", `{"public_key":"some-public-key","private_key":"some-private-key"}`)

			session := runCommand("set", "-n", "foo-rsa-key", "-u", "some-public-key", "-p", "some-private-key", "-t", "RSA")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: foo-rsa-key"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: rsa"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("handles newline characters", func() {
			setupSetServer("foo-rsa-key", "rsa", `{"public_key":"some\npublic\nkey","private_key":"some\nprivate\nkey"}`)

			session := runCommand("set", "-n", "foo-rsa-key", "-u", `some\npublic\nkey`, "-p", `some\nprivate\nkey`, "-t", "rsa", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).Should(MatchJSON(responseSetMyRSAWithNewlinesJson))
		})
	})

	Describe("setting password secrets", func() {

		It("puts a secret using explicit password type  and returns in yaml format", func() {
			setupSetServer("my-password", "password", `"potatoes"`)

			session := runCommand("set", "-n", "my-password", "-w", "potatoes", "-t", "password")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("prompts for value if value is not provided", func() {
			setupSetServer("my-password", "password", `"potatoes"`)

			session := runCommandWithStdin(strings.NewReader("potatoes\n"), "set", "-n", "my-password", "-t", "password")

			Eventually(string(session.Out.Contents())).Should(ContainSubstring("password: ********"))
			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("can set password that contains spaces interactively", func() {
			setupSetServer("my-password", "password", `"potatoes potatoes"`)

			session := runCommandWithStdin(strings.NewReader("potatoes potatoes\n"), "set", "-n", "my-password", "-t", "password")

			Eventually(string(session.Out.Contents())).Should(ContainSubstring("password:"))
			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("escapes special characters in the password", func() {
			setupSetServer("my-password", "password", `"{\"password\":\"some-still-bad-password\"}"`)

			session := runCommand("set", "-t", "password", "-n", "my-password", "-w", `{"password":"some-still-bad-password"}`)

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`value: <redacted>`))
		})

		It("puts a secret using explicit password type and returns in json format", func() {
			setupSetServer("my-password", "password", `"potatoes"`)

			session := runCommand("set", "-n", "my-password", "-w", "potatoes", "-t", "password", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(MatchJSON(responseSetMyPasswordPotatoesJson))
		})

		It("accepts case-insensitive type", func() {
			setupSetServer("my-password", "password", `"potatoes"`)

			session := runCommand("set", "-n", "my-password", "-w", "potatoes", "-t", "PASSWORD")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: password"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})
	})

	Describe("setting certificate secrets", func() {
		It("puts a secret using explicit certificate type and string values", func() {
			setupSetServer("my-secret", "certificate", `{"ca":"my-ca","certificate":"my-cert","private_key":"my-priv"}`)
			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root", "my-ca",
				"--certificate", "my-cert", "--private", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-secret"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: certificate"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret using explicit certificate type, string values, and certificate authority name", func() {
			setupSetServer("my-secret", "certificate", `{"ca": "", "ca_name":"my-ca","certificate":"my-cert","private_key":"my-priv"}`)

			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--ca-name", "my-ca",
				"--certificate", "my-cert", "--private", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-secret"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: certificate"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret using explicit certificate type and values read from files", func() {
			setupSetServer("my-secret", "certificate", `{"ca":"my-ca","certificate":"my-cert","private_key":"my-priv"}`)
			tempDir := test.CreateTempDir("certFilesForTesting")
			caFilename := test.CreateCredentialFile(tempDir, "ca.txt", "my-ca")
			certificateFilename := test.CreateCredentialFile(tempDir, "certificate.txt", "my-cert")
			privateFilename := test.CreateCredentialFile(tempDir, "private.txt", "my-priv")

			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root", caFilename,
				"--certificate", certificateFilename, "--private", privateFilename)

			os.RemoveAll(tempDir)
			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-secret"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: certificate"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		if runtime.GOOS != "windows" {
			It("fails to put a secret when reading from unreadable file", func() {
				testSetFileFailure("unreadable.txt", "", "")
				testSetFileFailure("", "unreadable.txt", "")
				testSetFileFailure("", "", "unreadable.txt")
			})
		}

		It("puts a secret using explicit certificate type and string values in json format", func() {
			setupSetServer("my-secret", "certificate", `{"ca":"my-ca","certificate":"my-cert","private_key":"my-priv"}`)

			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root", "my-ca",
				"--certificate", "my-cert", "--private", "my-priv", "--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"name": "my-secret"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"type": "certificate"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"value": "<redacted>"`))
		})

		It("accepts case insensitive type", func() {
			setupSetServer("my-secret", "certificate", `{"ca":"my-ca","certificate":"my-cert","private_key":"my-priv"}`)

			session := runCommand("set", "-n", "my-secret",
				"-t", "CERTIFICATE", "--root", "my-ca",
				"--certificate", "my-cert", "--private", "my-priv")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-secret"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: certificate"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("handles newline characters", func() {
			setupSetServer("my-secret", "certificate", `{"ca":"my\nca","certificate":"my\ncert","private_key":"my\npriv"}`)
			session := runCommand("set", "-n", "my-secret",
				"-t", "certificate", "--root", `my\nca`,
				"--certificate", `my\ncert`, "--private", `my\npriv`, "--output-json")
			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).Should(MatchJSON(responseSetMyCertificateWithNewlinesJson))
		})
	})

	Describe("setting User secrets", func() {
		It("puts a secret using explicit user type", func() {
			//SetupPutUserServer("my-username-credential", `{"username": "my-username", "password": "test-password"}`, "my-username", "test-password", "passw0rd-H4$h")
			setupSetServer("my-username-credential", "user", `{"username": "my-username", "password": "test-password"}`)
			session := runCommand("set", "-n", "my-username-credential", "-z", "my-username", "-w", "test-password", "-t", "user")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-username-credential"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: user"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`version_created_at: "2016-01-01T12:00:00Z"`))
		})

		It("should set password interactively for user", func() {
			setupSetServer("my-username-credential", "user", `{"username": "my-username", "password": "test-password"}`)
			session := runCommandWithStdin(strings.NewReader("test-password\n"), "set", "-n", "my-username-credential", "-t", "user", "--username", "my-username")

			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-username-credential"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: user"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
			Eventually(session).Should(Exit(0))
		})

		It("should set null username when it isn't provided", func() {
			setupSetServer("my-username-credential", "user", `{"username": "", "password": "test-password"}`)

			session := runCommandWithStdin(strings.NewReader("test-password\n"), "set", "-n", "my-username-credential", "-t", "user")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-username-credential"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: user"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
		})

		It("puts a secret using explicit user type in json format", func() {
			setupSetServer("my-username-credential", "user", `{"username": "my-username", "password": "test-password"}`)

			session := runCommand("set", "-n", "my-username-credential", "-z", "my-username", "-w", "test-password", "-t", "user",
				"--output-json")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"name": "my-username-credential"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"type": "user"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"value": "<redacted>"`))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`"version_created_at": "2016-01-01T12:00:00Z"`))
		})

		It("accepts case-insensitive type", func() {
			setupSetServer("my-username-credential", "user", `{"username": "my-username", "password": "test-password"}`)
			session := runCommand("set", "-n", "my-username-credential", "-z", "my-username", "-w", "test-password", "-t", "USER")

			Eventually(session).Should(Exit(0))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("name: my-username-credential"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("type: user"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring("value: <redacted>"))
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(`version_created_at: "2016-01-01T12:00:00Z"`))
		})
	})

	Describe("Help", func() {
		It("short flags", func() {
			Expect(commands.SetCommand{}).To(SatisfyAll(
				commands.HaveFlag("name", "n"),
				commands.HaveFlag("type", "t"),
				commands.HaveFlag("value", "v"),
				commands.HaveFlag("root", "r"),
				commands.HaveFlag("certificate", "c"),
				commands.HaveFlag("private", "p"),
			))
		})

		ItBehavesLikeHelp("set", "s", func(session *Session) {
			Expect(session.Err).To(Say("set"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("credential"))
		})

		It("displays missing 'n' option as required parameter", func() {
			session := runCommand("set", "-v", "potatoes")

			Eventually(session).Should(Exit(1))
			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flag `/n, /name' was not specified"))
			} else {
				Expect(session.Err).To(Say("the required flag `-n, --name' was not specified"))
			}
		})

		It("displays the server provided error when an error is received", func() {
			server.AppendHandlers(
				RespondWith(http.StatusBadRequest, `{"error": "test error"}`),
			)

			session := runCommand("set", "-n", "my-value", "-t", "value", "-v", "tomatoes")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("test error"))
		})
	})
})

const SET_CREDENTIAL_REQUEST_JSON = `{"type":"%s","name":"%s","value":%s}`
const SET_CREDENTIAL_RESPONSE_JSON = `{"type":"%s","id":"` + UUID + `","name":"%s","version_created_at":"` + TIMESTAMP + `","value":%s,"version_created_at":"` + TIMESTAMP + `"}`

func setupSetServer(name, keyType, value string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", "/api/v1/data"),
			VerifyJSON(fmt.Sprintf(SET_CREDENTIAL_REQUEST_JSON, keyType, name, value)),
			RespondWith(http.StatusOK, fmt.Sprintf(SET_CREDENTIAL_RESPONSE_JSON, keyType, name, value)),
		),
	)
}

func setupPutBadRequestServer(body string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("PUT", "/api/v1/data"),
			VerifyJSON(body),
			RespondWith(http.StatusBadRequest, `{"error":"test error"}`),
		),
	)
}

func testSetFileFailure(caFilename, certificateFilename, privateFilename string) {
	tempDir := test.CreateTempDir("certFilesForTesting")
	if caFilename == "unreadable.txt" {
		caFilename = test.CreateCredentialFile(tempDir, caFilename, "my-ca")
		err := os.Chmod(caFilename, 0222)
		Expect(err).To(BeNil())
	}
	if certificateFilename == "unreadable.txt" {
		certificateFilename = test.CreateCredentialFile(tempDir, certificateFilename, "my-cert")
		err := os.Chmod(certificateFilename, 0222)
		Expect(err).To(BeNil())
	}
	if privateFilename == "unreadable.txt" {
		privateFilename = test.CreateCredentialFile(tempDir, privateFilename, "my-priv")
		err := os.Chmod(privateFilename, 0222)
		Expect(err).To(BeNil())
	}

	session := runCommand("set", "-n", "my-secret",
		"-t", "certificate", "--root", caFilename,
		"--certificate", certificateFilename, "--private", privateFilename)

	os.RemoveAll(tempDir)
	Eventually(session).Should(Exit(1))
	Eventually(session.Err).Should(Say("A referenced file could not be opened. Please validate the provided filenames and permissions, then retry your request."))
}
