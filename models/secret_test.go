package models

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("String function", func() {

	It("renders string secrets", func() {
		stringSecret := NewSecret("stringSecret", SecretBody{ContentType: "value", Value: "my-value", UpdatedAt: "2016-01-01T12:00:00Z"})
		Expect(stringSecret.String()).To(Equal("" +
			"Type:          value\n" +
			"Name:          stringSecret\n" +
			"Value:         my-value\n" +
			"Updated:       2016-01-01T12:00:00Z"))
	})

	It("renders ssh secrets", func() {
		jsonBytes, _ := json.Marshal(Rsa{PublicKey: "my-pub", PrivateKey: "my-priv"})
		sshSecret := NewSecret("sshSecret", SecretBody{ContentType: "ssh", Value: unmarshal(jsonBytes), UpdatedAt: "2016-01-01T12:00:00Z"})
		Expect(sshSecret.String()).To(Equal("" +
			"Type:          ssh\n" +
			"Name:          sshSecret\n" +
			"Public Key:    my-pub\n" +
			"Private Key:   my-priv\n" +
			"Updated:       2016-01-01T12:00:00Z"))
	})

	It("renders rsa secrets", func() {
		jsonBytes, _ := json.Marshal(Rsa{PublicKey: "my-pub", PrivateKey: "my-priv"})
		sshSecret := NewSecret("rsaSecret", SecretBody{ContentType: "rsa", Value: unmarshal(jsonBytes), UpdatedAt: "2016-01-01T12:00:00Z"})
		Expect(sshSecret.String()).To(Equal("" +
			"Type:          rsa\n" +
			"Name:          rsaSecret\n" +
			"Public Key:    my-pub\n" +
			"Private Key:   my-priv\n" +
			"Updated:       2016-01-01T12:00:00Z"))
	})

	Describe("renders certificate secrets", func() {

		It("when fields have non-nil values", func() {
			jsonBytes, _ := json.Marshal(Certificate{Ca: "my-ca", Certificate: "my-cert", PrivateKey: "my-priv"})
			certificateSecret := NewSecret("nonNulledSecret", SecretBody{ContentType: "certificate", Value: unmarshal(jsonBytes), UpdatedAt: "2016-01-01T12:00:00Z"})
			Expect(certificateSecret.String()).To(Equal("" +
				"Type:          certificate\n" +
				"Name:          nonNulledSecret\n" +
				"Ca:            my-ca\n" +
				"Certificate:   my-cert\n" +
				"Private Key:   my-priv\n" +
				"Updated:       2016-01-01T12:00:00Z"))
		})

		It("when some fields have nil values", func() {
			jsonBytes, _ := json.Marshal(Certificate{Ca: "my-ca", Certificate: "", PrivateKey: "my-priv"})
			certificateSecret := NewSecret("nonNulledSecret", SecretBody{ContentType: "certificate", Value: unmarshal(jsonBytes), UpdatedAt: "2016-01-01T12:00:00Z"})
			Expect(certificateSecret.String()).To(Equal("" +
				"Type:          certificate\n" +
				"Name:          nonNulledSecret\n" +
				"Ca:            my-ca\n" +
				"Private Key:   my-priv\n" +
				"Updated:       2016-01-01T12:00:00Z"))
		})

		It("when fields all have nil values", func() {
			certificateSecret := NewSecret("nulledSecret", SecretBody{ContentType: "certificate", Value: unmarshal([]byte{}), UpdatedAt: "2016-01-01T12:00:00Z"})
			Expect(certificateSecret.String()).To(Equal("" +
				"Type:          certificate\n" +
				"Name:          nulledSecret\n" +
				"Updated:       2016-01-01T12:00:00Z"))
		})
	})
})

func unmarshal(jsonBytes []byte) map[string]interface{} {
	itemMap := map[string]interface{}{}
	json.Unmarshal(jsonBytes, &itemMap)
	return itemMap
}
