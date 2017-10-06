package credhub_test

import (
	version "github.com/hashicorp/go-version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
)

var _ = Describe("Options", func() {
	It("", func() {
		ch, _ := New("https://example.com", ServerVersion("2.2.2"))

		Expect(ch.ServerVersion()).To(Equal(version.Must(version.NewVersion("2.2.2"))))
	})
})
