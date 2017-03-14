package models

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Secret", func() {
	Describe("Terminal", func() {
		It("renders string secrets", func() {
			stringSecret := Secret{
				SecretBody: SecretBody{
					Name:             "stringSecret",
					SecretType:       "value",
					Value:            "my-value",
					VersionCreatedAt: "2016-01-01T12:00:00Z",
				},
			}

			Expect(stringSecret.Terminal()).To(Equal("" +
				"type: value\n" +
				"name: stringSecret\n" +
				"value: my-value\n" +
				"updated: 2016-01-01T12:00:00Z\n"))
		})

		It("renders ssh secrets", func() {
			ssh := RsaSsh{PublicKey: "my-pub", PrivateKey: "my-priv"}
			sshSecret := Secret{
				SecretBody: SecretBody{
					Name:             "sshSecret",
					SecretType:       "ssh",
					Value:            ssh,
					VersionCreatedAt: "2016-01-01T12:00:00Z",
				},
			}

			Expect(sshSecret.Terminal()).To(Equal("" +
				"type: ssh\n" +
				"name: sshSecret\n" +
				"value:\n" +
				"  public_key: my-pub\n" +
				"  private_key: my-priv\n" +
				"updated: 2016-01-01T12:00:00Z\n"))
		})

		It("renders rsa secrets", func() {
			rsa := RsaSsh{PublicKey: "my-pub", PrivateKey: "my-priv"}
			sshSecret := Secret{
				SecretBody: SecretBody{
					Name:             "rsaSecret",
					SecretType:       "rsa",
					Value:            rsa,
					VersionCreatedAt: "2016-01-01T12:00:00Z",
				},
			}

			Expect(sshSecret.Terminal()).To(Equal("" +
				"type: rsa\n" +
				"name: rsaSecret\n" +
				"value:\n" +
				"  public_key: my-pub\n" +
				"  private_key: my-priv\n" +
				"updated: 2016-01-01T12:00:00Z\n"))
		})

		Describe("renders certificate secrets", func() {
			It("when fields have non-nil values", func() {
				certificate := Certificate{Ca: "my-ca", Certificate: "my-cert", PrivateKey: "my-priv"}
				certificateSecret := Secret{
					SecretBody: SecretBody{
						Name:             "nonNulledSecret",
						SecretType:       "certificate",
						Value:            certificate,
						VersionCreatedAt: "2016-01-01T12:00:00Z",
					},
				}

				Expect(certificateSecret.Terminal()).To(Equal("" +
					"type: certificate\n" +
					"name: nonNulledSecret\n" +
					"value:\n" +
					"  ca: my-ca\n" +
					"  certificate: my-cert\n" +
					"  private_key: my-priv\n" +
					"updated: 2016-01-01T12:00:00Z\n"))
			})

			It("when some fields have nil values", func() {
				certificate := Certificate{Ca: "my-ca", Certificate: "", PrivateKey: "my-priv"}
				certificateSecret := Secret{
					SecretBody: SecretBody{
						Name:             "nonNulledSecret",
						SecretType:       "certificate",
						Value:            certificate,
						VersionCreatedAt: "2016-01-01T12:00:00Z",
					},
				}

				Expect(certificateSecret.Terminal()).To(Equal("" +
					"type: certificate\n" +
					"name: nonNulledSecret\n" +
					"value:\n" +
					"  ca: my-ca\n" +
					"  private_key: my-priv\n" +
					"updated: 2016-01-01T12:00:00Z\n"))
			})

			It("when fields all have nil values", func() {
				certificateSecret := Secret{
					SecretBody: SecretBody{
						Name:             "nulledSecret",
						SecretType:       "certificate",
						Value:            Certificate{},
						VersionCreatedAt: "2016-01-01T12:00:00Z",
					},
				}

				Expect(certificateSecret.Terminal()).To(Equal("" +
					"type: certificate\n" +
					"name: nulledSecret\n" +
					"value: {}\n" +
					"updated: 2016-01-01T12:00:00Z\n"))
			})
		})
	})

	Describe("JSON", func() {
		It("renders string secrets", func() {
			stringSecret := Secret{
				SecretBody: SecretBody{
					Name:             "stringSecret",
					SecretType:       "value",
					Value:            "my-value",
					VersionCreatedAt: "2016-01-01T12:00:00Z",
				},
			}

			Expect(stringSecret.Json()).To(MatchJSON(`{
				"type": "value",
				"name": "stringSecret",
				"value": "my-value",
				"version_created_at": "2016-01-01T12:00:00Z"
			}`))
		})

		It("renders ssh secrets", func() {
			ssh := RsaSsh{PublicKey: "my-pub", PrivateKey: "my-priv"}
			sshSecret := Secret{
				SecretBody: SecretBody{
					Name:             "sshSecret",
					SecretType:       "ssh",
					Value:            ssh,
					VersionCreatedAt: "2016-01-01T12:00:00Z",
				},
			}

			Expect(sshSecret.Json()).To(MatchJSON(`{
				"type": "ssh",
				"name": "sshSecret",
				"version_created_at": "2016-01-01T12:00:00Z",
				"value": {
					"public_key": "my-pub",
					"private_key": "my-priv"
				}
			}`))
		})

		It("renders rsa secrets", func() {
			rsa := RsaSsh{PublicKey: "my-pub", PrivateKey: "my-priv"}
			sshSecret := Secret{
				SecretBody: SecretBody{
					Name:             "rsaSecret",
					SecretType:       "rsa",
					Value:            rsa,
					VersionCreatedAt: "2016-01-01T12:00:00Z",
				},
			}

			Expect(sshSecret.Json()).To(MatchJSON(`{
				"type": "rsa",
				"name": "rsaSecret",
				"version_created_at": "2016-01-01T12:00:00Z",
				"value": {
					"public_key": "my-pub",
					"private_key": "my-priv"
				}
			}`))
		})

		Describe("renders certificate secrets", func() {
			It("when fields have non-nil values", func() {
				certificate := Certificate{Ca: "my-ca", Certificate: "my-cert", PrivateKey: "my-priv"}
				certificateSecret := Secret{
					SecretBody: SecretBody{
						Name:             "nonNulledSecret",
						SecretType:       "certificate",
						Value:            certificate,
						VersionCreatedAt: "2016-01-01T12:00:00Z",
					},
				}

				Expect(certificateSecret.Json()).To(MatchJSON(`{
					"type": "certificate",
					"name": "nonNulledSecret",
					"version_created_at": "2016-01-01T12:00:00Z",
					"value": {
						"ca": "my-ca",
						"certificate": "my-cert",
						"private_key": "my-priv"
					}
				}`))
			})

			It("when some fields have nil values", func() {
				certificate := Certificate{Ca: "my-ca", Certificate: "", PrivateKey: "my-priv"}
				certificateSecret := Secret{
					SecretBody: SecretBody{
						Name:             "nonNulledSecret",
						SecretType:       "certificate",
						Value:            certificate,
						VersionCreatedAt: "2016-01-01T12:00:00Z",
					},
				}

				Expect(certificateSecret.Json()).To(MatchJSON(`{
					"type": "certificate",
					"name": "nonNulledSecret",
					"version_created_at": "2016-01-01T12:00:00Z",
					"value": {
						"ca": "my-ca",
						"private_key": "my-priv"
					}
				}`))
			})

			It("when fields all have nil values", func() {
				certificateSecret := Secret{
					SecretBody: SecretBody{
						Name:             "nulledSecret",
						SecretType:       "certificate",
						Value:            Certificate{},
						VersionCreatedAt: "2016-01-01T12:00:00Z",
					},
				}

				Expect(certificateSecret.Json()).To(MatchJSON(`{
					"type": "certificate",
					"name": "nulledSecret",
					"value": {},
					"version_created_at": "2016-01-01T12:00:00Z"
				}`))
			})
		})
	})
})
