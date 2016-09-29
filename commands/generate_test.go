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

			doTest([]string{"generate", "-n", "my-password"})
		})

		It("prints the generated password secret", func() {
			setupPostServer("my-password", "potatoes", generateDefaultTypeRequestJson(`{}`, true))

			result := doTest([]string{"generate", "-n", "my-password"})

			Expect(result.Out).To(Say(responseMyPasswordPotatoes))
		})
	})

	Describe("with a variety of password parameters", func() {
		It("with with no-overwrite", func() {
			doPasswordOptionTest(`{}`, false, "--no-overwrite")
		})

		It("including length", func() {
			doPasswordOptionTest(`{"length":42}`, true, "-l", "42")
		})

		It("excluding upper case", func() {
			doPasswordOptionTest(`{"exclude_upper":true}`, true, "--exclude-upper")
		})

		It("excluding lower case", func() {
			doPasswordOptionTest(`{"exclude_lower":true}`, true, "--exclude-lower")
		})

		It("excluding special characters", func() {
			doPasswordOptionTest(`{"exclude_special":true}`, true, "--exclude-special")
		})

		It("excluding numbers", func() {
			doPasswordOptionTest(`{"exclude_number":true}`, true, "--exclude-number")
		})
	})

	Describe("with a variety of certificate parameters", func() {
		It("including common name", func() {
			doCertificateOptionTest(`{"common_name":"common.name.io"}`, true, "--common-name", "common.name.io")
		})

		It("including common name with no-overwrite", func() {
			doCertificateOptionTest(`{"common_name":"common.name.io"}`, false, "--common-name", "common.name.io", "--no-overwrite")
		})

		It("including organization", func() {
			doCertificateOptionTest(`{"organization":"organization.io"}`, true, "--organization", "organization.io")
		})

		It("including organization unit", func() {
			doCertificateOptionTest(`{"organization_unit":"My Unit"}`, true, "--organization-unit", "My Unit")
		})

		It("including locality", func() {
			doCertificateOptionTest(`{"locality":"My Locality"}`, true, "--locality", "My Locality")
		})

		It("including state", func() {
			doCertificateOptionTest(`{"state":"My State"}`, true, "--state", "My State")
		})

		It("including country", func() {
			doCertificateOptionTest(`{"country":"My Country"}`, true, "--country", "My Country")
		})

		It("including multiple alternative names", func() {
			doCertificateOptionTest(`{"alternative_names": [ "Alt1", "Alt2" ]}`, true, "--alternative-name", "Alt1", "--alternative-name", "Alt2")
		})

		It("including key length", func() {
			doCertificateOptionTest(`{"key_length":2048}`, true, "--key-length", "2048")
		})

		It("including duration", func() {
			doCertificateOptionTest(`{"duration":1000}`, true, "--duration", "1000")
		})

		It("including certificate authority", func() {
			doCertificateOptionTest(`{"ca":"my_ca"}`, true, "--ca", "my_ca")
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

func doPasswordOptionTest(optionJson string, overwrite bool, options ...string) {
	setupPostServer("my-password", "potatoes", generateRequestJson("password", optionJson, overwrite))

	doTest([]string{"generate", "-n", "my-password", "-t", "password"}, options...)
}

func doCertificateOptionTest(optionJson string, overwrite bool, options ...string) {
	setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", optionJson, overwrite))

	doTest([]string{"generate", "-n", "my-secret", "-t", "certificate"}, options...)
}

func doTest(leftOpts []string, options ...string) *Session {
	stuff := append(leftOpts, options...)
	session := runCommand(stuff...)

	Eventually(session).Should(Exit(0))
	return session
}
