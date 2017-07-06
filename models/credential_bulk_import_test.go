package models_test

import (
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/mitchellh/mapstructure"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CredentialBulkImport", func() {
	Describe("readBytes()", func() {
		It("parses YAML", func() {
			var credentialBulkImport models.CredentialBulkImport
			err := credentialBulkImport.ReadBytes(
				[]byte(
					`credentials:
- name: /director/deployment/blobstore - agent
  type: password
  value: gx4ll8193j5rw0wljgqo
- name: /director/deployment/blobstore - director
  type: value
  value: y14ck84ef51dnchgk4kp
- name: /director/deployment/bosh-ca
  type: certificate
  value:
    ca: |
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
    certificate: |
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      ...
      -----END RSA PRIVATE KEY-----
- name: /director/deployment/bosh-cert
  type: certificate
  value:
    ca_name: /dan-cert
    certificate: |
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      ...
      -----END RSA PRIVATE KEY-----`))

			Expect(err).To(BeNil())
			Expect(len(credentialBulkImport.Credentials)).To(Equal(4))
			Expect(credentialBulkImport.Credentials[0].Name).To(Equal("/director/deployment/blobstore - agent"))
			Expect(credentialBulkImport.Credentials[1].Name).To(Equal("/director/deployment/blobstore - director"))
			Expect(credentialBulkImport.Credentials[2].Name).To(Equal("/director/deployment/bosh-ca"))
			Expect(credentialBulkImport.Credentials[3].Name).To(Equal("/director/deployment/bosh-cert"))
			Expect(credentialBulkImport.Credentials[0].Type).To(Equal("password"))
			Expect(credentialBulkImport.Credentials[1].Type).To(Equal("value"))
			Expect(credentialBulkImport.Credentials[2].Type).To(Equal("certificate"))
			Expect(credentialBulkImport.Credentials[3].Type).To(Equal("certificate"))
			Expect(credentialBulkImport.Credentials[0].Value.(string)).To(Equal("gx4ll8193j5rw0wljgqo"))
			Expect(credentialBulkImport.Credentials[1].Value.(string)).To(Equal("y14ck84ef51dnchgk4kp"))

			var certificate1 models.Certificate
			err = mapstructure.Decode(credentialBulkImport.Credentials[2].Value, &certificate1)
			Expect(err).To(BeNil())
			Expect(certificate1.Ca).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
			Expect(certificate1.Certificate).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
			Expect(certificate1.PrivateKey).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))
			Expect(certificate1.CaName).To(Equal(""))

			var certificate2 models.Certificate
			err = mapstructure.Decode(credentialBulkImport.Credentials[3].Value, &certificate2)
			Expect(err).To(BeNil())
			Expect(certificate2.Ca).To(Equal(""))
			Expect(certificate2.Certificate).To(ContainSubstring(`-----BEGIN CERTIFICATE-----`))
			Expect(certificate2.PrivateKey).To(ContainSubstring(`-----BEGIN RSA PRIVATE KEY-----`))
			Expect(certificate2.CaName).To(Equal("/dan-cert"))
		})
	})

	Describe("readFile()", func() {
		It("parses YAML from an input file", func() {
			var credentialBulkImport models.CredentialBulkImport
			err := credentialBulkImport.ReadFile("../test/test_import_file.yml")

			Expect(err).To(BeNil())
			Expect(len(credentialBulkImport.Credentials)).To(Equal(4))
			Expect(credentialBulkImport.Credentials[0].Name).To(Equal("/director/deployment/blobstore - agent"))
			Expect(credentialBulkImport.Credentials[1].Name).To(Equal("/director/deployment/blobstore - director"))
			Expect(credentialBulkImport.Credentials[2].Name).To(Equal("/director/deployment/bosh-ca"))
			Expect(credentialBulkImport.Credentials[3].Name).To(Equal("/director/deployment/bosh-cert"))
			Expect(credentialBulkImport.Credentials[0].Type).To(Equal("password"))
			Expect(credentialBulkImport.Credentials[1].Type).To(Equal("value"))
			Expect(credentialBulkImport.Credentials[2].Type).To(Equal("certificate"))
			Expect(credentialBulkImport.Credentials[3].Type).To(Equal("certificate"))
			Expect(credentialBulkImport.Credentials[0].Value.(string)).To(Equal("gx4ll8193j5rw0wljgqo"))
			Expect(credentialBulkImport.Credentials[1].Value.(string)).To(Equal("y14ck84ef51dnchgk4kp"))

			var certificate1 models.Certificate
			err = mapstructure.Decode(credentialBulkImport.Credentials[2].Value, &certificate1)
			Expect(err).To(BeNil())
			Expect(certificate1.Ca).To(ContainSubstring(`ca-certificate`))
			Expect(certificate1.Certificate).To(ContainSubstring(`certificate`))
			Expect(certificate1.PrivateKey).To(ContainSubstring(`private-key`))
			Expect(certificate1.CaName).To(Equal(""))

			var certificate2 models.Certificate
			err = mapstructure.Decode(credentialBulkImport.Credentials[3].Value, &certificate2)
			Expect(err).To(BeNil())
			Expect(certificate2.Ca).To(Equal(""))
			Expect(certificate2.Certificate).To(ContainSubstring(`certificate`))
			Expect(certificate2.PrivateKey).To(ContainSubstring(`private-key`))
			Expect(certificate2.CaName).To(Equal("/dan-cert"))
		})
	})
})
