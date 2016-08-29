package repositories_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/credhub-cli/repositories"

	"bytes"
	"io/ioutil"
	"net/http"

	"errors"

	"github.com/pivotal-cf/credhub-cli/client/clientfakes"
	"github.com/pivotal-cf/credhub-cli/config"
	cmcli_errors "github.com/pivotal-cf/credhub-cli/errors"
)

var _ = Describe("Repository", func() {
	var (
		httpClient clientfakes.FakeHttpClient
		cfg        config.Config
	)

	BeforeEach(func() {
		cfg = config.Config{
			ApiURL:  "http://example.com",
			AuthURL: "http://uaa.example.com",
		}
	})

	Describe("DoSendRequest", func() {
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

			response, err := DoSendRequest(&httpClient, request)
			Expect(err).ToNot(HaveOccurred())
			Expect(response).To(Equal(&responseObj))
		})

		Describe("Errors", func() {
			It("returns NetworkError when there is a network error", func() {
				httpClient.DoReturns(nil, errors.New("hello"))

				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)
				_, error := DoSendRequest(&httpClient, request)
				Expect(error).To(MatchError(cmcli_errors.NewNetworkError()))
			})

			It("returns a error when response is 400", func() {
				responseObj := http.Response{
					StatusCode: 400,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "My error"}`))),
				}

				httpClient.DoReturns(&responseObj, nil)

				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)
				_, error := DoSendRequest(&httpClient, request)

				Expect(error.Error()).To(Equal("My error"))
			})

			It("returns UnauthorizedError when the CM server returns Unauthorized", func() {
				responseObj := http.Response{
					StatusCode: 401,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "invalid_token","error_description":"Access token expired: "}`))),
				}
				httpClient.DoReturns(&responseObj, nil)
				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

				_, error := DoSendRequest(&httpClient, request)
				Expect(error).To(MatchError(cmcli_errors.NewUnauthorizedError()))
			})

			It("returns ForbiddenError when the CM server returns Forbidden", func() {
				responseObj := http.Response{
					StatusCode: 403,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "insufficient_scope","error_description":"Insufficient scope for this resource"}`))),
				}
				httpClient.DoReturns(&responseObj, nil)
				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

				_, error := DoSendRequest(&httpClient, request)
				Expect(error).To(MatchError(cmcli_errors.NewForbiddenError()))
			})
		})
	})
})
