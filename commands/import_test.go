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

		It("prevents misconfigured yaml from creating a request", func() {
			SetupOverwritePutValueServer("/director/deployment/blobstore1", "password", "test_password_1", true)
			SetupOverwritePutValueServer("/director/deployment/blobstore2", "password", "test_password_2", true)

			session := runCommand("import", "-f", "../test/test_password_import_file.yml")

			Eventually(session).Should(Exit(0))
			Eventually(session.Err).Should(Say(`Interface conversion error`))
		})
	})

	Describe("importing a file with value credentials", func() {
		It("sets the value credentials", func() {
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

		It("prevents misconfigured yaml from creating a request", func() {
			SetupOverwritePutValueServer("/director/deployment/blobstore3", "value", "test_value_1", true)
			SetupOverwritePutValueServer("/director/deployment/blobstore4", "value", "test_value_2", true)

			session := runCommand("import", "-f", "../test/test_value_import_file.yml")

			Eventually(session).Should(Exit(0))
			Eventually(session.Err).Should(Say(`Interface conversion error`))
		})
	})

	Describe("importing a file with certificate credentials", func() {
		It("sets the password credentials", func() {
			SetupPutCertificateServer("/director/deployment/bosh-ca-cert",
				`ca-certificate`,
				`certificate`,
				`private-key`)
			SetupPutCertificateWithCaNameServer("/director/deployment/bosh-cert-secondary",
				"/dan-cert",
				`certificate`,
				`private-key`)

			session := runCommand("import", "-f", "../test/test_certificate_import_file.yml")

			Eventually(session).Should(Exit(0))

			Eventually(session.Out).Should(Say(`name: /director/deployment/bosh-ca-cert
type: certificate
value:
  ca: ca-certificate
  certificate: certificate
  private_key: private-key`))
			Eventually(session.Out).Should(Say(`name: /director/deployment/bosh-cert-secondary
type: certificate
value:
  ca: known-ca-value
  certificate: certificate
  private_key: private-key`))
		})

		It("prevents misconfigured yaml from creating a request", func() {
			SetupPutCertificateServer("/director/deployment/bosh-ca-cert",
				`ca-certificate`,
				`certificate`,
				`private-key`)
			SetupPutCertificateWithCaNameServer("/director/deployment/bosh-cert-secondary",
				"/dan-cert",
				`certificate`,
				`private-key`)

			session := runCommand("import", "-f", "../test/test_certificate_import_file.yml")

			Eventually(session).Should(Exit(0))
			Eventually(session.Err).Should(Say(`'' expected a map, got 'string'`))
		})
	})

	Describe("importing a file with rsa credentials", func() {
		It("sets the rsa credentials", func() {
			SetupPutRsaSshServer("/director/deployment/rsa1", "rsa", "public-key", "private-key", true)

			session := runCommand("import", "-f", "../test/test_rsa_import_file.yml")

			Eventually(session).Should(Exit(0))

			Eventually(session.Out).Should(Say(`name: /director/deployment/rsa1
type: rsa
value:
  private_key: private-key
  public_key: public-key`))
		})

		It("prevents misconfigured yaml from creating a request", func() {
			SetupPutRsaSshServer("/director/deployment/rsa1", "rsa", "public-key", "private-key", true)

			session := runCommand("import", "-f", "../test/test_rsa_import_file.yml")

			Eventually(session).Should(Exit(0))
			Eventually(session.Err).Should(Say(`'' expected a map, got 'string'`))
		})
	})

	Describe("importing a file with ssh credentials", func() {
		It("sets the ssh credentials", func() {
			SetupPutRsaSshServer("/director/deployment/ssh1", "ssh", "ssh-public-key", "private-key", true)

			session := runCommand("import", "-f", "../test/test_ssh_import_file.yml")

			Eventually(session).Should(Exit(0))

			Eventually(session.Out).Should(Say(`name: /director/deployment/ssh1
type: ssh
value:
  private_key: private-key
  public_key: ssh-public-key`))
		})

		It("prevents misconfigured yaml from creating a request", func() {
			SetupPutRsaSshServer("/director/deployment/ssh1", "ssh", "ssh-public-key", "private-key", true)

			session := runCommand("import", "-f", "../test/test_ssh_import_file.yml")

			Eventually(session).Should(Exit(0))
			Eventually(session.Err).Should(Say(`'' expected a map, got 'string'`))
		})
	})

	Describe("importing a file with mixed credentials", func() {
		It("sets the all credentials", func() {
			SetupOverwritePutValueServer("/director/deployment/blobstore - agent", "password", "gx4ll8193j5rw0wljgqo", true)
			SetupOverwritePutValueServer("/director/deployment/blobstore - director", "value", "y14ck84ef51dnchgk4kp", true)
			SetupPutCertificateServer("/director/deployment/bosh-ca",
				`ca-certificate`,
				`certificate`,
				`private-key`)
			SetupPutCertificateWithCaNameServer("/director/deployment/bosh-cert",
				"/dan-cert",
				`certificate`,
				`private-key`)
			SetupPutRsaSshServer("/director/deployment/rsa", "rsa", "public-key", "private-key", true)
			SetupPutRsaSshServer("/director/deployment/ssh", "ssh", "ssh-public-key", "private-key", true)

			session := runCommand("import", "-f", "../test/test_import_file.yml")

			Eventually(session).Should(Exit(0))

			Eventually(session.Out).Should(Say(`name: /director/deployment/blobstore - agent
type: password
value: gx4ll8193j5rw0wljgqo`))
			Eventually(session.Out).Should(Say(`name: /director/deployment/blobstore - director
type: value
value: y14ck84ef51dnchgk4kp`))
			Eventually(session.Out).Should(Say(`name: /director/deployment/bosh-ca
type: certificate
value:
  ca: ca-certificate
  certificate: certificate
  private_key: private-key`))
			Eventually(session.Out).Should(Say(`name: /director/deployment/bosh-cert
type: certificate
value:
  ca: known-ca-value
  certificate: certificate
  private_key: private-key`))
			Eventually(session.Out).Should(Say(`name: /director/deployment/rsa
type: rsa
value:
  private_key: private-key
  public_key: public-key`))
			Eventually(session.Out).Should(Say(`name: /director/deployment/ssh
type: ssh
value:
  private_key: private-key
  public_key: ssh-public-key`))
		})
	})
})
