package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("String function", func() {
	It("when fields have non-nil values", func() {
		params := CaParameters{Certificate: "my-cert", Private: "my-priv"}
		stringCa := NewCa("stringCa", CaBody{ContentType: "root", Ca: &params})
		Expect(stringCa.String()).To(Equal("" +
			"Type:		root\n" +
			"Name:		stringCa\n" +
			"Certificate:		my-cert\n" +
			"Private:	my-priv"))
	})
})
