//go:build !windows
// +build !windows

package config_test

import (
	"os"

	"path"

	"code.cloudfoundry.org/credhub-cli/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var cfg config.Config
	var cachedConfig config.Config

	BeforeEach(func() {
		cachedConfig = config.ReadConfig()
		cfg = config.Config{
			ConfigWithoutSecrets: config.ConfigWithoutSecrets{
				ApiURL:  "http://api.example.com",
				AuthURL: "http://auth.example.com",
			},
		}
	})

	AfterEach(func() {
		config.WriteConfig(cachedConfig)
	})

	It("set appropriate permissions for persisted files", func() {
		config.WriteConfig(cfg)

		parentStat, _ := os.Stat(path.Dir(config.ConfigPath()))
		Expect(parentStat.Mode().String()).To(Equal("drwxr-xr-x"))

		fileStat, _ := os.Stat(config.ConfigPath())
		Expect(fileStat.Mode().String()).To(Equal("-rw-------"))
	})

	Describe("CREDHUB_HOME", func() {
		var originalCredhubHome string
		var originalHome string
		var hadCredhubHome bool

		BeforeEach(func() {
			originalCredhubHome, hadCredhubHome = os.LookupEnv("CREDHUB_HOME")
			originalHome = os.Getenv("HOME")
		})

		AfterEach(func() {
			if hadCredhubHome {
				os.Setenv("CREDHUB_HOME", originalCredhubHome)
			} else {
				os.Unsetenv("CREDHUB_HOME")
			}
			os.Setenv("HOME", originalHome)
		})

		It("uses CREDHUB_HOME instead of HOME when set", func() {
			os.Setenv("HOME", "/home/original-user")
			os.Setenv("CREDHUB_HOME", "/custom/credhub/home")

			Expect(config.ConfigDir()).To(Equal("/custom/credhub/home/.credhub"))
		})

		It("falls back to HOME when CREDHUB_HOME is not set", func() {
			os.Setenv("HOME", "/home/original-user")
			os.Unsetenv("CREDHUB_HOME")

			Expect(config.ConfigDir()).To(Equal("/home/original-user/.credhub"))
		})
	})
})
