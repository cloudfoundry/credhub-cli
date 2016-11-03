package repositories_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/credhub-cli/repositories"

	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/pivotal-cf/credhub-cli/client/clientfakes"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/models"

	cmcli_errors "github.com/pivotal-cf/credhub-cli/errors"
)

var _ = Describe("FindRepository", func() {
	var (
		repository Repository
		httpClient clientfakes.FakeHttpClient
		cfg        config.Config
	)

	BeforeEach(func() {
		cfg = config.Config{
			ApiURL:  "http://example.com",
			AuthURL: "http://uaa.example.com",
		}
	})

	Describe("SendRequest", func() {
		BeforeEach(func() {
			repository = NewSecretQueryRepository(&httpClient)
		})

		It("sends a request to the server", func() {
			request, _ := http.NewRequest("GET", "http://example.com/api/v1/data?name-like=find-me", nil)

			responseObj := http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
					"credentials": [
							{
								"name": "dan.password",
								"updated_at": "2016-09-06T23:26:58Z"
							},
							{
								"name": "deploy1/dan/id.key",
								"updated_at": "2016-09-06T23:26:58Z"
							}
					]
				}`))),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			expectedFindResponseBody := models.SecretQueryResponseBody{
				Credentials: []models.SecretQueryCredential{
					{
						Name:      "dan.password",
						UpdatedAt: "2016-09-06T23:26:58Z",
					},
					{
						Name:      "deploy1/dan/id.key",
						UpdatedAt: "2016-09-06T23:26:58Z",
					},
				},
			}

			findResponseBody, err := repository.SendRequest(request, "")

			Expect(err).ToNot(HaveOccurred())
			Expect(findResponseBody).To(Equal(expectedFindResponseBody))
		})

		Describe("Errors", func() {
			It("returns a NewResponseError when the JSON response cannot be parsed", func() {
				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("adasdasdasasd"))),
				}
				httpClient.DoReturns(&responseObj, nil)
				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

				_, error := repository.SendRequest(request, "foo")
				Expect(error).To(MatchError(cmcli_errors.NewResponseError()))
			})
		})
	})

})
