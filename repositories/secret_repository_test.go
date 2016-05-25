package repositories_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/repositories"

	"bytes"
	"io/ioutil"
	"net/http"

	"errors"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/client/clientfakes"
	cmcli_errors "github.com/pivotal-cf/cm-cli/errors"
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
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"value":"my-value"}`))),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			secret, err := subject.SendRequest(request)

			Expect(err).ToNot(HaveOccurred())
			Expect(secret).To(Equal(client.NewSecret("my-secret", client.SecretBody{Value: "my-value"})))
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
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(``))),
				}

				httpClient.DoReturns(&responseObj, nil)

				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)
				_, error := subject.SendRequest(request)

				Expect(error).To(MatchError(cmcli_errors.NewInvalidStatusError()))
			})

			It("returns a response error when response json cannot be parsed", func() {
				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("adasdasdasasd"))),
				}

				httpClient.DoReturns(&responseObj, nil)

				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)
				_, error := subject.SendRequest(request)

				Expect(error).To(Equal(cmcli_errors.NewResponseError()))
			})
		})
	})
})
