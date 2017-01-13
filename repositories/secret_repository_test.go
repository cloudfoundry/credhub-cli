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
	cmcli_errors "github.com/pivotal-cf/credhub-cli/errors"
	"github.com/pivotal-cf/credhub-cli/models"
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

				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"name":"foo","id":"some-id","type":"value","value":"my-value","version_created_at":"2016-12-07T22:57:04Z"}`))),
				}

				httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
					Expect(req).To(Equal(request))

					return &responseObj, nil
				}

				expectedSecretBody := models.SecretBody{
					Name:             "foo",
					SecretType:       "value",
					Value:            "my-value",
					VersionCreatedAt: "2016-12-07T22:57:04Z",
				}

				expectedSecret := models.Secret{
					SecretBody: expectedSecretBody,
				}

				secret, err := repository.SendRequest(request, "foo")

				Expect(err).ToNot(HaveOccurred())
				Expect(secret).To(Equal(expectedSecret))
			})

			It("sends a request to the server for an array of credentials", func() {
				request, _ := http.NewRequest("GET", "http://example.com/bar", nil)

				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"data":[{"name":"bar","id":"some-id","type":"password","value":"my-password","version_created_at":"2016-12-07T22:57:04Z"}]}`))),
				}

				httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
					Expect(req).To(Equal(request))

					return &responseObj, nil
				}

				expectedSecretBody := models.SecretBody{
					Name:             "bar",
					SecretType:       "password",
					Value:            "my-password",
					VersionCreatedAt: "2016-12-07T22:57:04Z",
				}

				expectedSecret := models.Secret{
					SecretBody: expectedSecretBody,
				}

				secret, err := repository.SendRequest(request, "foo")

				Expect(err).ToNot(HaveOccurred())
				Expect(secret).To(Equal(expectedSecret))
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
