package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("String function", func() {
	It("when fields have non-nil values", func() {
		params := CaParameters{Public: "my-pub", Private: "my-priv"}
		stringCa := NewCa("stringCa", CaBody{Ca: &params})
		Expect(stringCa.String()).To(Equal("" +
			"Name:		stringCa\n" +
			"Public:		my-pub\n" +
			"Private:	my-priv"))
	})
})
