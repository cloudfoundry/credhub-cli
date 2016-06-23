package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("String function", func() {

	It("renders string secrets", func() {
		stringSecret := NewSecret("stringSecret", SecretBody{ContentType: "value", Value: "my-value", UpdatedAt: "2016-01-01T12:00:00Z"})
		Expect(stringSecret.String()).To(Equal("" +
			"Type:		value\n" +
			"Name:		stringSecret\n" +
			"Value:		my-value\n" +
			"Updated:	2016-01-01T12:00:00Z"))
	})

	Describe("renders certificate secrets", func() {

		It("when fields have non-nil values", func() {
			cert := Certificate{Ca: "my-ca", Public: "my-pub", Private: "my-priv"}
			certificateSecret := NewSecret("nonNulledSecret", SecretBody{ContentType: "certificate", Certificate: &cert, UpdatedAt: "2016-01-01T12:00:00Z"})
			Expect(certificateSecret.String()).To(Equal("" +
				"Type:		certificate\n" +
				"Name:		nonNulledSecret\n" +
				"CA:		my-ca\n" +
				"Public:		my-pub\n" +
				"Private:	my-priv\n" +
				"Updated:	2016-01-01T12:00:00Z"))
		})

		It("when some fields have nil values", func() {
			cert := Certificate{Ca: "my-ca", Public: "", Private: "my-priv"}
			certificateSecret := NewSecret("nonNulledSecret", SecretBody{ContentType: "certificate", Certificate: &cert, UpdatedAt: "2016-01-01T12:00:00Z"})
			Expect(certificateSecret.String()).To(Equal("" +
				"Type:		certificate\n" +
				"Name:		nonNulledSecret\n" +
				"CA:		my-ca\n" +
				"Private:	my-priv\n" +
				"Updated:	2016-01-01T12:00:00Z"))
		})

		It("when fields all have nil values", func() {
			cert := Certificate{}
			certificateSecret := NewSecret("nulledSecret", SecretBody{ContentType: "certificate", Certificate: &cert, UpdatedAt: "2016-01-01T12:00:00Z"})
			Expect(certificateSecret.String()).To(Equal("" +
				"Type:		certificate\n" +
				"Name:		nulledSecret\n" +
				"Updated:	2016-01-01T12:00:00Z"))
		})
	})
})
