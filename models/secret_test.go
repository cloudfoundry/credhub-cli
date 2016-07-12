package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("String function", func() {

	It("renders string secrets", func() {
		stringSecret := NewSecret("stringSecret", SecretBody{ContentType: "value", Credential: "my-value", UpdatedAt: "2016-01-01T12:00:00Z"})
		Expect(stringSecret.String()).To(Equal("" +
			"Type:		value\n" +
			"Name:		stringSecret\n" +
			"Credential:	my-value\n" +
			"Updated:	2016-01-01T12:00:00Z"))
	})

	Describe("renders certificate secrets", func() {

		It("when fields have non-nil values", func() {
			cert := Certificate{Root: "my-ca", Certificate: "my-cert", Private: "my-priv"}
			certificateSecret := NewSecret("nonNulledSecret", SecretBody{ContentType: "certificate", Credential: &cert, UpdatedAt: "2016-01-01T12:00:00Z"})
			Expect(certificateSecret.String()).To(Equal("" +
				"Type:		certificate\n" +
				"Name:		nonNulledSecret\n" +
				"Root:		my-ca\n" +
				"Certificate:		my-cert\n" +
				"Private:	my-priv\n" +
				"Updated:	2016-01-01T12:00:00Z"))
		})

		It("when some fields have nil values", func() {
			cert := Certificate{Root: "my-ca", Certificate: "", Private: "my-priv"}
			certificateSecret := NewSecret("nonNulledSecret", SecretBody{ContentType: "certificate", Credential: &cert, UpdatedAt: "2016-01-01T12:00:00Z"})
			Expect(certificateSecret.String()).To(Equal("" +
				"Type:		certificate\n" +
				"Name:		nonNulledSecret\n" +
				"Root:		my-ca\n" +
				"Private:	my-priv\n" +
				"Updated:	2016-01-01T12:00:00Z"))
		})

		It("when fields all have nil values", func() {
			cert := Certificate{}
			certificateSecret := NewSecret("nulledSecret", SecretBody{ContentType: "certificate", Credential: &cert, UpdatedAt: "2016-01-01T12:00:00Z"})
			Expect(certificateSecret.String()).To(Equal("" +
				"Type:		certificate\n" +
				"Name:		nulledSecret\n" +
				"Updated:	2016-01-01T12:00:00Z"))
		})
	})
})
