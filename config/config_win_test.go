//go:build windows
// +build windows

package config_test

import (
	"os"
	"syscall"

	"code.cloudfoundry.org/credhub-cli/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config (windows specific)", func() {
	var cfg config.Config

	BeforeEach(func() {
		cfg = config.Config{
			ConfigWithoutSecrets: config.ConfigWithoutSecrets{
				ApiURL:  "http://api.example.com",
				AuthURL: "http://auth.example.com",
			},
		}
	})

	It("hides the config directory", func() {
		err := config.WriteConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		p, err := syscall.UTF16PtrFromString(config.ConfigDir())
		Expect(err).ToNot(HaveOccurred())

		attrs, err := syscall.GetFileAttributes(p)
		Expect(err).ToNot(HaveOccurred())

		Expect(attrs & syscall.FILE_ATTRIBUTE_HIDDEN).To(Equal(uint32(syscall.FILE_ATTRIBUTE_HIDDEN)))
	})

	Describe("CREDHUB_HOME", func() {
		var originalCredhubHome string
		var originalUserProfile string
		var hadCredhubHome bool

		BeforeEach(func() {
			originalCredhubHome, hadCredhubHome = os.LookupEnv("CREDHUB_HOME")
			originalUserProfile = os.Getenv("USERPROFILE")
		})

		AfterEach(func() {
			if hadCredhubHome {
				os.Setenv("CREDHUB_HOME", originalCredhubHome)
			} else {
				os.Unsetenv("CREDHUB_HOME")
			}
			os.Setenv("USERPROFILE", originalUserProfile)
		})

		It("uses CREDHUB_HOME instead of USERPROFILE when set", func() {
			os.Setenv("USERPROFILE", `C:\Users\original-user`)
			os.Setenv("CREDHUB_HOME", `C:\custom\credhub\home`)

			Expect(config.ConfigDir()).To(Equal(`C:\custom\credhub\home\.credhub`))
		})

		It("falls back to USERPROFILE when CREDHUB_HOME is not set", func() {
			os.Setenv("USERPROFILE", `C:\Users\original-user`)
			os.Unsetenv("CREDHUB_HOME")

			Expect(config.ConfigDir()).To(Equal(`C:\Users\original-user\.credhub`))
		})
	})
})
