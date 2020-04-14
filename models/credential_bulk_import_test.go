package models_test

import (
	"code.cloudfoundry.org/credhub-cli/errors"
	"code.cloudfoundry.org/credhub-cli/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CredentialBulkImport", func() {
	Describe("ReadFile()", func() {
		It("parses YAML", func() {
			var credentialBulkImport models.CredentialBulkImport
			err := credentialBulkImport.ReadFile("../test/test_import_file.yml", false)

			Expect(err).To(BeNil())
			Expect(len(credentialBulkImport.Credentials)).To(Equal(7))

			expectedPassword := make(map[string]interface{})
			expectedPassword["name"] = "/test/password"
			expectedPassword["type"] = "password"
			expectedPassword["value"] = "test-password-value"
			expectedPassword["overwrite"] = true

			Expect(credentialBulkImport.Credentials[0]).To(Equal(expectedPassword))

			expectedValue := make(map[string]interface{})
			expectedValue["name"] = "/test/value"
			expectedValue["type"] = "value"
			expectedValue["value"] = "test-value"
			expectedValue["metadata"] = map[string]interface{}{
				"some": "thing",
				"nested": map[string]interface{}{
					"with": "value",
				},
			}
			expectedValue["overwrite"] = true

			Expect(credentialBulkImport.Credentials[1]).To(Equal(expectedValue))

			expectedCertificate := make(map[string]interface{})
			expectedCertificate["name"] = "/test/certificate"
			expectedCertificate["type"] = "certificate"
			expectedCertificate["value"] = map[string]interface{}{
				"certificate": "certificate",
				"private_key": "private-key",
				"ca":          "ca-certificate",
			}
			expectedCertificate["metadata"] = map[string]interface{}{
				"certs": "metadata",
			}
			expectedCertificate["overwrite"] = true

			Expect(credentialBulkImport.Credentials[2]).To(Equal(expectedCertificate))

			expectedRsa := make(map[string]interface{})
			expectedRsa["name"] = "/test/rsa"
			expectedRsa["type"] = "rsa"
			expectedRsa["value"] = map[string]interface{}{
				"public_key":  "public-key",
				"private_key": "private-key",
			}
			expectedRsa["overwrite"] = true

			Expect(credentialBulkImport.Credentials[3]).To(Equal(expectedRsa))

			expectedSsh := make(map[string]interface{})
			expectedSsh["name"] = "/test/ssh"
			expectedSsh["type"] = "ssh"
			expectedSsh["value"] = map[string]interface{}{
				"public_key":  "ssh-public-key",
				"private_key": "private-key",
			}
			expectedSsh["overwrite"] = true

			Expect(credentialBulkImport.Credentials[4]).To(Equal(expectedSsh))

			expectedUser := make(map[string]interface{})
			expectedUser["name"] = "/test/user"
			expectedUser["type"] = "user"
			expectedUser["value"] = map[string]interface{}{
				"username": "covfefe",
				"password": "test-user-password",
			}
			expectedUser["overwrite"] = true

			Expect(credentialBulkImport.Credentials[5]).To(Equal(expectedUser))

			expectedJSON := make(map[string]interface{})
			expectedJSON["name"] = "/test/json"
			expectedJSON["type"] = "json"
			expectedJSON["value"] = map[string]interface{}{
				"arbitrary_object": map[string]interface{}{
					"nested_array": []interface{}{
						"array_val1",
						map[string]interface{}{"array_object_subvalue": "covfefe"},
					},
				},
				"1":    "key is not a string",
				"3.14": "pi",
				"true": "key is a bool",
			}
			expectedJSON["overwrite"] = true

			Expect(credentialBulkImport.Credentials[6]).To(Equal(expectedJSON))
		})
		It("parses JSON", func() {
			var credentialBulkImport models.CredentialBulkImport
			err := credentialBulkImport.ReadFile("../test/test_import_file.json", true)

			Expect(err).To(BeNil())
			Expect(len(credentialBulkImport.Credentials)).To(Equal(7))

			expectedPassword := make(map[string]interface{})
			expectedPassword["name"] = "/test/password"
			expectedPassword["type"] = "password"
			expectedPassword["value"] = "test-password-value"
			expectedPassword["overwrite"] = true

			Expect(credentialBulkImport.Credentials[0]).To(Equal(expectedPassword))

			expectedValue := make(map[string]interface{})
			expectedValue["name"] = "/test/value"
			expectedValue["type"] = "value"
			expectedValue["value"] = "test-value"
			expectedValue["metadata"] = map[string]interface{}{
				"some": "thing",
				"nested": map[string]interface{}{
					"with": "value",
				},
			}
			expectedValue["overwrite"] = true

			Expect(credentialBulkImport.Credentials[1]).To(Equal(expectedValue))

			expectedCertificate := make(map[string]interface{})
			expectedCertificate["name"] = "/test/certificate"
			expectedCertificate["type"] = "certificate"
			expectedCertificate["value"] = map[string]interface{}{
				"ca":          "ca-certificate",
				"certificate": "certificate",
				"private_key": "private-key",
			}
			expectedCertificate["metadata"] = map[string]interface{}{
				"certs": "metadata",
			}
			expectedCertificate["overwrite"] = true

			Expect(credentialBulkImport.Credentials[2]).To(Equal(expectedCertificate))

			expectedRsa := make(map[string]interface{})
			expectedRsa["name"] = "/test/rsa"
			expectedRsa["type"] = "rsa"
			expectedRsa["value"] = map[string]interface{}{
				"public_key":  "public-key",
				"private_key": "private-key",
			}
			expectedCertificate["metadata"] = map[string]interface{}{
				"certs": "metadata",
			}
			expectedRsa["overwrite"] = true

			Expect(credentialBulkImport.Credentials[3]).To(Equal(expectedRsa))

			expectedSsh := make(map[string]interface{})
			expectedSsh["name"] = "/test/ssh"
			expectedSsh["type"] = "ssh"
			expectedSsh["value"] = map[string]interface{}{
				"public_key":  "ssh-public-key",
				"private_key": "private-key",
			}
			expectedSsh["overwrite"] = true

			Expect(credentialBulkImport.Credentials[4]).To(Equal(expectedSsh))

			expectedUser := make(map[string]interface{})
			expectedUser["name"] = "/test/user"
			expectedUser["type"] = "user"
			expectedUser["value"] = map[string]interface{}{
				"username": "covfefe",
				"password": "test-user-password",
			}
			expectedUser["overwrite"] = true

			Expect(credentialBulkImport.Credentials[5]).To(Equal(expectedUser))

			expectedJSON := make(map[string]interface{})
			expectedJSON["name"] = "/test/json"
			expectedJSON["type"] = "json"
			expectedJSON["value"] = map[string]interface{}{
				"arbitrary_object": map[string]interface{}{
					"nested_array": []interface{}{
						"array_val1",
						map[string]interface{}{"array_object_subvalue": "covfefe"},
					},
				},
				"1":    "key is not a string",
				"3.14": "pi",
				"true": "key is a bool",
			}
			expectedJSON["overwrite"] = true

			Expect(credentialBulkImport.Credentials[6]).To(Equal(expectedJSON))
		})
	})
	Describe("formatting", func() {
		var credentialBulkImport *models.CredentialBulkImport
		BeforeEach(func() {
			credentialBulkImport = &models.CredentialBulkImport{}
		})
		Context("when the import file is of type json", func() {
			Context("when first line is credentials tag", func() {
				It("does not return an error", func() {
					credentials := `{
  "credentials": [
    {
      "name": "/test/password",
      "type": "password",
      "value": "test-password-value"
    }
  ]
}`
					err := credentialBulkImport.ReadBytes([]byte(credentials), true)
					Expect(err).To(BeNil())
				})
			})

			Context("when first line is not credentials tag", func() {
				It("raises an error", func() {
					credentials := `{
  "not-credentials": [
    {
      "name": "/test/password",
      "type": "password",
      "value": "test-password-value"
    }
  ]
}`
					err := credentialBulkImport.ReadBytes([]byte(credentials), true)
					Expect(err).To(Equal(errors.NewNoCredentialsTagError()))
				})
			})

			Context("when yaml is incorrect", func() {
				It("raises an error", func() {
					credentials := `{
  "credentials": [
    {
      "name: "/test/password",
      "type": "password"
      5
  ]
}`
					err := credentialBulkImport.ReadBytes([]byte(credentials), true)
					Expect(err).To(Equal(errors.NewInvalidImportJSONError()))
				})
			})
		})

		Context("when the import file is of type yaml", func() {
			Context("when first line is credentials tag", func() {
				It("does not return an error", func() {
					credentials := `credentials:
- name: /test/password
  type: password
  value: test-password-value`
					err := credentialBulkImport.ReadBytes([]byte(credentials), false)
					Expect(err).To(BeNil())
				})
			})

			Context("when first line is credentials tag with trailing white space", func() {
				It("does not return an error", func() {
					credentials := "credentials:   \n" +
						`- name: /test/password
  type: password
  value: test-password-value`
					err := credentialBulkImport.ReadBytes([]byte(credentials), false)
					Expect(err).To(BeNil())
				})
			})

			Context("when first line is not credentials tag", func() {
				It("raises an error", func() {
					credentials := `not-credentials:
- name: /test/password
  type: password
  value: test-password-value`
					err := credentialBulkImport.ReadBytes([]byte(credentials), false)
					Expect(err).To(Equal(errors.NewNoCredentialsTagError()))
				})
			})

			Context("when yaml is incorrect", func() {
				It("raises an error", func() {
					credentials := `credentials:
1
2
			`
					err := credentialBulkImport.ReadBytes([]byte(credentials), false)
					Expect(err).To(Equal(errors.NewInvalidImportYamlError()))
				})
			})
		})
	})
})
