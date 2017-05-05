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
	cmcli_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
)

var _ = Describe("SecretRepository", func() {
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
			repository = NewSecretRepository(&httpClient)
		})

		Context("when there is a response body", func() {
			It("sends a request to the server which responds with a single credential", func() {
				request, _ := http.NewRequest("GET", "http://example.com/foo", nil)

				expectedSecretJson := `{"name":"foo","id":"some-id","type":"value","value":"my-value","version_created_at":"2016-12-07T22:57:04Z"}`
				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(expectedSecretJson))),
				}

				httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
					Expect(req).To(Equal(request))
					return &responseObj, nil
				}

				secret, err := repository.SendRequest(request, "foo")
				Expect(err).ToNot(HaveOccurred())
				Expect(secret.ToJson()).To(MatchJSON(expectedSecretJson))
			})

			It("sends a request to the server for an array of credentials", func() {
				request, _ := http.NewRequest("GET", "http://example.com/bar", nil)

				expectedSecretJson := `{"name":"bar","id":"some-id","type":"password","value":"my-password","version_created_at":"2016-12-07T22:57:04Z"}`
				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(expectedSecretJson))),
				}

				httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
					Expect(req).To(Equal(request))

					return &responseObj, nil
				}

				secret, err := repository.SendRequest(request, "foo")
				Expect(err).ToNot(HaveOccurred())
				Expect(secret.ToJson()).To(MatchJSON(expectedSecretJson))
			})
		})

		Describe("Deletion", func() {
			It("sends a delete request to the server", func() {
				request, _ := http.NewRequest("DELETE", "http://example.com/foo", nil)

				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
				}

				httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
					Expect(req).To(Equal(request))

					return &responseObj, nil
				}

				secret, err := repository.SendRequest(request, "foo")

				Expect(err).ToNot(HaveOccurred())
				Expect(secret).To(Equal(models.Secret{}))
			})
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
