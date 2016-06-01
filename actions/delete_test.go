package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/actions"

	"net/http"

	"errors"

	"bytes"
	"io/ioutil"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/client/clientfakes"
	"github.com/pivotal-cf/cm-cli/config"
	cm_errors "github.com/pivotal-cf/cm-cli/errors"
)

var _ = Describe("Delete", func() {

	var (
		subject    Delete
		httpClient clientfakes.FakeHttpClient
	)

	BeforeEach(func() {
		myConfig := config.Config{ApiURL: "pivotal.io"}
		subject = NewDelete(&httpClient, myConfig)
	})

	Describe("Delete", func() {
		It("deletes a secret from the server", func() {
			request := client.NewDeleteSecretRequest("pivotal.io", "my-secret")

			responseObj := http.Response{
				StatusCode: 200,
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			err := subject.Delete("my-secret")

			Expect(err).ToNot(HaveOccurred())
		})

		Describe("Errors", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewDelete(&httpClient, config.Config{})

				error := subject.Delete("my-secret")

				Expect(error).To(MatchError(cm_errors.NewNoTargetUrlError()))
			})

			It("returns NewNetworkError when there is a network error", func() {
				httpClient.DoReturns(nil, errors.New("hello"))

				error := subject.Delete("my-secret")
				Expect(error).To(MatchError(cm_errors.NewNetworkError()))
			})

			It("returns a not-found error when response is 404", func() {
				responseObj := http.Response{
					StatusCode: 404,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "My error"}`))),
				}

				httpClient.DoReturns(&responseObj, nil)

				error := subject.Delete("my-secret")
				Expect(error.Error()).To(Equal("My error"))
			})

			It("returns a bad request error when response is 500", func() {
				responseObj := http.Response{
					StatusCode: 500,
				}

				httpClient.DoReturns(&responseObj, nil)

				error := subject.Delete("my-secret")
				Expect(error).To(MatchError(cm_errors.NewSecretBadRequestError()))
			})
		})
	})
})
