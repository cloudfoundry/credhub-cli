package commands_test

import (
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Token", func() {
	Context("when the config file has a token", func() {
		BeforeEach(func() {
			cfg := config.ReadConfig()
			cfg.AccessToken = "FAKETOKEN"
			config.WriteConfig(cfg)
		})

		It("displays the token with --token", func() {
			session := runCommand("--token")

			Eventually(session).Should(Exit(0))
			sout := string(session.Out.Contents())
			Expect(sout).To(ContainSubstring("Bearer FAKETOKEN"))
		})
	})

	Context("when the config file does not have a token", func() {
		BeforeEach(func() {
			cfg := config.ReadConfig()
			cfg.AccessToken = ""
			config.WriteConfig(cfg)
		})

		It("displays nothing", func() {
			session := runCommand("--token")

			Eventually(session).Should(Exit(0))
			sout := string(session.Out.Contents())
			Expect(sout).To(Equal(""))
		})
	})
})
