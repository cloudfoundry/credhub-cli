package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"
	"testing"

	"io/ioutil"
	"os"
	"runtime"

	"net/http"

	"fmt"

	"io"

	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

const TIMESTAMP = `2016-01-01T12:00:00Z`
const SECRET_VALUE_REQUEST_JSON = `{"type":"value", "value":"%s"}`
const SECRET_VALUE_RESPONSE_JSON = `{"type":"value", "value":"%s", "updated_at":"` + TIMESTAMP + `"}`
const SECRET_VALUE_RESPONSE_TABLE = `Type:		value\nName:		%s\nValue:	%s\nUpdated:	` + TIMESTAMP
const SECRET_CERTIFICATE_REQUEST_JSON = `{"type":"certificate","value":{"root":"%s","certificate":"%s","private":"%s"}}`
const SECRET_CERTIFICATE_RESPONSE_JSON = `{"type":"certificate","value":{"root":"%s","certificate":"%s","private":"%s"},"updated_at":"` + TIMESTAMP + `"}`
const SECRET_CERTIFICATE_RESPONSE_TABLE = `Type:		certificate\nName:		%s\nRoot:		%s\nCertificate:		%s\nPrivate:	%s\nUpdated:	` + TIMESTAMP
const GENERATE_REQUEST_JSON = `{"type":"%s","parameters":%s}`
const CA_REQUEST_JSON = `{"type":"%s","value":{"certificate":"%s","private":"%s"}}`
const CA_RESPONSE_JSON = `{"type":"%s","value":{"certificate":"%s","private":"%s"},"updated_at":"` + TIMESTAMP + `"}`
const CA_RESPONSE_TABLE = `Type:		%s\nName:		%s\nCertificate:		%s\nPrivate:	%s\nUpdated:	` + TIMESTAMP

var responseMySecretPotatoes = fmt.Sprintf(SECRET_VALUE_RESPONSE_TABLE, "my-secret", "potatoes")
var responseMySecretCertificate = fmt.Sprintf(SECRET_CERTIFICATE_RESPONSE_TABLE, "my-secret", "my-ca", "my-cert", "my-priv")
var responseMyCertificateAuthority = fmt.Sprintf(CA_RESPONSE_TABLE, "root", "my-secret", "my-cert", "my-priv")

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

	server = NewTLSServer()

	authServer = NewTLSServer()

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
	executable_path, err := Build("github.com/pivotal-cf/cm-cli")
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
