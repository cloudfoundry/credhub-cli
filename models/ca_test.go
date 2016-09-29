package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/credhub-cli/models"
)

var _ = Describe("String function", func() {
	It("when fields have non-nil values", func() {
		params := CaParameters{Certificate: "my-cert", PrivateKey: "my-priv"}
		stringCa := NewCa("stringCa", CaBody{ContentType: "root", Value: &params})
		Expect(stringCa.String()).To(Equal("" +
			"Type:          root\n" +
			"Name:          stringCa\n" +
			"Certificate:   my-cert\n" +
			"Private Key:   my-priv"))
	})
})
