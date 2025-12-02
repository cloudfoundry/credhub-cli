package models_test

import (
	"encoding/json"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.yaml.in/yaml/v3"
)

var _ = Describe("ExportCredentials", func() {
	credentials := []credentials.Credential{
		{
			Base: credentials.Base{
				Id:               "valueID",
				Name:             "valueName",
				Type:             "value",
				VersionCreatedAt: "valueCreatedAt",
			},
			Value: "test",
		},
		{
			Base: credentials.Base{
				Id:               "passwordID",
				Name:             "passwordName",
				Type:             "password",
				VersionCreatedAt: "passwordCreatedAt",
			},
			Value: "test",
		},
	}

	Describe("outputJSON is set to true", func() {
		It("returns a JSON map with a root credential object", func() {
			exportCreds, err := models.ExportCredentials(credentials, true)
			Expect(err).ToNot(HaveOccurred())

			var v map[string]interface{}
			var mapOfInterfaces []interface{}

			err = json.Unmarshal(exportCreds.Bytes, &v)

			Expect(err).To(BeNil())
			Expect(v["Credentials"]).NotTo(BeNil())
			Expect(v["Credentials"]).To(BeAssignableToTypeOf(mapOfInterfaces))
		})

		It("lists each credential", func() {
			exportCreds, _ := models.ExportCredentials(credentials, true)

			var v map[string]interface{}
			_ = json.Unmarshal(exportCreds.Bytes, &v)

			exportedCredentials := v["Credentials"].([]interface{})

			Expect(exportedCredentials).To(HaveLen(len(credentials)))
		})

		It("includes only a name, type and value in each credential", func() {
			expectedKeys := []string{"Name", "Type", "Value", "Metadata"}
			exportCreds, _ := models.ExportCredentials(credentials, true)

			var v map[string]interface{}
			_ = json.Unmarshal(exportCreds.Bytes, &v)

			exportedCredentials := v["Credentials"].([]interface{})

			for _, credential := range exportedCredentials {
				c := credential.(map[string]interface{})

				for k := range c {
					Expect(expectedKeys).To(ContainElement(k))
				}
			}
		})

		It("produces JSON that can be reimported", func() {
			exportCreds, _ := models.ExportCredentials(credentials, true)
			credImporter := &models.CredentialBulkImport{}

			err := credImporter.ReadBytes(exportCreds.Bytes, true)

			Expect(err).To(BeNil())
		})
	})

	Describe("outputJSON is set to false", func() {
		It("returns a YAML map with a root credential object", func() {
			exportCreds, err := models.ExportCredentials(credentials, false)
			Expect(err).ToNot(HaveOccurred())

			var v map[string]interface{}
			var mapOfInterfaces []interface{}

			err = yaml.Unmarshal(exportCreds.Bytes, &v)

			Expect(err).To(BeNil())
			Expect(v["credentials"]).NotTo(BeNil())
			Expect(v["credentials"]).To(BeAssignableToTypeOf(mapOfInterfaces))
		})

		It("lists each credential", func() {
			exportCreds, _ := models.ExportCredentials(credentials, false)

			var v map[string]interface{}
			_ = yaml.Unmarshal(exportCreds.Bytes, &v)

			exportedCredentials := v["credentials"].([]interface{})

			Expect(exportedCredentials).To(HaveLen(len(credentials)))
		})

		It("includes only a name, type and value in each credential", func() {
			expectedKeys := []string{"name", "type", "value", "metadata"}
			exportCreds, _ := models.ExportCredentials(credentials, false)

			var v map[string]interface{}
			_ = yaml.Unmarshal(exportCreds.Bytes, &v)

			exportedCredentials := v["credentials"].([]interface{})

			for _, credential := range exportedCredentials {
				c := credential.(map[string]interface{})

				for k := range c {
					Expect(expectedKeys).To(ContainElement(k))
				}
			}
		})

		It("produces YAML that can be reimported", func() {
			exportCreds, _ := models.ExportCredentials(credentials, false)
			credImporter := &models.CredentialBulkImport{}

			err := credImporter.ReadBytes(exportCreds.Bytes, false)

			Expect(err).To(BeNil())
		})
	})
})

var _ = Describe("CredentialBulkExport", func() {
	Describe("String", func() {
		testString := "test"
		subject := models.CredentialBulkExport{[]byte(testString)}

		It("returns a string representation of the Bytes", func() {
			Expect(subject.String()).To(Equal(testString))
		})
	})
})
