package commands_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Import", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("import", "-f", "../test/test_import_file.yml")
	ItRequiresAnAPIToBeSet("import", "-f", "../test/test_import_file.yml")

	Describe("importing a file with mixed credentials", func() {
		It("sets all the credentials", func() {
			setUpImportRequests()

			session := runCommand("import", "-f", "../test/test_import_file.yml")

			Eventually(session).Should(Exit(0))

			Eventually(session.Out).Should(Say(`name: /test/password
type: password
value: test-password-value`))
			Eventually(session.Out).Should(Say(`name: /test/value
type: value
value: test-value`))
			Eventually(session.Out).Should(Say(`name: /test/certificate
type: certificate
value:
  ca: ca-certificate
  certificate: certificate
  private_key: private-key`))
			Eventually(session.Out).Should(Say(`name: /test/rsa
type: rsa
value:
  private_key: private-key
  public_key: public-key`))
			Eventually(session.Out).Should(Say(`name: /test/ssh
type: ssh
value:
  private_key: private-key
  public_key: ssh-public-key`))
			Eventually(session.Out).Should(Say(`name: /test/user
type: user
value:
  password: test-user-password
  password_hash: P455W0rd-H45H
  username: covfefe`))
			Eventually(session.Out).Should(Say(`name: /test/json
type: json
value:
  "1": key is not a string
  "3.14": pi
  arbitrary_object:
    nested_array:
    - array_val1
    - array_object_subvalue: covfefe
  "true": key is a bool
`))
			Eventually(session.Out).Should(Say(`Import complete.
Successfully set: 7
Failed to set: 0
`))
		})
	})

	Describe("when the yaml file starts with ---", func() {
		It("sets all the credentials", func() {
			setUpImportRequests()

			session := runCommand("import", "-f", "../test/test_import_file_with_document_end.yml")

			Eventually(session).Should(Exit(0))

			Eventually(session.Out).Should(Say(`name: /test/password
type: password
value: test-password-value`))
			Eventually(session.Out).Should(Say(`name: /test/value
type: value
value: test-value`))
			Eventually(session.Out).Should(Say(`name: /test/certificate
type: certificate
value:
  ca: ca-certificate
  certificate: certificate
  private_key: private-key`))
			Eventually(session.Out).Should(Say(`name: /test/rsa
type: rsa
value:
  private_key: private-key
  public_key: public-key`))
			Eventually(session.Out).Should(Say(`name: /test/ssh
type: ssh
value:
  private_key: private-key
  public_key: ssh-public-key`))
			Eventually(session.Out).Should(Say(`name: /test/user
type: user
value:
  password: test-user-password
  password_hash: P455W0rd-H45H
  username: covfefe`))
			Eventually(session.Out).Should(Say(`name: /test/json
type: json
value:
  "1": key is not a string
  "3.14": pi
  arbitrary_object:
    nested_array:
    - array_val1
    - array_object_subvalue: covfefe
  "true": key is a bool
`))
			Eventually(session.Out).Should(Say(`Import complete.
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

			Eventually(session.Out).Should(Say(`name: /test/password
type: password
value: test-password-value`))
			Eventually(session.Out).Should(Say(`name: /test/value
type: value
value: test-value`))
			Eventually(session.Out).Should(Say(`name: /test/certificate
type: certificate
value:
  ca: ca-certificate
  certificate: certificate
  private_key: private-key`))
			Eventually(session.Out).Should(Say(`name: /test/rsa
type: rsa
value:
  private_key: private-key
  public_key: public-key`))
			Eventually(session.Out).Should(Say(`name: /test/ssh
type: ssh
value:
  private_key: private-key
  public_key: ssh-public-key`))
			Eventually(session.Out).Should(Say(`name: /test/user
type: user
value:
  password: test-user-password
  password_hash: P455W0rd-H45H
  username: covfefe`))
			Eventually(session.Out).Should(Say(`name: /test/json
type: json
value:
  "1": key is not a string
  "3.14": pi
  arbitrary_object:
    nested_array:
    - array_val1
    - array_object_subvalue: covfefe
  "true": key is a bool
`))
			Eventually(session.Out).Should(Say(`Import complete.
Successfully set: 7
Failed to set: 0
`))
		})
	})

	Describe("when importing file with no name specified", func() {
		It("passes through the server error", func() {
			jsonBody := `{"name":"","type":"password","value":"test-password","overwrite":true}`
			SetupPutBadRequestServer(jsonBody)

			session := runCommand("import", "-f", "../test/test_import_missing_name.yml")

			Eventually(session.Out).Should(Say(`test error`))
		})
	})

	Describe("when importing file with incorrect YAML", func() {
		It("returns an error message", func() {
			errorMessage := `The referenced file does not contain valid yaml structure. Please update and retry your request.`

			session := runCommand("import", "-f", "../test/test_import_incorrect_yaml.yml")

			Eventually(session.Err).Should(Say(errorMessage))
		})
	})

	Describe("when some credentials fail to set it prints errors in summary", func() {
		It("should display error message", func() {
			error := "The request does not include a valid type. Valid values include 'value', 'json', 'password', 'user', 'certificate', 'ssh' and 'rsa'."

			request := `{"type":"invalid_type","name":"/test/invalid_type","value":"some string","overwrite":true}`
			request1 := `{"type":"invalid_type","name":"/test/invalid_type1","value":"some string","overwrite":true}`
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/api/v1/data"),
					VerifyJSON(request),
					RespondWith(http.StatusBadRequest, `{"error": "`+error+`"}`)),
				CombineHandlers(
					VerifyRequest("PUT", "/api/v1/data"),
					VerifyJSON(request1),
					RespondWith(http.StatusBadRequest, `{"error": "`+error+`"}`)),
			)
			SetupPutUserServer("/test/user", `{"username": "covfefe", "password": "test-user-password"}`, "covfefe", "test-user-password", "P455W0rd-H45H", true)

			session := runCommand("import", "-f", "../test/test_import_partial_fail_set.yml")
			successfulSetMessage := `Credential '/test/invalid_type' at index 0 could not be set: The request does not include a valid type. Valid values include 'value', 'json', 'password', 'user', 'certificate', 'ssh' and 'rsa'.

Credential '/test/invalid_type1' at index 1 could not be set: The request does not include a valid type. Valid values include 'value', 'json', 'password', 'user', 'certificate', 'ssh' and 'rsa'.

id: 5a2edd4f-1686-4c8d-80eb-5daa866f9f86
version_created_at: 2016-01-01T12:00:00Z
name: /test/user
type: user
value:
  password: test-user-password
  password_hash: P455W0rd-H45H
  username: covfefe`
			summaryMessage := `Import complete.
Successfully set: 1
Failed to set: 2
 - Credential '/test/invalid_type' at index 0 could not be set: The request does not include a valid type. Valid values include 'value', 'json', 'password', 'user', 'certificate', 'ssh' and 'rsa'.
 - Credential '/test/invalid_type1' at index 1 could not be set: The request does not include a valid type. Valid values include 'value', 'json', 'password', 'user', 'certificate', 'ssh' and 'rsa'.
`
			Eventually(session.Out).Should(Say(successfulSetMessage))
			Eventually(session.Out).Should(Say(summaryMessage))
		})
	})

	Describe("when no credential tag present in import file", func() {
		It("prints correct error message", func() {
			session := runCommand("import", "-f", "../test/test_import_incorrect_format.yml")

			noCredentialTagError := "The referenced import file does not begin with the key 'credentials'. The import file must contain a list of credentials under the key 'credentials'. Please update and retry your request."
			Eventually(session.Err).Should(Say(noCredentialTagError))
		})
	})
})

func setUpImportRequests() {
	SetupOverwritePutValueServer("/test/password", "password", "test-password-value", true)
	SetupOverwritePutValueServer("/test/value", "value", "test-value", true)
	SetupPutCertificateServer("/test/certificate",
		`ca-certificate`,
		`certificate`,
		`private-key`)
	SetupPutRsaSshServer("/test/rsa", "rsa", "public-key", "private-key", true)
	SetupPutRsaSshServer("/test/ssh", "ssh", "ssh-public-key", "private-key", true)
	SetupPutUserServer("/test/user", `{"username": "covfefe", "password": "test-user-password"}`, "covfefe", "test-user-password", "P455W0rd-H45H", true)
	setupPutJsonServer("/test/json", `{"1":"key is not a string","3.14":"pi","true":"key is a bool","arbitrary_object":{"nested_array":["array_val1",{"array_object_subvalue":"covfefe"}]}}`)
}
