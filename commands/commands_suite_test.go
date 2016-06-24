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

const TIMESTAMP = `2016-01-01T12:00:00Z`
const VALUE_REQUEST_JSON = `{"type":"value", "value":"%s"}`
const VALUE_RESPONSE_JSON = `{"type":"value", "value":"%s", "updated_at":"` + TIMESTAMP + `"}`
const VALUE_RESPONSE_TABLE = `Type:		value\nName:		%s\nValue:		%s\nUpdated:	` + TIMESTAMP
const CERTIFICATE_REQUEST_JSON = `{"type":"certificate","certificate":{"ca":"%s","public":"%s","private":"%s"}}`
const CERTIFICATE_RESPONSE_JSON = `{"type":"certificate","certificate":{"ca":"%s","public":"%s","private":"%s"},"updated_at":"` + TIMESTAMP + `"}`
const CERTIFICATE_RESPONSE_TABLE = `Type:		certificate\nName:		%s\nCA:		%s\nPublic:		%s\nPrivate:	%s\nUpdated:	` + TIMESTAMP
const CA_CERTIFICATE_REQUEST_JSON = `{"root":{"public":"%s","private":"%s"}}`
const CA_CERTIFICATE_RESPONSE_JSON = `{"root":{"public":"%s","private":"%s"},"updated_at":"` + TIMESTAMP + `"}`
const CA_CERTIFICATE_RESPONSE_TABLE = `Name:		%s\nPublic:		%s\nPrivate:	%s\nUpdated:	` + TIMESTAMP

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
