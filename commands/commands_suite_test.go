package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"

	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

const TIMESTAMP = `2016-01-01T12:00:00Z`
const UUID = `5a2edd4f-1686-4c8d-80eb-5daa866f9f86`

const STRING_SECRET_OVERWRITE_REQUEST_JSON = `{"type":"%s","name":"%s","value":"%s","overwrite":%t}`
const CERTIFICATE_SECRET_REQUEST_JSON = `{"type":"certificate","name":"%s","value":{"ca":"%s","certificate":"%s","private_key":"%s"},"overwrite":%t}`
const GENERATE_SECRET_REQUEST_JSON = `{"name":"%s","type":"%s","overwrite":%t,"parameters":%s}`
const RSA_SSH_SECRET_REQUEST_JSON = `{"type":"%s","name":"%s","value":{"public_key":"%s","private_key":"%s"},"overwrite":%t}`
const GENERATE_DEFAULT_TYPE_REQUEST_JSON = `{"name":"%s","type":"password","overwrite":%t,"parameters":%s}`

const STRING_SECRET_RESPONSE_JSON = `{"type":"%s","id":"` + UUID + `","name":"%s","version_created_at":"` + TIMESTAMP + `","value":"%s"}`
const CERTIFICATE_SECRET_RESPONSE_JSON = `{"type":"certificate","id":"` + UUID + `","name":"%s","version_created_at":"` + TIMESTAMP + `","value":{"ca":"%s","certificate":"%s","private_key":"%s"}}`
const RSA_SSH_SECRET_RESPONSE_JSON = `{"type":"%s","id":"` + UUID + `","name":"%s","version_created_at":"` + TIMESTAMP + `","value":{"public_key":"%s","private_key":"%s"},"version_created_at":"` + TIMESTAMP + `"}`

const STRING_SECRET_ARRAY_RESPONSE_JSON = `{"data":[` + STRING_SECRET_RESPONSE_JSON + `]}`
const CERTIFICATE_SECRET_ARRAY_RESPONSE_JSON = `{"data":[` + CERTIFICATE_SECRET_RESPONSE_JSON + `]}`
const RSA_SSH_SECRET_ARRAY_RESPONSE_JSON = `{"data":[` + RSA_SSH_SECRET_RESPONSE_JSON + `]}`

const STRING_SECRET_RESPONSE_TABLE = "type: %s\nname: %s\nvalue: %s\nupdated: " + TIMESTAMP
const CERTIFICATE_SECRET_RESPONSE_TABLE = "type: certificate\nname: %s\nvalue:\n  ca: %s\n  certificate: %s\n  private_key: %s\nupdated: " + TIMESTAMP
const SSH_SECRET_RESPONSE_TABLE = "type: ssh\nname: %s\nvalue:\n  public_key: %s\n  private_key: %s\nupdated: " + TIMESTAMP
const RSA_SECRET_RESPONSE_TABLE = "type: rsa\nname: %s\nvalue:\n  public_key: %s\n  private_key: %s\nupdated: " + TIMESTAMP

var responseMyValuePotatoes = fmt.Sprintf(STRING_SECRET_RESPONSE_TABLE, "value", "my-value", "potatoes")
var responseMyPasswordPotatoes = fmt.Sprintf(STRING_SECRET_RESPONSE_TABLE, "password", "my-password", "potatoes")
var responseMyCertificate = fmt.Sprintf(CERTIFICATE_SECRET_RESPONSE_TABLE, "my-secret", "my-ca", "my-cert", "my-priv")
var responseMySSHFoo = fmt.Sprintf(SSH_SECRET_RESPONSE_TABLE, "foo-ssh-key", "some-public-key", "some-private-key")
var responseMyRSAFoo = fmt.Sprintf(RSA_SECRET_RESPONSE_TABLE, "foo-rsa-key", "some-public-key", "some-private-key")

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var (
	commandPath string
	homeDir     string
	server      *Server
	authServer  *Server
)

var _ = BeforeEach(func() {
	var err error
	homeDir, err = ioutil.TempDir("", "cm-test")
	Expect(err).NotTo(HaveOccurred())

	if runtime.GOOS == "windows" {
		os.Setenv("USERPROFILE", homeDir)
	} else {
		os.Setenv("HOME", homeDir)
	}

	server = NewServer()

	authServer = NewServer()

	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("GET", "/info"),
			RespondWith(http.StatusOK, `{
					"app":{"version":"my-version","name":"CredHub"},
					"auth-server":{"url":"`+authServer.URL()+`"}
					}`),
		),
	)

	session := runCommand("api", server.URL())
	Eventually(session).Should(Exit(0))
})

var _ = AfterEach(func() {
	server.Close()
	authServer.Close()
	os.RemoveAll(homeDir)
})

var _ = SynchronizedBeforeSuite(func() []byte {
	executable_path, err := Build("github.com/cloudfoundry-incubator/credhub-cli", "-ldflags", "-X github.com/cloudfoundry-incubator/credhub-cli/version.Version=test-version")
	Expect(err).NotTo(HaveOccurred())
	return []byte(executable_path)
}, func(data []byte) {
	commandPath = string(data)
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	CleanupBuildArtifacts()
})

func runCommand(args ...string) *Session {
	cmd := exec.Command(commandPath, args...)
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

func runCommandWithStdin(stdin io.Reader, args ...string) *Session {
	cmd := exec.Command(commandPath, args...)
	cmd.Stdin = stdin
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	<-session.Exited

	return session
}

func createTempDir(prefix string) string {
	name, err := ioutil.TempDir("", prefix)
	if err != nil {
		panic(err)
	}
	return name
}

func createSecretFile(dir, filename string, contents string) string {
	path := dir + "/" + filename
	err := ioutil.WriteFile(path, []byte(contents), 0644)
	if err != nil {
		panic(err)
	}
	return path
}

func ItBehavesLikeHelp(command string, alias string, validate func(*Session)) {
	It("displays help", func() {
		session := runCommand(command, "-h")
		Eventually(session).Should(Exit(1))
		validate(session)
	})

	It("displays help using the alias", func() {
		session := runCommand(alias, "-h")
		Eventually(session).Should(Exit(1))
		validate(session)
	})
}
