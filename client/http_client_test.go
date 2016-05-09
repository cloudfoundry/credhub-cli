package client_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/cm-cli/client"
)

var _ = Describe("HttpClient", func() {
	var (
		server        *ghttp.Server
		c             *client.HttpClient
		serverHandler http.HandlerFunc

		route string

		requestStruct struct {
			SomeRequestField string `json:"SomeRequestField"`
		}
	)

	BeforeEach(func() {
		server = ghttp.NewServer()

		c = &client.HttpClient{
			BaseURL: server.URL(),
		}

		route = "/some/route"
	})

	JustBeforeEach(func() {
		server.AppendHandlers(serverHandler)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Put", func() {
		BeforeEach(func() {
			requestStruct.SomeRequestField = "SomeData"

			serverHandler = ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", route),
				ghttp.VerifyContentType("application/json"),
				ghttp.VerifyJSON(` { "SomeRequestField": "SomeData" } `),
				ghttp.RespondWith(http.StatusOK, `{ "SomeResponseField": "some value" }`),
			)
		})

		It("should send the request body as JSON", func() {
			var responseStruct struct {
				SomeResponseField string
			}
			err := c.Put(route, requestStruct, &responseStruct)
			Expect(err).NotTo(HaveOccurred())
			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(responseStruct.SomeResponseField).To(Equal("some value"))
		})
	})

})
