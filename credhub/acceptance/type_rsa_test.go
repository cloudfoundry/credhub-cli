package acceptance_test

import (
	"crypto/x509"
	"encoding/pem"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

var _ = Describe("RSA Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("some-rsa")
		opts := generate.RSA{KeyLength: 2048}

		By("generate rsa keys with path " + name)
		generatedRSA, err := credhubClient.GenerateRSA(name, opts, false)
		Expect(err).ToNot(HaveOccurred())
		block, _ := pem.Decode([]byte(generatedRSA.Value.PrivateKey))
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(privateKey.N.BitLen()).To(Equal(2048))

		By("generate the rsa keys again without overwrite returns same rsa")
		rsa, err := credhubClient.GenerateRSA(name, opts, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsa).To(Equal(generatedRSA))

		By("setting the rsa keys again without overwrite returns the same")
		newRSA := values.RSA{PrivateKey: "private key", PublicKey: "public key"}
		rsa, err = credhubClient.SetRSA(name, newRSA, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsa).To(Equal(generatedRSA))

		By("overwriting with generate")
		rsa, err = credhubClient.GenerateRSA(name, opts, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsa).ToNot(Equal(generatedRSA))

		By("overwriting with set")
		rsa, err = credhubClient.SetRSA(name, newRSA, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsa.Value).To(Equal(newRSA))

		By("getting the rsa credential")
		rsa, err = credhubClient.GetRSA(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsa.Value).To(Equal(newRSA))

		By("deleting the rsa credential")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetRSA(name)
		Expect(err).To(HaveOccurred())
	})
})
