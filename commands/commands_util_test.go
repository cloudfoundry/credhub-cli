package commands_test

import (
	"os"

	"github.com/pivotal-cf/cm-cli/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmcli_errors "github.com/pivotal-cf/cm-cli/errors"
)

var _ = Describe("Util", func() {
	Describe("readFile", func() {
		It("reads a file into memory", func() {
			tempDir := createTempDir("filesForTesting")
			fileContents := "My Test String"
			filename := createSecretFile(tempDir, "file.txt", fileContents)
			readContents, err := commands.ReadFile(filename)
			Expect(readContents).To(Equal(fileContents))
			Expect(err).To(BeNil())
			os.RemoveAll(tempDir)
		})

		It("returns an error message if a file cannot be read", func() {
			readContents, err := commands.ReadFile("Foo")
			Expect(readContents).To(Equal(""))
			Expect(err).To(MatchError(cmcli_errors.NewFileLoadError()))
		})
	})
})
