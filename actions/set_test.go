package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/actions"

	"net/http"

	"errors"

	"bytes"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	cmcli_errors "github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories/repositoriesfakes"
)

var _ = Describe("Set", func() {

	var (
		subject          Set
		secretRepository repositoriesfakes.FakeSecretRepository
	)

	BeforeEach(func() {
		config := config.Config{ApiURL: "pivotal.io"}

		subject = NewSet(&secretRepository, config)
	})

	Describe("Set", func() {
		It("set uses supplied request and returns a SecretBody from the server", func() {
			myJsonRequest := `"{foo":"bar","obj":{"wild":"strawberries"}}`
			request, _ := http.NewRequest("PUT", "my-url", bytes.NewReader([]byte(myJsonRequest)))
			request.Header.Set("Content-Type", "application/json")

			expectedBody := models.SecretBody{
				ContentType: "value",
				Value:       "abcd",
				Certificate: &models.Certificate{
					Ca: "duh",
				},
			}
			expectedSecret := models.NewSecret("my-secret", expectedBody)
			secretRepository.SendRequestStub = func(req *http.Request) (models.SecretBody, error) {
				Expect(req).To(Equal(request))
				return expectedBody, nil
			}

			secret, err := subject.Set(request, "my-secret")

			Expect(err).ToNot(HaveOccurred())
			Expect(secret).To(Equal(expectedSecret))
		})

		Describe("Errors", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewSet(&secretRepository, config.Config{})

				req := client.NewPutValueRequest("pivotal.io", "my-secret", "abcd")
				_, error := subject.Set(req, "my-secret")

				Expect(error).To(MatchError(cmcli_errors.NewNoTargetUrlError()))
			})

			It("returns an error if the request fails", func() {
				request := client.NewPutValueRequest("pivotal.io", "my-secret", "abcd")
				expectedError := errors.New("My Special Error")
				secretRepository.SendRequestStub = func(req *http.Request) (models.SecretBody, error) {
					Expect(req).To(Equal(request))
					return models.SecretBody{}, expectedError
				}

				req := client.NewPutValueRequest("pivotal.io", "my-secret", "abcd")
				_, err := subject.Set(req, "my-secret")

				Expect(err).To(Equal(expectedError))
			})
		})
	})
})
