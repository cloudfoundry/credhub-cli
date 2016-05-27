package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/actions"

	"net/http"

	"errors"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/client/clientfakes"
	cm_errors "github.com/pivotal-cf/cm-cli/errors"
)

var _ = Describe("Api", func() {

	var (
		subject    Api
		httpClient clientfakes.FakeHttpClient
	)

	BeforeEach(func() {
		subject = NewApi(&httpClient)
	})

	Describe("Api Validation", func() {
		It("returns no error when server is valid", func() {
			expectedRequest := client.NewInfoRequest("pivotal.io")

			responseObj := http.Response{
				StatusCode: 200,
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(expectedRequest))

				return &responseObj, nil
			}

			err := subject.ValidateTarget("pivotal.io")

			Expect(err).ToNot(HaveOccurred())
		})

		Describe("Errors", func() {
			It("returns our network error message when there is a network error", func() {
				expectedError := errors.New("hello")
				httpClient.DoReturns(nil, expectedError)

				err := subject.ValidateTarget("pivotal.io")

				Expect(err).To(Equal(cm_errors.NewNetworkError()))
			})

			It("returns an invalid target error when status code is not 200", func() {
				responseObj := http.Response{
					StatusCode: 404,
				}

				httpClient.DoReturns(&responseObj, nil)

				err := subject.ValidateTarget("pivotal.io")

				Expect(err).To(Equal(cm_errors.NewInvalidTargetError()))
			})
		})
	})
})
