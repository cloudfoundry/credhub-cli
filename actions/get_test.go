package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/actions"

	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/client/clientfakes"
	"github.com/pivotal-cf/cm-cli/config"
	cm_errors "github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
)

var _ = Describe("Get", func() {

	var (
		subject    Get
		httpClient clientfakes.FakeHttpClient
	)

	BeforeEach(func() {
		myConfig := config.Config{ApiURL: "pivotal.io"}
		subject = NewGet(&httpClient, myConfig)
	})

	Describe("GetSecret", func() {
		It("gets and returns a secret from the server", func() {
			request := client.NewGetSecretRequest("pivotal.io", "my-secret")
			expectedBody := models.SecretBody{
				ContentType: "value",
				Value:       "potatoes",
			}

			responseObj := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"type":"value","value":"potatoes"}`))),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			secret, err := subject.GetSecret("my-secret")

			Expect(err).ToNot(HaveOccurred())
			Expect(secret).To(Equal(models.NewSecret("my-secret", expectedBody)))
		})

		Describe("Errors", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewGet(&httpClient, config.Config{})

				_, error := subject.GetSecret("my-secret")

				Expect(error).To(MatchError(cm_errors.NewNoTargetUrlError()))
			})

			It("returns server error when response is 404", func() {
				responseObj := http.Response{
					StatusCode: 404,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "My error"}`))),
				}

				httpClient.DoReturns(&responseObj, nil)

				_, error := subject.GetSecret("my-secret")
				Expect(error.Error()).To(Equal("My error"))
			})

			It("returns NewNetworkError when there is a network error", func() {
				httpClient.DoReturns(nil, errors.New("hello"))

				_, error := subject.GetSecret("my-secret")
				Expect(error).To(MatchError(cm_errors.NewNetworkError()))
			})

			It("returns a response error when response json cannot be parsed", func() {
				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("adasdasdasasd"))),
				}

				httpClient.DoReturns(&responseObj, nil)

				_, error := subject.GetSecret("my-secret")
				Expect(error).To(MatchError(cm_errors.NewResponseError()))
			})
		})
	})
})
