package commands_test

import (
	"net/http"

	"fmt"

	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/credhub-cli/commands"
)

var _ = Describe("Generate", func() {
	Describe("Without parameters", func() {
		It("uses default parameters", func() {
			setupPostServer("my-password", "potatoes", generateDefaultTypeRequestJson(`{}`, true))

			session := runCommand("generate", "-n", "my-password")

			Eventually(session).Should(Exit(0))
		})

		It("prints the generated password secret", func() {
			setupPostServer("my-password", "potatoes", generateDefaultTypeRequestJson(`{}`, true))

			session := runCommand("generate", "-n", "my-password")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(responseMyPasswordPotatoes))
		})
	})

	Describe("with a variety of password parameters", func() {
		It("with with no-overwrite", func() {
			setupPostServer("my-password", "potatoes", generateRequestJson("password", `{}`, false))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupPostServer("my-password", "potatoes", generateRequestJson("password", `{"length":42}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "-l", "42")
			Eventually(session).Should(Exit(0))
		})

		It("excluding upper case", func() {
			setupPostServer("my-password", "potatoes", generateRequestJson("password", `{"exclude_upper":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-upper")
			Eventually(session).Should(Exit(0))
		})

		It("excluding lower case", func() {
			setupPostServer("my-password", "potatoes", generateRequestJson("password", `{"exclude_lower":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-lower")
			Eventually(session).Should(Exit(0))
		})

		It("excluding special characters", func() {
			setupPostServer("my-password", "potatoes", generateRequestJson("password", `{"exclude_special":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-special")
			Eventually(session).Should(Exit(0))
		})

		It("excluding numbers", func() {
			setupPostServer("my-password", "potatoes", generateRequestJson("password", `{"exclude_number":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-number")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of certificate parameters", func() {
		It("including common name", func() {
			setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", `{"common_name":"common.name.io"}`, true))
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io")
			Eventually(session).Should(Exit(0))
		})

		It("including common name with no-overwrite", func() {
			setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", `{"common_name":"common.name.io"}`, false))
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including organization", func() {
			setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", `{"organization":"organization.io"}`, true))
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--organization", "organization.io")
			Eventually(session).Should(Exit(0))
		})

		It("including organization unit", func() {
			setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", `{"organization_unit":"My Unit"}`, true))
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--organization-unit", "My Unit")
			Eventually(session).Should(Exit(0))
		})

		It("including locality", func() {
			setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", `{"locality":"My Locality"}`, true))
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--locality", "My Locality")
			Eventually(session).Should(Exit(0))
		})

		It("including state", func() {
			setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", `{"state":"My State"}`, true))
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--state", "My State")
			Eventually(session).Should(Exit(0))
		})

		It("including country", func() {
			setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", `{"country":"My Country"}`, true))
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--country", "My Country")
			Eventually(session).Should(Exit(0))
		})

		It("including multiple alternative names", func() {
			setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", `{"alternative_names": [ "Alt1", "Alt2" ]}`, true))
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--alternative-name", "Alt1", "--alternative-name", "Alt2")
			Eventually(session).Should(Exit(0))
		})

		It("including key length", func() {
			setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", `{"key_length":2048}`, true))
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--key-length", "2048")
			Eventually(session).Should(Exit(0))
		})

		It("including duration", func() {
			setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", `{"duration":1000}`, true))
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--duration", "1000")
			Eventually(session).Should(Exit(0))
		})

		It("including certificate authority", func() {
			setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", `{"ca":"my_ca"}`, true))
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--ca", "my_ca")
			Eventually(session).Should(Exit(0))
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
				commands.HaveFlag("exclude-special", "S"),
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
				RespondWith(http.StatusBadRequest, `{"error": "you fail."}`),
			)

			session := runCommand("generate", "-n", "my-value")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("you fail."))
		})
	})
})

func setupPostServer(name string, value string, requestJson string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", fmt.Sprintf("/api/v1/data/%s", name)),
			VerifyJSON(requestJson),
			RespondWith(http.StatusOK, fmt.Sprintf(SECRET_STRING_RESPONSE_JSON, "password", value)),
		),
	)
}

func generateRequestJson(secretType string, params string, overwrite bool) string {
	return fmt.Sprintf(GENERATE_SECRET_REQUEST_JSON, secretType, overwrite, params)
}

func generateDefaultTypeRequestJson(params string, overwrite bool) string {
	return fmt.Sprintf(GENERATE_DEFAULT_TYPE_REQUEST_JSON, overwrite, params)
}
