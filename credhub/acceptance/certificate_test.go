package acceptance_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

var _ = Describe("Certificate", func() {
	var ch *CredHub

	BeforeEach(func() {
		var err error
		ch, err = New("https://localhost:9000",
			SkipTLSValidation(),
			AuthBuilder(uaa.PasswordGrantBuilder("credhub_cli", "", "credhub", "password")))

		Expect(err).ToNot(HaveOccurred())
	})

	Specify("certificate lifecycle", func() {

		name := fmt.Sprintf("/acceptance/certificate/%v", time.Now().UnixNano())
		generateCert := generate.Certificate{
			CommonName: "example.com",
			SelfSign:   true,
		}

		setCert := values.Certificate{
			Ca:          "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
			Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
			PrivateKey:  "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
		}

		By("generate a certificate with path " + name)
		certificate, err := ch.GenerateCertificate(name, generateCert, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value.Certificate).ToNot(BeEmpty())
		Expect(certificate.Value.PrivateKey).ToNot(BeEmpty())
		firstGeneratedCertificate := certificate.Value

		By("generate the certificate again without overwrite returns same certificate")
		certificate, err = ch.GenerateCertificate(name, generateCert, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).To(Equal(firstGeneratedCertificate))

		By("setting the certificate again without overwrite returns same certificate")
		certificate, err = ch.SetCertificate(name, setCert, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).To(Equal(firstGeneratedCertificate))

		By("overwriting the certificate with generate")
		certificate, err = ch.GenerateCertificate(name, generateCert, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).ToNot(Equal(firstGeneratedCertificate))

		By("overwriting the certificate with set")
		certificate, err = ch.SetCertificate(name, setCert, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).To(Equal(setCert))

		By("getting the certificate")
		certificate, err = ch.GetCertificate(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).To(Equal(setCert))

		By("deleting the certificate")
		err = ch.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = ch.GetCertificate(name)
		Expect(err).To(HaveOccurred())
	})
})
