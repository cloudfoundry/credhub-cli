package acceptance_test

import (
	"crypto/x509"
	"encoding/pem"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

var _ = Describe("SSH Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("some-ssh")
		opts := generate.SSH{KeyLength: 2048}

		By("generate ssh keys with path " + name)
		generatedSSH, err := credhubClient.GenerateSSH(name, opts, false)
		Expect(err).ToNot(HaveOccurred())
		block, _ := pem.Decode([]byte(generatedSSH.Value.PrivateKey))
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(privateKey.N.BitLen()).To(Equal(2048))

		By("generate the ssh keys again without overwrite returns same ssh")
		ssh, err := credhubClient.GenerateSSH(name, opts, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(ssh).To(Equal(generatedSSH))

		By("setting the ssh keys again without overwrite returns the same")
		newSSH := values.SSH{PrivateKey: "private key", PublicKey: "public key"}
		ssh, err = credhubClient.SetSSH(name, newSSH, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(ssh).To(Equal(generatedSSH))

		By("overwriting with generate")
		ssh, err = credhubClient.GenerateSSH(name, opts, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(ssh).ToNot(Equal(generatedSSH))

		By("overwriting with set")
		ssh, err = credhubClient.SetSSH(name, newSSH, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(ssh.Value).To(Equal(newSSH))

	})
})
