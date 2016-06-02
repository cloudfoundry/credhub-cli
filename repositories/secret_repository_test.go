package repositories_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/repositories"

	"bytes"
	"io/ioutil"
	"net/http"

	"errors"

	"github.com/pivotal-cf/cm-cli/client/clientfakes"
	cmcli_errors "github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
)

var _ = Describe("SecretRepository", func() {

	var (
		subject    SecretRepository
		httpClient clientfakes.FakeHttpClient
	)

	BeforeEach(func() {
		subject = NewSecretRepository(&httpClient)
	})

	Describe("SendRequest", func() {
		It("sends a request to the server", func() {
			request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

			responseObj := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"type":"value","value":"my-value"}`))),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			expectedSecretBody := models.SecretBody{
				ContentType: "value",
				Value:       "my-value",
			}

			secretBody, err := subject.SendRequest(request)

			Expect(err).ToNot(HaveOccurred())
			Expect(secretBody).To(Equal(expectedSecretBody))
		})

		Describe("Errors", func() {
			It("returns NewNetworkError when there is a network error", func() {
				httpClient.DoReturns(nil, errors.New("hello"))

				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)
				_, error := subject.SendRequest(request)
				Expect(error).To(MatchError(cmcli_errors.NewNetworkError()))
			})

			It("returns a error when response is 400", func() {
				responseObj := http.Response{
					StatusCode: 400,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "My error"}`))),
				}

				httpClient.DoReturns(&responseObj, nil)

				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)
				_, error := subject.SendRequest(request)

				Expect(error.Error()).To(Equal("My error"))
			})
		})
	})
})
