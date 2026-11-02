package commands_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Import", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("import", "-f", "../test/test_import_file.yml")
	ItRequiresAnAPIToBeSet("import", "-f", "../test/test_import_file.yml")

	Describe("importing a file with mixed credentials", func() {
		Describe("when importing yaml", func() {
			It("sets all the credentials", func() {
				setUpImportRequests()

				session := runCommand("import", "-f", "../test/test_import_file.yml")

				Eventually(session).Should(Exit(0))

				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 7
Failed to set: 0
`))
			})
		})
		Describe("when importing json", func() {
			It("sets all the credentials", func() {
				setUpImportRequests()

				session := runCommand("import", "-f", "../test/test_import_file.json", "-j")

				Eventually(session).Should(Exit(0))

				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 7
Failed to set: 0
`))
			})
		})
	})

	Describe("when the yaml file starts with ---", func() {
		It("sets all the credentials", func() {
			setUpImportRequests()

			session := runCommand("import", "-f", "../test/test_import_file_with_document_end.yml")

			Eventually(session).Should(Exit(0))

			Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 7
Failed to set: 0
`))
		})
	})

	Describe("when the yaml file starts with --- and spaces for Iryna", func() {
		It("sets all the credentials", func() {
			setUpImportRequests()

			session := runCommand("import", "-f", "../test/test_import_file_with_document_end_and_spaces.yml")

			Eventually(session).Should(Exit(0))

			Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 7
Failed to set: 0
`))
		})
	})

	Describe("when importing file with no name specified", func() {
		Describe("when importing yaml", func() {
			It("passes through the server error", func() {
				jsonBody := `{"name":"","type":"password","value":"test-password"}`
				setupPutBadRequestServer(jsonBody)

				session := runCommand("import", "-f", "../test/test_import_missing_name.yml")

				Eventually(session.Out).Should(Say(`test error`))
			})
		})
		Describe("when importing json", func() {
			It("passes through the server error", func() {
				jsonBody := `{"name":"","type":"password","value":"test-password"}`
				setupPutBadRequestServer(jsonBody)

				session := runCommand("import", "-f", "../test/test_import_missing_name.json", "-j")

				Eventually(session.Out).Should(Say(`test error`))
			})
		})
	})

	Describe("when importing file with incorrect structure", func() {
		Describe("when importing yaml", func() {
			It("returns an error message", func() {
				errorMessage := "The referenced file does not contain valid yaml structure. Please update and retry your request."

				session := runCommand("import", "-f", "../test/test_import_incorrect_yaml.yml")

				Eventually(session.Err).Should(Say(errorMessage))
			})
		})
		Describe("when importing json", func() {
			It("returns an error message", func() {
				errorMessage := `The referenced file does not contain valid json structure. Please update and retry your request.`

				session := runCommand("import", "-f", "../test/test_import_incorrect_json.json", "-j")

				Eventually(session.Err).Should(Say(errorMessage))
			})
		})
	})

	Describe("when some credentials fail to set it prints errors in summary", func() {
		Describe("when importing yaml", func() {
			It("should display error message", func() {
				request := `{"type":"invalid_type","name":"/test/invalid_type","value":"some string"}`
				request1 := `{"type":"invalid_type","name":"/test/invalid_type1","value":"some string"}`
				setupPutBadRequestServer(request)
				setupPutBadRequestServer(request1)
				setupSetServer("/test/user", "user", `{"username": "covfefe", "password": "test-user-password"}`)

				session := runCommand("import", "-f", "../test/test_import_partial_fail_set.yml")
				summaryMessage := `Import complete.
Successfully set: 1
Failed to set: 2
`
				Eventually(session.Out).Should(Say(`Credential '/test/invalid_type' at index 0 could not be set: test error`))
				Eventually(session.Out).Should(Say(`Credential '/test/invalid_type1' at index 1 could not be set: test error`))
				Expect(session.ExitCode()).ToNot(Equal(0))
				Expect(session.Out.Contents()).NotTo(ContainSubstring(`id: 5a2edd4f-1686-4c8d-80eb-5daa866f9f86`))
				Expect(session.Err.Contents()).To(ContainSubstring(`One or more credentials failed to import`))
				Eventually(session.Out).Should(Say(summaryMessage))
			})
		})
		Describe("when importing json", func() {
			It("should display error message", func() {
				request := `{"type":"invalid_type","name":"/test/invalid_type","value":"some string"}`
				request1 := `{"type":"invalid_type","name":"/test/invalid_type1","value":"some string"}`
				setupPutBadRequestServer(request)
				setupPutBadRequestServer(request1)
				setupSetServer("/test/user", "user", `{"username": "covfefe", "password": "test-user-password"}`)

				session := runCommand("import", "-f", "../test/test_import_partial_fail_set.json", "-j")
				summaryMessage := `Import complete.
Successfully set: 1
Failed to set: 2
`
				Eventually(session.Out).Should(Say(`Credential '/test/invalid_type' at index 0 could not be set: test error`))
				Eventually(session.Out).Should(Say(`Credential '/test/invalid_type1' at index 1 could not be set: test error`))
				Expect(session.ExitCode()).ToNot(Equal(0))
				Expect(session.Out.Contents()).NotTo(ContainSubstring(`id: 5a2edd4f-1686-4c8d-80eb-5daa866f9f86`))
				Expect(session.Err.Contents()).To(ContainSubstring(`One or more credentials failed to import`))
				Eventually(session.Out).Should(Say(summaryMessage))
			})
		})
	})

	Describe("when no credential tag present in import file", func() {
		Describe("when importing yaml", func() {
			It("prints correct error message", func() {
				session := runCommand("import", "-f", "../test/test_import_incorrect_format.yml")

				noCredentialTagError := "The referenced file does not contain valid yaml structure. Please update and retry your request."
				Eventually(session.Err).Should(Say(noCredentialTagError))
			})
		})
		Describe("when importing json", func() {
			It("prints correct error message", func() {
				session := runCommand("import", "-f", "../test/test_import_incorrect_format.json", "-j")

				noCredentialTagError := `The referenced file does not contain valid json structure. Please update and retry your request.`
				Eventually(session.Err).Should(Say(noCredentialTagError))
			})
		})
	})

	Describe("when importing an ssh type with key public_key_fingerprint", func() {
		Describe("when importing yaml", func() {
			It("ignore public_key_fingerprint", func() {
				setupSetServer("/test/sshCred", "ssh", `{"public_key":"some-key","private_key":"some-private-key"}`)

				session := runCommand("import", "-f", "../test/test_import_ssh_type_with_public_key_fingerprint.yml")
				Eventually(session).Should(Exit(0))
				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})
		Describe("when importing json", func() {
			It("ignore public_key_fingerprint", func() {
				setupSetServer("/test/sshCred", "ssh", `{"public_key":"some-key","private_key":"some-private-key"}`)

				session := runCommand("import", "-f", "../test/test_import_ssh_type_with_public_key_fingerprint.json", "-j")
				Eventually(session).Should(Exit(0))
				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})

	})

	Describe("when importing a user type with password_hash", func() {
		Describe("when importing yaml", func() {
			It("ignore password_hash", func() {
				setupSetServer("/test/userCred", "user", `{"username": "sample-username", "password": "test-user-password"}`)

				session := runCommand("import", "-f", "../test/test_import_user_type_with_password_hash.yml")
				Eventually(session).Should(Exit(0))
				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})
		Describe("when importing json", func() {
			It("ignore password_hash", func() {
				setupSetServer("/test/userCred", "user", `{"username": "sample-username", "password": "test-user-password"}`)

				session := runCommand("import", "-f", "../test/test_import_user_type_with_password_hash.json", "-j")
				Eventually(session).Should(Exit(0))
				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})

	})

	Describe("when importing a value as a integer", func() {
		Describe("when importing yaml", func() {
			It("casts int to string and successfully imports", func() {
				setupSetServer("/test/intStringValue", "value", `"123"`)

				session := runCommand("import", "-f", "../test/test_import_with_int_for_value.yml")

				Eventually(session).Should(Exit(0))

				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})
		Describe("when importing json", func() {
			It("casts int to string and successfully imports", func() {
				setupSetServer("/test/intStringValue", "value", `"123"`)

				session := runCommand("import", "-f", "../test/test_import_with_int_for_value.json", "-j")

				Eventually(session).Should(Exit(0))

				Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
			})
		})

	})

	Describe("when importing certificate chain", func() {
		Context("and leaf comes after signing CA", func() {
			Describe("when importing yaml", func() {
				It("imports the signing CA first", func() {
					setupSetServer("/root_ca", "certificate", `{"ca":"root-ca","certificate":"root-certificate","private_key":"root-private-key"}`)
					setupSetServer("/intermediate_ca", "certificate", `{"ca_name":"/root_ca","certificate":"intermediate-certificate","private_key":"intermediate-private-key"}`)
					setupSetServer("/leaf_cert", "certificate", `{"ca_name":"/intermediate_ca","certificate":"leaf-certificate","private_key":"leaf-private-key"}`)

					session := runCommand("import", "-f", "../test/certificate-chain.yml")
					Eventually(session).Should(Exit(0))
					Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 3
Failed to set: 0
`))

				})
			})
			Describe("when importing json", func() {
				It("imports the signing CA first", func() {
					setupSetServer("/root_ca", "certificate", `{"ca":"root-ca","certificate":"root-certificate","private_key":"root-private-key"}`)
					setupSetServer("/intermediate_ca", "certificate", `{"ca_name":"/root_ca","certificate":"intermediate-certificate","private_key":"intermediate-private-key"}`)
					setupSetServer("/leaf_cert", "certificate", `{"ca_name":"/intermediate_ca","certificate":"leaf-certificate","private_key":"leaf-private-key"}`)

					session := runCommand("import", "-f", "../test/certificate-chain.json", "-j")
					Eventually(session).Should(Exit(0))
					Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 3
Failed to set: 0
`))
				})
			})

		})
	})

	Describe("when the ca_name graph contains a cycle", func() {
		for _, format := range []struct {
			name string
			ext  string
			args []string
		}{
			{name: "yaml", ext: ".yml"},
			{name: "json", ext: ".json", args: []string{"-j"}},
		} {
			format := format
			importFixture := func(base string) *Session {
				return runCommand(append([]string{"import", "-f", "../test/" + base + format.ext}, format.args...)...)
			}

			Describe("when importing "+format.name, func() {
				It("reports a two-certificate cycle as failed without setting either", func() {
					session := importFixture("certificate-cycle-two")

					Eventually(session).Should(Exit(1))
					Expect(receivedSetRequests()).To(Equal(0))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Import complete.
Successfully set: 0
Failed to set: 2
`))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Credential '/cert_a' at index 0 could not be set: circular ca_name reference: /cert_a -> /cert_b -> /cert_a`))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Credential '/cert_b' at index 1 could not be set: circular ca_name reference: /cert_a -> /cert_b -> /cert_a`))
					Expect(string(session.Err.Contents())).To(ContainSubstring(`One or more credentials failed to import`))
				})

				It("reports a three-certificate cycle as failed without setting any", func() {
					session := importFixture("certificate-cycle-three")

					Eventually(session).Should(Exit(1))
					Expect(receivedSetRequests()).To(Equal(0))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Import complete.
Successfully set: 0
Failed to set: 3
`))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Credential '/cert_a' at index 0 could not be set: circular ca_name reference: /cert_a -> /cert_b -> /cert_c -> /cert_a`))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Credential '/cert_b' at index 1 could not be set: circular ca_name reference: /cert_a -> /cert_b -> /cert_c -> /cert_a`))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Credential '/cert_c' at index 2 could not be set: circular ca_name reference: /cert_a -> /cert_b -> /cert_c -> /cert_a`))
				})

				It("imports the valid chain and reports the cycle in the same file", func() {
					setupSetServer("/root_ca", "certificate", `{"ca":"root-ca","certificate":"root-ca-certificate","private_key":"root-ca-private-key"}`)
					setupSetServer("/valid_leaf", "certificate", `{"ca_name":"/root_ca","certificate":"valid-leaf-certificate","private_key":"valid-leaf-private-key"}`)

					session := importFixture("certificate-cycle-mixed")

					Eventually(session).Should(Exit(1))
					Expect(receivedSetRequests()).To(Equal(2))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Import complete.
Successfully set: 2
Failed to set: 2
`))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Credential '/cert_a' at index 1 could not be set: circular ca_name reference: /cert_a -> /cert_b -> /cert_a`))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Credential '/cert_b' at index 2 could not be set: circular ca_name reference: /cert_a -> /cert_b -> /cert_a`))
				})

				It("sends a leaf signed by a cycle member to the server and reports its rejection", func() {
					setupPutBadRequestServer(`{"type":"certificate","name":"/leaf_cert","value":{"ca_name":"/cert_a","certificate":"leaf-cert-certificate","private_key":"leaf-cert-private-key"}}`)

					session := importFixture("certificate-cycle-with-leaf")

					Eventually(session).Should(Exit(1))
					Expect(receivedSetRequests()).To(Equal(1))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Import complete.
Successfully set: 0
Failed to set: 3
`))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Credential '/leaf_cert' at index 0 could not be set: test error`))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Credential '/cert_a' at index 1 could not be set: circular ca_name reference: /cert_a -> /cert_b -> /cert_a`))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Credential '/cert_b' at index 2 could not be set: circular ca_name reference: /cert_a -> /cert_b -> /cert_a`))
				})
			})

			Describe("when importing "+format.name+" a certificate whose ca_name is its own name", func() {
				It("leaves the decision to the server and imports it when the server accepts", func() {
					setupSetServer("/self_ca", "certificate", `{"ca_name":"/self_ca","certificate":"self-ca-certificate","private_key":"self-ca-private-key"}`)

					session := importFixture("certificate-self-reference")

					Eventually(session).Should(Exit(0))
					Expect(receivedSetRequests()).To(Equal(1))
					Expect(string(session.Out.Contents())).To(Equal(`Import complete.
Successfully set: 1
Failed to set: 0
`))
				})

				It("reports the server's rejection as an ordinary failure, not a cycle", func() {
					setupPutBadRequestServer(`{"type":"certificate","name":"/self_ca","value":{"ca_name":"/self_ca","certificate":"self-ca-certificate","private_key":"self-ca-private-key"}}`)

					session := importFixture("certificate-self-reference")

					Eventually(session).Should(Exit(1))
					Expect(receivedSetRequests()).To(Equal(1))
					Expect(string(session.Out.Contents())).To(ContainSubstring(`Credential '/self_ca' at index 0 could not be set: test error`))
					Expect(string(session.Out.Contents())).NotTo(ContainSubstring(`circular`))
				})
			})
		}
	})
})

func receivedSetRequests() int {
	count := 0
	for _, request := range server.ReceivedRequests() {
		if request.Method == "PUT" && request.URL.Path == "/api/v1/data" {
			count++
		}
	}
	return count
}

func setUpImportRequests() {
	setupSetServer("/test/password", "password", `"test-password-value"`)
	setupSetServerWithMetadata("/test/value", "value", `"test-value"`, `{"some":"thing", "nested":{"with":"value"}}`)
	setupSetServerWithMetadata("/test/certificate", "certificate", `{"ca":"ca-certificate","certificate":"certificate","private_key":"private-key"}`, `{"certs":"metadata"}`)
	setupSetServer("/test/rsa", "rsa", `{"public_key":"public-key","private_key":"private-key"}`)
	setupSetServer("/test/ssh", "ssh", `{"public_key":"ssh-public-key","private_key":"private-key"}`)
	setupSetServer("/test/user", "user", `{"username": "covfefe", "password": "test-user-password"}`)
	setupSetServer("/test/json", "json", `{"1":"key is not a string","3.14":"pi","true":"key is a bool","arbitrary_object":{"nested_array":["array_val1",{"array_object_subvalue":"covfefe"}]}}`)
}
