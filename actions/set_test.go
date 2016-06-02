package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/actions"

	"net/http"

	"errors"

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

	Describe("SetValueSecret", func() {
		It("sets and returns a secret from the server", func() {
			request := client.NewPutValueRequest("pivotal.io", "my-secret", "abcd")
			expectedBody := models.SecretBody{
				ContentType: "value",
				Value:       "abcd",
			}
			expectedSecret := models.NewSecret("my-secret", expectedBody)
			secretRepository.SendRequestStub = func(req *http.Request) (models.SecretBody, error) {
				Expect(req).To(Equal(request))
				return expectedBody, nil
			}

			secret, err := subject.SetValue("my-secret", "abcd")

			Expect(err).ToNot(HaveOccurred())
			Expect(secret).To(Equal(expectedSecret))
		})

		Describe("Errors", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewSet(&secretRepository, config.Config{})

				_, error := subject.SetValue("my-secret", "abcd")

				Expect(error).To(MatchError(cmcli_errors.NewNoTargetUrlError()))
			})

			It("returns an error if the request fails", func() {
				request := client.NewPutValueRequest("pivotal.io", "my-secret", "abcd")
				expectedError := errors.New("My Special Error")
				secretRepository.SendRequestStub = func(req *http.Request) (models.SecretBody, error) {
					Expect(req).To(Equal(request))
					return models.SecretBody{}, expectedError
				}

				_, err := subject.SetValue("my-secret", "abcd")

				Expect(err).To(Equal(expectedError))
			})
		})
	})
})
