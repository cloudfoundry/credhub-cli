package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

var _ = Describe("Certificate", func() {
	Specify("certificate lifecycle", func() {
		name := testCredentialPath("some-certificate")

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
		certificate, err := credhubClient.GenerateCertificate(name, generateCert, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value.Certificate).ToNot(BeEmpty())
		Expect(certificate.Value.PrivateKey).ToNot(BeEmpty())
		firstGeneratedCertificate := certificate.Value

		By("generate the certificate again without overwrite returns same certificate")
		certificate, err = credhubClient.GenerateCertificate(name, generateCert, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).To(Equal(firstGeneratedCertificate))

		By("setting the certificate again without overwrite returns same certificate")
		certificate, err = credhubClient.SetCertificate(name, setCert, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).To(Equal(firstGeneratedCertificate))

		By("overwriting the certificate with generate")
		certificate, err = credhubClient.GenerateCertificate(name, generateCert, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).ToNot(Equal(firstGeneratedCertificate))

		By("overwriting the certificate with set")
		certificate, err = credhubClient.SetCertificate(name, setCert, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).To(Equal(setCert))

		By("getting the certificate")
		certificate, err = credhubClient.GetCertificate(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(certificate.Value).To(Equal(setCert))

		By("deleting the certificate")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetCertificate(name)
		Expect(err).To(HaveOccurred())
	})
})
