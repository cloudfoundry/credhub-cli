package models_test

import (
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ImportFile", func() {
	Describe("readBytes()", func() {
		It("parses YAML", func() {
			var importFile models.ImportFile
			err := importFile.ReadBytes(
				[]byte(
					`credentials:
- name: /director/deployment/blobstore - agent
  type: password
  value: gx4ll8193j5rw0wljgqo
- name: /director/deployment/blobstore - director
  type: value
  value: y14ck84ef51dnchgk4kp`))

			Expect(err).To(BeNil())
			Expect(len(importFile.Credentials)).To(Equal(2))
			Expect(importFile.Credentials[0].Name).To(Equal("/director/deployment/blobstore - agent"))
			Expect(importFile.Credentials[1].Name).To(Equal("/director/deployment/blobstore - director"))
			Expect(importFile.Credentials[0].Type).To(Equal("password"))
			Expect(importFile.Credentials[1].Type).To(Equal("value"))
			Expect(importFile.Credentials[0].Value).To(Equal("gx4ll8193j5rw0wljgqo"))
			Expect(importFile.Credentials[1].Value).To(Equal("y14ck84ef51dnchgk4kp"))
		})
	})

	Describe("readFile()", func() {
		It("parses YAML from an input file", func() {
			var importFile models.ImportFile
			err := importFile.ReadFile("../test/test_import_file.yml")

			Expect(err).To(BeNil())
			Expect(len(importFile.Credentials)).To(Equal(2))
			Expect(importFile.Credentials[0].Name).To(Equal("/director/deployment/blobstore - agent"))
			Expect(importFile.Credentials[1].Name).To(Equal("/director/deployment/blobstore - director"))
			Expect(importFile.Credentials[0].Type).To(Equal("password"))
			Expect(importFile.Credentials[1].Type).To(Equal("value"))
			Expect(importFile.Credentials[0].Value).To(Equal("gx4ll8193j5rw0wljgqo"))
			Expect(importFile.Credentials[1].Value).To(Equal("y14ck84ef51dnchgk4kp"))
		})
	})
})
