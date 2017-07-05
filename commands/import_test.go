package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Import", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("get", "-n", "test-credential")

	Describe("importing a file with password credentials", func() {
		It("sets the password credentials", func() {
			SetupOverwritePutValueServer("/director/deployment/blobstore1", "password", "test_password_1", true)
			SetupOverwritePutValueServer("/director/deployment/blobstore2", "password", "test_password_2", true)

			session := runCommand("import", "-f", "../test/test_password_import_file.yml")

			Eventually(session).Should(Exit(0))

			Eventually(session.Out).Should(Say(`name: /director/deployment/blobstore1
type: password
value: test_password_1`))
			Eventually(session.Out).Should(Say(`name: /director/deployment/blobstore2
type: password
value: test_password_2`))
		})
	})

	Describe("importing a file with value credentials", func() {
		It("sets the password credentials", func() {
			SetupOverwritePutValueServer("/director/deployment/blobstore3", "value", "test_value_1", true)
			SetupOverwritePutValueServer("/director/deployment/blobstore4", "value", "test_value_2", true)

			session := runCommand("import", "-f", "../test/test_value_import_file.yml")

			Eventually(session).Should(Exit(0))

			Eventually(session.Out).Should(Say(`name: /director/deployment/blobstore3
type: value
value: test_value_1`))
			Eventually(session.Out).Should(Say(`name: /director/deployment/blobstore4
type: value
value: test_value_2`))
		})
	})

	Describe("importing a file with mixed credentials", func() {
		It("sets the password credentials", func() {
			SetupOverwritePutValueServer("/director/deployment/blobstore - agent", "password", "gx4ll8193j5rw0wljgqo", true)
			SetupOverwritePutValueServer("/director/deployment/blobstore - director", "value", "y14ck84ef51dnchgk4kp", true)

			session := runCommand("import", "-f", "../test/test_import_file.yml")

			Eventually(session).Should(Exit(0))

			Eventually(session.Out).Should(Say(`name: /director/deployment/blobstore - agent
type: password
value: gx4ll8193j5rw0wljgqo`))
			Eventually(session.Out).Should(Say(`name: /director/deployment/blobstore - director
type: value
value: y14ck84ef51dnchgk4kp`))
		})
	})
})
