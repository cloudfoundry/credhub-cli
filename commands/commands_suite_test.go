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
const SECRET_STRING_OVERWRITE_REQUEST_JSON = `{"type":"%s","name":"%s","value":"%s","overwrite":%t}`
const SECRET_STRING_RESPONSE_JSON = `{"data":[{"type":"%s", "value":"%s", "updated_at":"` + TIMESTAMP + `"}]}`
const SECRET_SET_STRING_RESPONSE_JSON = `{"type":"%s", "value":"%s", "updated_at":"` + TIMESTAMP + `"}`
const SECRET_STRING_RESPONSE_TABLE = "Type:          %s\nName:          %s\nValue:         %s\nUpdated:       " + TIMESTAMP
const SECRET_CERTIFICATE_REQUEST_JSON = `{"type":"certificate","name":"%s","value":{"ca":"%s","certificate":"%s","private_key":"%s"},"overwrite":%t}`
const SECRET_CERTIFICATE_RESPONSE_JSON = `{"data":[{"type":"certificate","value":{"ca":"%s","certificate":"%s","private_key":"%s"},"updated_at":"` + TIMESTAMP + `"}]}`
const SECRET_CERTIFICATE_RESPONSE_TABLE = "Type:          certificate\nName:          %s\nCa:            %s\nCertificate:   %s\nPrivate Key:   %s\nUpdated:       " + TIMESTAMP
const SECRET_RSA_SSH_REQUEST_JSON = `{"type":"%s","name":"%s","value":{"public_key":"%s","private_key":"%s"},"overwrite":%t}`
const SECRET_RSA_SSH_RESPONSE_JSON = `{"data":[{"type":"%s","value":{"public_key":"%s","private_key":"%s"},"updated_at":"` + TIMESTAMP + `"}]}`
const SECRET_SSH_RESPONSE_TABLE = "Type:          ssh\nName:          %s\nPublic Key:    %s\nPrivate Key:   %s\nUpdated:       " + TIMESTAMP
const SECRET_RSA_RESPONSE_TABLE = "Type:          rsa\nName:          %s\nPublic Key:    %s\nPrivate Key:   %s\nUpdated:       " + TIMESTAMP
const GENERATE_SECRET_REQUEST_JSON = `{"type":"%s","overwrite":%t,"parameters":%s}`
const GENERATE_DEFAULT_TYPE_REQUEST_JSON = `{"type":"password","overwrite":%t,"parameters":%s}`
const CA_REQUEST_JSON = `{"name":"%s","type":"%s","value":{"certificate":"%s","private_key":"%s"}}`
const CA_SET_RESPONSE_JSON = `{"type":"%s","value":{"certificate":"%s","private_key":"%s"},"updated_at":"` + TIMESTAMP + `"}`
const CA_RESPONSE_JSON = `{"data":[{"type":"%s","value":{"certificate":"%s","private_key":"%s"},"updated_at":"` + TIMESTAMP + `"}]}`
const CA_RESPONSE_TABLE = "Type:          %s\nName:          %s\nCertificate:   %s\nPrivate Key:   %s\nUpdated:       " + TIMESTAMP

var responseMyValuePotatoes = fmt.Sprintf(SECRET_STRING_RESPONSE_TABLE, "value", "my-value", "potatoes")
var responseMyPasswordPotatoes = fmt.Sprintf(SECRET_STRING_RESPONSE_TABLE, "password", "my-password", "potatoes")
var responseMyCertificate = fmt.Sprintf(SECRET_CERTIFICATE_RESPONSE_TABLE, "my-secret", "my-ca", "my-cert", "my-priv")
var responseMyCertificateAuthority = fmt.Sprintf(CA_RESPONSE_TABLE, "root", "my-secret", "my-cert", "my-priv")
var responseMySSHFoo = fmt.Sprintf(SECRET_SSH_RESPONSE_TABLE, "foo-ssh-key", "some-public-key", "some-private-key")
var responseMyRSAFoo = fmt.Sprintf(SECRET_RSA_RESPONSE_TABLE, "foo-rsa-key", "some-public-key", "some-private-key")

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
					"app":{"version":"my-version","name":"Pivotal Credential Manager"},
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
	executable_path, err := Build("github.com/pivotal-cf/credhub-cli", "-ldflags", "-X github.com/pivotal-cf/credhub-cli/version.Version=test-version")
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
