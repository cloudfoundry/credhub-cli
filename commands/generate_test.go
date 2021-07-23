package commands_test

import (
	"net/http"
	"runtime"

	"fmt"

	"code.cloudfoundry.org/credhub-cli/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

const generateRequestJSONWithMetadata = `{"type":"%s","name":"%s","parameters":%s,"overwrite":%t, "metadata":%s}`
const generateResponseJSONWithMetadata = `{"type":"%s","id":"` + uuid + `","name":"%s","version_created_at":"` + timestamp + `","value":%s,"metadata":%s}`
const generateResponseJSONWithDurationOverriddenFlag = `{"type":"%s","id":"` + uuid + `","name":"%s","version_created_at":"` + timestamp + `","value":%s,"duration_overridden":true}`

var _ = Describe("Generate", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("generate", "-n", "test-credential", "-t", "password")
	ItRequiresAnAPIToBeSet("generate", "-n", "test-credential", "-t", "password")
	testAutoLogin := []TestAutoLogin{
		{
			method:              "POST",
			responseFixtureFile: "generate_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v1/data",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogin, "generate", "-n", "test-credential", "-t", "password")

	It("requires a type", func() {
		session := runCommand("generate", "-n", "my-credential")
		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(Say("A type must be specified when generating a credential. Valid types include 'password', 'user', 'certificate', 'ssh' and 'rsa'."))
	})

	Describe("Without password parameters", func() {
		BeforeEach(func() {
			setupGenerateServer("password", "my-password", `"potatoes"`, `{}`, true)
		})
		It("uses default parameters", func() {
			session := runCommand("generate", "-n", "my-password", "-t", "password")
			Eventually(session).Should(Exit(0))
		})

		It("prints the generated password secret", func() {
			session := runCommand("generate", "-n", "my-password", "-t", "password")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: my-password"))
			Eventually(session.Out).Should(Say("type: password"))
			Eventually(session.Out).Should(Say("value: <redacted>"))
		})

		It("can print the generated password secret as JSON", func() {
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(fmt.Sprintf(redactedResponseJSON, "password", "my-password", "{}")))
		})

		It("allows the type to be any case", func() {
			session := runCommand("generate", "-n", "my-password", "-t", "PASSWORD")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of password parameters", func() {
		It("can print the secret as JSON", func() {
			setupGenerateServer("password", "my-password", `"potatoes"`, `{}`, true)

			session := runCommand(
				"generate",
				"-n", "my-password",
				"-t", "password",
				"--output-json",
			)

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(fmt.Sprintf(redactedResponseJSON, "password", "my-password", "{}")))
		})

		It("with with no-overwrite", func() {
			setupGenerateServer("password", "my-password", `"potatoes"`, `{}`, false)
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupGenerateServer("password", "my-password", `"potatoes"`, `{"length":42}`, true)
			session := runCommand("generate", "-n", "my-password", "-t", "password", "-l", "42")
			Eventually(session).Should(Exit(0))
		})

		It("excluding upper case", func() {
			setupGenerateServer("password", "my-password", `"potatoes"`, `{"exclude_upper":true}`, true)
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-upper")
			Eventually(session).Should(Exit(0))
		})

		It("excluding lower case", func() {
			setupGenerateServer("password", "my-password", `"potatoes"`, `{"exclude_lower":true}`, true)
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-lower")
			Eventually(session).Should(Exit(0))
		})

		It("including special characters", func() {
			setupGenerateServer("password", "my-password", `"potatoes"`, `{"include_special":true}`, true)
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--include-special")
			Eventually(session).Should(Exit(0))
		})

		It("excluding numbers", func() {
			setupGenerateServer("password", "my-password", `"potatoes"`, `{"exclude_number":true}`, true)
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-number")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of SSH parameters", func() {
		It("prints the SSH key", func() {
			setupGenerateServer("ssh", "foo-ssh-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{}`, true)

			session := runCommand("generate", "-n", "foo-ssh-key", "-t", "ssh")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: foo-ssh-key"))
			Eventually(session.Out).Should(Say("value: <redacted>"))
		})

		It("can print the SSH key as JSON", func() {
			setupGenerateServer("ssh", "foo-ssh-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{}`, true)

			session := runCommand("generate", "-n", "foo-ssh-key", "-t", "ssh", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(fmt.Sprintf(redactedResponseJSON, "ssh", "foo-ssh-key", "{}")))
		})

		It("with with no-overwrite", func() {
			setupGenerateServer("ssh", "foo-ssh-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{}`, false)

			session := runCommand("generate", "-n", "foo-ssh-key", "-t", "ssh", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupGenerateServer("ssh", "foo-ssh-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{"key_length":3072}`, true)
			session := runCommand("generate", "-n", "foo-ssh-key", "-t", "ssh", "-k", "3072")
			Eventually(session).Should(Exit(0))
		})

		It("including comment", func() {
			setupGenerateServer("ssh", "foo-ssh-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{"ssh_comment":"i am an ssh comment"}`, true)
			session := runCommand("generate", "-n", "foo-ssh-key", "-t", "ssh", "-m", "i am an ssh comment")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of RSA parameters", func() {
		It("prints the RSA key", func() {
			setupGenerateServer("rsa", "foo-rsa-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{}`, true)

			session := runCommand("generate", "-n", "foo-rsa-key", "-t", "rsa")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: foo-rsa-key"))
			Eventually(session.Out).Should(Say("value: <redacted>"))
		})

		It("allows the type to be any case", func() {
			setupGenerateServer("rsa", "foo-rsa-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{}`, true)

			session := runCommand("generate", "-n", "foo-rsa-key", "-t", "RSA")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: foo-rsa-key"))
		})

		It("can print the RSA key as JSON", func() {
			setupGenerateServer("rsa", "foo-rsa-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{}`, true)

			session := runCommand("generate", "-n", "foo-rsa-key", "-t", "rsa", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(fmt.Sprintf(redactedResponseJSON, "rsa", "foo-rsa-key", "{}")))
		})

		It("with with no-overwrite", func() {
			setupGenerateServer("rsa", "foo-rsa-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{}`, false)
			session := runCommand("generate", "-n", "foo-rsa-key", "-t", "rsa", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupGenerateServer("rsa", "foo-rsa-key", `{"public_key":"some-public-key","private_key":"some-private-key"}`, `{"key_length":3072}`, true)
			session := runCommand("generate", "-n", "foo-rsa-key", "-t", "rsa", "-k", "3072")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of certificate parameters", func() {
		It("prints the certificate", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"common_name":"common.name.io"}`, true)

			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: my-secret"))
			Eventually(session.Out).Should(Say("type: certificate"))
			Eventually(session.Out).Should(Say("value: <redacted>"))
			Eventually(session.Out).ShouldNot(Say("duration_overridden:"))
		})

		It("allows the type to be any case", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"common_name":"common.name.io"}`, true)

			session := runCommand("generate", "-n", "my-secret", "-t", "CERTIFICATE", "--common-name", "common.name.io")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: my-secret"))
		})

		It("can print the certificate as JSON", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"common_name":"common.name.io"}`, true)

			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(fmt.Sprintf(redactedResponseJSON, "certificate", "my-secret", "{}")))
		})

		It("including common name with no-overwrite", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"common_name":"common.name.io"}`, false)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including organization", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"organization":"organization.io"}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--organization", "organization.io")
			Eventually(session).Should(Exit(0))
		})

		It("including organization unit", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"organization_unit":"My Unit"}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--organization-unit", "My Unit")
			Eventually(session).Should(Exit(0))
		})

		It("including locality", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"locality":"My Locality"}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--locality", "My Locality")
			Eventually(session).Should(Exit(0))
		})

		It("including state", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"state":"My State"}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--state", "My State")
			Eventually(session).Should(Exit(0))
		})

		It("including country", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"country":"My Country"}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--country", "My Country")
			Eventually(session).Should(Exit(0))
		})

		It("including multiple alternative names", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"alternative_names": [ "Alt1", "Alt2" ]}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--alternative-name", "Alt1", "--alternative-name", "Alt2")
			Eventually(session).Should(Exit(0))
		})

		It("including multiple extended key usage settings", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"extended_key_usage": [ "server_auth", "client_auth" ]}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "-e", "server_auth", "--ext-key-usage=client_auth")
			Eventually(session).Should(Exit(0))
		})

		It("including multiple key usage settings", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"key_usage": ["digital_signature", "non_repudiation"]}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "-g", "digital_signature", "--key-usage=non_repudiation")
			Eventually(session).Should(Exit(0))
		})

		It("including key length", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"key_length":2048}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--key-length", "2048")
			Eventually(session).Should(Exit(0))
		})

		It("including duration", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"duration":1000}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--duration", "1000")
			Eventually(session).Should(Exit(0))
		})

		It("including duration less than minimum duration", func() {
			setupGenerateServerWithDurationOverriddenFlag("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"duration":365}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--duration", "365")
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("duration_overridden: true"))
		})

		It("including certificate authority", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"my_ca","certificate":"my-cert","private_key":"my-priv"}`, `{"ca":"my_ca"}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--ca", "my_ca")
			Eventually(session).Should(Exit(0))
		})

		It("including self-signed flag", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"","certificate":"my-cert","private_key":"my-priv"}`, `{"self_sign": true, "common_name": "my.name.io"}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "-c", "my.name.io", "--self-sign")
			Eventually(session).Should(Exit(0))
		})

		It("including is-ca flag", func() {
			setupGenerateServer("certificate", "my-secret", `{"ca":"my-cert","certificate":"my-cert","private_key":"my-priv"}`, `{"is_ca": true, "common_name": "my.name.io"}`, true)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "-c", "my.name.io", "--is-ca")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of user parameters", func() {
		It("prints the secret", func() {
			setupGenerateServer(
				"user",
				"my-username-credential",
				`{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4Sh"}`,
				`{}`,
				true)

			session := runCommand("generate", "-n", "my-username-credential", "-t", "user")

			Eventually(session).Should(Exit(0))
			Expect(session.Out.Contents()).To(ContainSubstring("name: my-username-credential"))
			Expect(session.Out.Contents()).To(ContainSubstring("type: user"))
			Expect(session.Out.Contents()).To(ContainSubstring("value: <redacted>"))
		})

		It("allows the type to be any case", func() {
			setupGenerateServer(
				"user",
				"my-username-credential",
				`{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4Sh"}`,
				`{}`,
				true)

			session := runCommand("generate", "-n", "my-username-credential", "-t", "USER")

			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("name: my-username-credential"))
		})

		It("should accept a statically provided username", func() {
			setupGenerateServerWithValue(
				"user",
				"my-username-credential",
				`{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4Sh"}`,
				`{}`,
				`{"username": "my-username"}`,
				true)

			session := runCommand("generate", "-n", "my-username-credential", "-t", "user", "-z", "my-username")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(ContainSubstring("name: my-username-credential"))
			Expect(string(session.Out.Contents())).To(ContainSubstring("value: <redacted>"))
		})

		It("with with no-overwrite", func() {
			setupGenerateServer(
				"user",
				"my-username-credential",
				`{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4Sh"}`,
				`{}`,
				false)
			session := runCommand("generate", "-n", "my-username-credential", "-t", "user", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupGenerateServer(
				"user",
				"my-username-credential",
				`{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4Sh"}`,
				`{"length": 42}`,
				true)
			session := runCommand("generate", "-n", "my-username-credential", "-t", "user", "-l", "42")
			Eventually(session).Should(Exit(0))
		})

		It("excluding upper case", func() {
			setupGenerateServer(
				"user",
				"my-username-credential",
				`{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4Sh"}`,
				`{"exclude_upper": true}`,
				true)
			session := runCommand("generate", "-n", "my-username-credential", "-t", "user", "--exclude-upper")
			Eventually(session).Should(Exit(0))
		})

		It("excluding lower case", func() {
			setupGenerateServer(
				"user",
				"my-username-credential",
				`{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4Sh"}`,
				`{"exclude_lower": true}`,
				true)
			session := runCommand("generate", "-n", "my-username-credential", "-t", "user", "--exclude-lower")
			Eventually(session).Should(Exit(0))
		})

		It("including special characters", func() {
			setupGenerateServer(
				"user",
				"my-username-credential",
				`{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4Sh"}`,
				`{"include_special": true}`,
				true)
			session := runCommand("generate", "-n", "my-username-credential", "-t", "user", "--include-special")
			Eventually(session).Should(Exit(0))
		})

		It("excluding numbers", func() {
			setupGenerateServer(
				"user",
				"my-username-credential",
				`{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4Sh"}`,
				`{"exclude_number": true}`,
				true)
			session := runCommand("generate", "-n", "my-username-credential", "-t", "user", "--exclude-number")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with metadata", func() {
		It("generates a secret with metadata", func() {
			setupGenerateServerWithMetadata(
				"user",
				"my-cred",
				`{"username":"my-username", "password":"test-password", "password_hash":"passw0rd-H4Sh"}`,
				`{}`,
				true,
				`{"some":{"example":"metadata"}, "array":["metadata"]}`,
			)

			session := runCommand("generate", "-t", "user", "-n", "my-cred", "--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			Eventually(session).Should(Exit(0))
			metadataOutput := `
metadata:
  array:
  - metadata
  some:
    example: metadata`
			Eventually(string(session.Out.Contents())).Should(ContainSubstring(metadataOutput))
		})

		It("errors when metadata is malformed", func() {
			setupGenerateServerWithMetadata(
				"user",
				"my-cred",
				`{error}`,
				`{}`,
				true,
				`"not-valid-json"`)

			session := runCommand("generate", "-t", "user", "-n", "my-cred", "--metadata", `"not-valid-json"`)

			Eventually(session).Should(Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("The argument for --metadata is not a valid json object. Please update and retry your request."))
		})

		It("errors when server does not support metadata", func() {
			setCachedServerVersion("2.5.0")

			session := runCommand("generate", "-t", "user", "-n", "my-cred", "--metadata", `{"some":{"example":"metadata"}, "array":["metadata"]}`)

			Eventually(session).Should(Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("The --metadata flag is not supported for this version of the credhub server (requires >= 2.6.x). Please remove the flag and retry your request."))
		})
	})

	Describe("When username parameter is included for non-user types", func() {
		It("returns a sensible error", func() {
			session := runCommand("generate", "-n", "test-ssh-value", "-t", "ssh", "-z", "my-username")
			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("Username parameter is not valid for this credential type."))
		})
	})

	Describe("Help", func() {
		ItBehavesLikeHelp("generate", "n", func(session *Session) {
			Expect(session.Err).To(Say("generate"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("length"))
		})

		It("short flags", func() {
			Expect(commands.GenerateCommand{}).To(SatisfyAll(
				commands.HaveFlag("name", "n"),
				commands.HaveFlag("type", "t"),
				commands.HaveFlag("no-overwrite", "O"),
				commands.HaveFlag("length", "l"),
				commands.HaveFlag("include-special", "S"),
				commands.HaveFlag("exclude-number", "N"),
				commands.HaveFlag("exclude-upper", "U"),
				commands.HaveFlag("exclude-lower", "L"),
				commands.HaveFlag("common-name", "c"),
				commands.HaveFlag("organization", "o"),
				commands.HaveFlag("organization-unit", "u"),
				commands.HaveFlag("locality", "i"),
				commands.HaveFlag("state", "s"),
				commands.HaveFlag("country", "y"),
				commands.HaveFlag("alternative-name", "a"),
				commands.HaveFlag("key-length", "k"),
				commands.HaveFlag("duration", "d"),
			))
		})

		It("displays missing 'n' option as required parameters", func() {
			session := runCommand("generate")

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

			session := runCommand("generate", "-n", "my-value", "-t", "value")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("test error"))
		})
	})
})

func setupGenerateServer(keyType, name, generatedValue, params string, overwrite bool) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", "/api/v1/data"),
			VerifyJSON(fmt.Sprintf(generateRequestJSON, keyType, name, params, overwrite)),
			RespondWith(http.StatusOK, fmt.Sprintf(generateResponseJSON, keyType, name, generatedValue)),
		),
	)
}

func setupGenerateServerWithMetadata(keyType, name, generatedValue, params string, overwrite bool, metadata string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", "/api/v1/data"),
			VerifyJSON(fmt.Sprintf(generateRequestJSONWithMetadata, keyType, name, params, overwrite, metadata)),
			RespondWith(http.StatusOK, fmt.Sprintf(generateResponseJSONWithMetadata, keyType, name, generatedValue, metadata)),
		),
	)
}

func setupGenerateServerWithValue(keyType, name, generatedValue, params, value string, overwrite bool) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", "/api/v1/data"),
			VerifyJSON(fmt.Sprintf(generateWithValueRequestJSON, keyType, name, params, overwrite, value)),
			RespondWith(http.StatusOK, fmt.Sprintf(generateResponseJSON, keyType, name, generatedValue)),
		),
	)
}

func setupGenerateServerWithDurationOverriddenFlag(keyType, name, generatedValue, params string, overwrite bool) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", "/api/v1/data"),
			VerifyJSON(fmt.Sprintf(generateRequestJSON, keyType, name, params, overwrite)),
			RespondWith(http.StatusOK, fmt.Sprintf(generateResponseJSONWithDurationOverriddenFlag, keyType, name, generatedValue)),
		),
	)
}
