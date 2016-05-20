package actions_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/cm-cli/actions"
	"github.com/pivotal-cf/cm-cli/actions/actionsfakes"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
)

var _ = Describe("Version", func() {
	var (
		subject    actions.Version
		httpClient actionsfakes.FakeHttpClient
	)

	BeforeEach(func() {
		config := config.Config{ApiURL: "omfgdogs.com"}
		subject = actions.NewVersion(&httpClient, config)
	})

	Describe("Version", func() {
		It("returns the version of the cli and server", func() {
			request := client.NewInfoRequest("omfgdogs.com")

			responseObj := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"app":{"name":"Pivotal Credential Manager","version":"my-version"}}`))),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			serverVersion := subject.GetServerVersion()
			Expect(serverVersion).To(Equal("my-version"))
		})

		It("returns Not Found if server returned a non 200 status code", func() {
			responseObj := http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"app":{"name":"Pivotal Credential Manager","version":"my-version"}}`))),
			}

			httpClient.DoReturns(&responseObj, nil)

			serverVersion := subject.GetServerVersion()
			Expect(serverVersion).To(Equal("Not Found"))
		})

		It("returns Not Found if server has a network error", func() {
			responseObj := http.Response{
				StatusCode: 200,
			}

			httpClient.DoReturns(&responseObj, errors.New("dogs are gone"))

			serverVersion := subject.GetServerVersion()
			Expect(serverVersion).To(Equal("Not Found"))
		})

		It("returns Not Found if server returns bad json", func() {
			responseObj := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`sdafasdfasdf`))),
			}

			httpClient.DoReturns(&responseObj, nil)

			serverVersion := subject.GetServerVersion()
			Expect(serverVersion).To(Equal("Not Found"))
		})
	})
})
