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
)

var _ = Describe("Generate", func() {
	It("without parameters", func() {
		doValueDefaultTypeOptionTest(`{}`)
	})

	Describe("with a variety of value parameters", func() {
		It("including length", func() {
			doValueOptionTest(`{"length":42}`, "-l", "42")
		})

		It("excluding upper case", func() {
			doValueOptionTest(`{"exclude_upper":true}`, "--exclude-upper")
		})

		It("excluding lower case", func() {
			doValueOptionTest(`{"exclude_lower":true}`, "--exclude-lower")
		})

		It("excluding special characters", func() {
			doValueOptionTest(`{"exclude_special":true}`, "--exclude-special")
		})

		It("excluding numbers", func() {
			doValueOptionTest(`{"exclude_number":true}`, "--exclude-number")
		})
	})

	Describe("with a variety of password parameters", func() {
		It("including length", func() {
			doPasswordOptionTest(`{"length":42}`, "-l", "42")
		})

		It("excluding upper case", func() {
			doPasswordOptionTest(`{"exclude_upper":true}`, "--exclude-upper")
		})

		It("excluding lower case", func() {
			doPasswordOptionTest(`{"exclude_lower":true}`, "--exclude-lower")
		})

		It("excluding special characters", func() {
			doPasswordOptionTest(`{"exclude_special":true}`, "--exclude-special")
		})

		It("excluding numbers", func() {
			doPasswordOptionTest(`{"exclude_number":true}`, "--exclude-number")
		})
	})

	Describe("with a variety of certificate parameters", func() {
		It("including common name", func() {
			doCertificateOptionTest(`{"common_name":"common.name.io"}`, "--common-name", "common.name.io")
		})

		It("including organization", func() {
			doCertificateOptionTest(`{"organization":"organization.io"}`, "--organization", "organization.io")
		})

		It("including organization unit", func() {
			doCertificateOptionTest(`{"organization_unit":"My Unit"}`, "--organization-unit", "My Unit")
		})

		It("including locality", func() {
			doCertificateOptionTest(`{"locality":"My Locality"}`, "--locality", "My Locality")
		})

		It("including state", func() {
			doCertificateOptionTest(`{"state":"My State"}`, "--state", "My State")
		})

		It("including country", func() {
			doCertificateOptionTest(`{"country":"My Country"}`, "--country", "My Country")
		})

		It("including multiple alternative names", func() {
			doCertificateOptionTest(`{"alternative_names": [ "Alt1", "Alt2" ]}`, "--alternative-name", "Alt1", "--alternative-name", "Alt2")
		})

		It("including key length", func() {
			doCertificateOptionTest(`{"key_length":2048}`, "--key-length", "2048")
		})

		It("including duration", func() {
			doCertificateOptionTest(`{"duration":1000}`, "--duration", "1000")
		})

		It("including certificate authority", func() {
			doCertificateOptionTest(`{"ca":"my_ca"}`, "--ca", "my_ca")
		})
	})

	Describe("Help", func() {
		It("displays help", func() {
			session := runCommand("generate", "-h")

			Eventually(session).Should(Exit(1))
			Expect(session.Err).To(Say("generate"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("length"))
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
			RespondWith(http.StatusOK, fmt.Sprintf(SECRET_STRING_RESPONSE_JSON, "value", value)),
		),
	)
}

func generateRequestJson(secretType string, params string) string {
	return fmt.Sprintf(GENERATE_REQUEST_JSON, secretType, params)
}

func generateDefaultTypeRequestJson(params string) string {
	return fmt.Sprintf(GENERATE_DEFAULT_TYPE_REQUEST_JSON, params)
}

func doValueDefaultTypeOptionTest(optionJson string, options ...string) {
	setupPostServer("my-password", "potatoes", generateDefaultTypeRequestJson(optionJson))

	doTest([]string{"generate", "-n", "my-password"}, options...)
}

func doValueOptionTest(optionJson string, options ...string) {
	setupPostServer("my-value", "potatoes", generateRequestJson("value", optionJson))

	doTest([]string{"generate", "-n", "my-value", "-t", "value"}, options...)
}

func doPasswordOptionTest(optionJson string, options ...string) {
	setupPostServer("my-password", "potatoes", generateRequestJson("password", optionJson))

	doTest([]string{"generate", "-n", "my-password", "-t", "password"}, options...)
}

func doCertificateOptionTest(optionJson string, options ...string) {
	setupPostServer("my-secret", "potatoes", generateRequestJson("certificate", optionJson))

	doTest([]string{"generate", "-n", "my-secret", "-t", "certificate"}, options...)
}

func doTest(leftOpts []string, options ...string) {
	stuff := append(leftOpts, options...)
	session := runCommand(stuff...)

	Eventually(session).Should(Exit(0))
}
