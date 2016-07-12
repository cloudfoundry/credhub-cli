package repositories_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/repositories"

	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/pivotal-cf/cm-cli/client/clientfakes"
	"github.com/pivotal-cf/cm-cli/models"
)

var _ = Describe("CaRepository", func() {
	var (
		repository Repository
		httpClient clientfakes.FakeHttpClient
	)

	Describe("SendRequest", func() {
		Context("when there is a response body", func() {
			BeforeEach(func() {
				repository = NewCaRepository(&httpClient)
			})

			It("sends a request to the server", func() {
				request, _ := http.NewRequest("PUT", "http://example.com/foo", nil)

				responseObj := http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"ca":{"certificate":"my-cert","private":"my-priv"}}`))),
				}

				httpClient.DoStub = func(req *http.Request) (resp *http.Response, err error) {
					Expect(req).To(Equal(request))

					return &responseObj, nil
				}

				caParams := models.CaParameters{
					Certificate: "my-cert",
					Private:     "my-priv",
				}
				expectedCaBody := models.CaBody{
					Ca: &caParams,
				}
				expectedCa := models.Ca{
					Name:   "foo",
					CaBody: expectedCaBody,
				}

				ca, err := repository.SendRequest(request, "foo")

				Expect(err).ToNot(HaveOccurred())
				Expect(ca).To(Equal(expectedCa))
			})
		})
	})
})
