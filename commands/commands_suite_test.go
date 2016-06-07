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

	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

const VALUE_REQUEST_JSON = `{"type":"value", "value":"%s"}`
const VALUE_RESPONSE_JSON = VALUE_REQUEST_JSON
const VALUE_RESPONSE_TABLE = `Type:	value\nName:	%s\nValue:	%s`
const CERTIFICATE_REQUEST_JSON = `{"type":"certificate","certificate":{"ca":"%s","public":"%s","private":"%s"}}`
const CERTIFICATE_RESPONSE_JSON = CERTIFICATE_REQUEST_JSON
const CERTIFICATE_RESPONSE_TABLE = `Type:		certificate\nName:		%s\nCA:		%s\nPublic:		%s\nPrivate:	%s`

var responseMyPotatoes = fmt.Sprintf(VALUE_RESPONSE_TABLE, "my-secret", "potatoes")

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var (
	commandPath string
	homeDir     string
	server      *Server
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

	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("GET", "/info"),
			RespondWith(http.StatusOK, ""),
		),
	)

	runCommand("api", server.URL())
})

var _ = AfterEach(func() {
	server.Close()
	os.RemoveAll(homeDir)
})

var _ = SynchronizedBeforeSuite(func() []byte {
	path, err := Build("github.com/pivotal-cf/cm-cli")
	Expect(err).NotTo(HaveOccurred())
	return []byte(path)
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
