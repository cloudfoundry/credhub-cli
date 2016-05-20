package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/actions"

	"bytes"
	"io/ioutil"
	"net/http"

	"errors"

	"github.com/pivotal-cf/cm-cli/actions/actionsfakes"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	cmcli_errors "github.com/pivotal-cf/cm-cli/errors"
)

var _ = Describe("Set", func() {

	var (
		subject    Set
		httpClient actionsfakes.FakeHttpClient
	)

	BeforeEach(func() {
		config := config.Config{ApiURL: "pivotal.io"}

		subject = NewSet(&httpClient, config)
	})

	Describe("SetSecret", func() {
		It("sets and returns a secret from the server", func() {
			request := client.NewPutSecretRequest("pivotal.io", "my-secret", "my-value")

			responseObj := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"value":"my-value"}`))),
			}

			httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
				Expect(req).To(Equal(request))

				return &responseObj, nil
			}

			secret, err := subject.SetSecret("my-secret", "my-value")

			Expect(err).ToNot(HaveOccurred())
			Expect(secret).To(Equal(client.NewSecret("my-secret", client.SecretBody{Value: "my-value"})))
		})

		Describe("Errors", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewSet(&httpClient, config.Config{})

				_, error := subject.SetSecret("my-secret", "my-value")

				Expect(error).To(MatchError(cmcli_errors.NewNoTargetUrlError()))
			})

			It("returns NewNetworkError when there is a network error", func() {
				httpClient.DoReturns(nil, errors.New("hello"))

				_, error := subject.SetSecret("my-secret", "my-value")
				Expect(error).To(MatchError(cmcli_errors.NewNetworkError()))
			})

			It("returns a error when response is 400", func() {
				responseObj := http.Response{
					StatusCode: 400,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(``))),
				}

				httpClient.DoReturns(&responseObj, nil)

				_, error := subject.SetSecret("my-secret", "my-value")

				Expect(error).To(MatchError(cmcli_errors.NewInvalidStatusError()))
			})

			It("returns a response error when response json cannot be parsed", func() {
				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("adasdasdasasd"))),
				}

				httpClient.DoReturns(&responseObj, nil)

				_, error := subject.SetSecret("my-secret", "my-value")

				Expect(error).To(Equal(cmcli_errors.NewResponseError()))
			})
		})
	})
})
