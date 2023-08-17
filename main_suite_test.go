package main_test

import (
	"runtime"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var (
	commandPath string
	commandName string
)

var _ = SynchronizedBeforeSuite(func() []byte {
	executable_path, err := Build("code.cloudfoundry.org/credhub-cli", "-ldflags", "-X code.cloudfoundry.org/credhub-cli/version.Version=test-version")
	Expect(err).NotTo(HaveOccurred())
	return []byte(executable_path)
}, func(data []byte) {
	commandPath = string(data)
	commandName = getLeafFileName(commandPath)
})

func getLeafFileName(path string) string {
	var pathArray []string
	if runtime.GOOS == "windows" {
		pathArray = strings.Split(path, "\\")
	} else {
		pathArray = strings.Split(path, "/")
	}
	return pathArray[len(pathArray)-1]
}
