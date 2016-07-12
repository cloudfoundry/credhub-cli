package actions_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/client/clientfakes"
	"github.com/pivotal-cf/cm-cli/config"
)

var _ = Describe("Info", func() {
	var (
		subject    actions.Version
		httpClient clientfakes.FakeHttpClient
	)

	BeforeEach(func() {
		config := config.Config{ApiURL: "omfgdogs.com"}
		subject = actions.NewInfo(&httpClient, config)
	})

	Describe("Version", func() {
		It("returns the version of the cli and server", func() {
			request := client.NewInfoRequest("omfgdogs.com")

			responseObj := http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"app":{"version":"my-version","name":"Pivotal Credential Manager"},
					"auth-server":{"url":"https://example.com"}
					}`)),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			serverInfo, _ := subject.GetServerInfo()
			Expect(serverInfo.App.Version).To(Equal("my-version"))
			Expect(serverInfo.AuthServer.Url).To(Equal("https://example.com"))
		})

		It("returns error if server returned a non 200 status code", func() {
			responseObj := http.Response{
				StatusCode: 400,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"app":{"version":"my-version","name":"Pivotal Credential Manager"},
					"auth-server":{"url":"https://example.com"}
					}`)),
			}

			httpClient.DoReturns(&responseObj, nil)

			_, err := subject.GetServerInfo()
			Expect(err).NotTo(BeNil())
		})

		It("returns error if server has a network error", func() {
			responseObj := http.Response{
				StatusCode: 200,
			}

			httpClient.DoReturns(&responseObj, errors.New("dogs are gone"))

			_, err := subject.GetServerInfo()
			Expect(err).NotTo(BeNil())
		})

		It("returns error if server returns bad json", func() {
			responseObj := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`sdafasdfasdf`)),
			}

			httpClient.DoReturns(&responseObj, nil)

			_, err := subject.GetServerInfo()
			Expect(err).NotTo(BeNil())
		})
	})
})
