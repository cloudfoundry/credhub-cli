package client_test

import (
	"net/http"

	. "github.com/pivotal-cf/cm-cli/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("API", func() {
	Describe("NewInfoRequest", func() {
		It("Returns a request for the info endpoint", func() {
			httpRequest, _ := http.NewRequest("GET", "fake_target.com/info", nil)
			apiTarget := "fake_target.com"

			request := NewInfoRequest(apiTarget)

			Expect(request).To(Equal(httpRequest))
		})
	})
})
