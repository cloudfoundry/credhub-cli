// +build !windows

package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/cm-cli/config"
)

var _ = Describe("Config", func() {
	var cfg config.Config

	BeforeEach(func() {
		cfg = config.Config{
			ApiURL:  "http://api.example.com",
			AuthURL: "http://auth.example.com",
		}
	})

	It("places the config file in .cm in the home directory", func() {
		Expect(config.ConfigPath()).To(HaveSuffix(`/.cm/config.json`))
	})
})
