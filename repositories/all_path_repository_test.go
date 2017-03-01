package repositories_test

import (
	. "github.com/cloudfoundry-incubator/credhub-cli/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/client/clientfakes"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/models"

	cmcli_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
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
			repository = NewAllPathRepository(&httpClient)
		})

		It("sends a request to the server", func() {
			request, _ := http.NewRequest("GET", "http://example.com/data?paths=true", nil)

			responseObj := http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
				"paths": [
						{
							"path": "deploy123/"
						},
						{
							"path": "deploy123/dan/"
						},
						{
							"path": "deploy123/dan/consul/"
						},
						{
							"path": "deploy12/"
						},
						{
							"path": "consul/"
						},
						{
							"path": "consul/deploy123/"
						}
				]
			}`))),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			expectedFindResponseBody := models.AllPathResponseBody{
				Paths: []models.Path{
					{
						Path: "deploy123/",
					},
					{
						Path: "deploy123/dan/",
					},
					{
						Path: "deploy123/dan/consul/",
					},
					{
						Path: "deploy12/",
					},
					{
						Path: "consul/",
					},
					{
						Path: "consul/deploy123/",
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
