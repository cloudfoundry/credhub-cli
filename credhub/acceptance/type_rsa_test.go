package acceptance_test

import (
	"encoding/pem"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"crypto/x509"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
)

var _ = Describe("RSA Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("some-generatedRSA")
		opts := generate.RSA{KeyLength: 4096}

		By("generate an rsa key with path " + name)
		generatedRSA, err := credhubClient.GenerateRSA(name, opts, false)
		Expect(err).ToNot(HaveOccurred())
		block, _ := pem.Decode([]byte(generatedRSA.Value.PrivateKey))
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(privateKey.N.BitLen()).To(Equal(4096))

		By("generate the rsa again without overwrite returns same rsa")
		rsa, err := credhubClient.GenerateRSA(name, opts, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsa).To(Equal(generatedRSA))

		By("overwriting the user with generate")
		rsa, err = credhubClient.GenerateRSA(name, opts, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsa).ToNot(Equal(generatedRSA))
	})
})
