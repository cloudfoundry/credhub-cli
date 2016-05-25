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

	Describe("SetSecret", func() {
		It("sets and returns a secret from the server", func() {
			request := client.NewPutSecretRequest("pivotal.io", "my-secret", "my-value")
			expectedBody := client.SecretBody{Value: "my-value"}
			expectedSecret := client.NewSecret("my-secret", expectedBody)
			secretRepository.SendRequestStub = func(req *http.Request) (client.SecretBody, error) {
				Expect(req).To(Equal(request))
				return expectedBody, nil
			}

			secret, err := subject.SetSecret("my-secret", "my-value")

			Expect(err).ToNot(HaveOccurred())
			Expect(secret).To(Equal(expectedSecret))
		})

		It("generates a secret", func() {
			request := client.NewGenerateSecretRequest("pivotal.io", "my-name")
			expectedBody := client.SecretBody{Value: "my-value"}
			expectedSecret := client.NewSecret("my-name", expectedBody)
			secretRepository.SendRequestStub = func(req *http.Request) (client.SecretBody, error) {
				Expect(req).To(Equal(request))
				return expectedBody, nil
			}

			secret, err := subject.GenerateSecret("my-name")

			Expect(err).ToNot(HaveOccurred())
			Expect(secret).To(Equal(expectedSecret))
		})

		Describe("Errors when Setting", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewSet(&secretRepository, config.Config{})

				_, error := subject.SetSecret("my-secret", "my-value")

				Expect(error).To(MatchError(cmcli_errors.NewNoTargetUrlError()))
			})

			It("returns an error if the request fails", func() {
				request := client.NewPutSecretRequest("pivotal.io", "my-secret", "my-value")
				expectedError := errors.New("My Special Error")
				secretRepository.SendRequestStub = func(req *http.Request) (client.SecretBody, error) {
					Expect(req).To(Equal(request))
					return client.SecretBody{}, expectedError
				}

				_, err := subject.SetSecret("my-secret", "my-value")

				Expect(err).To(Equal(expectedError))
			})
		})

		Describe("Errors when Generating", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewSet(&secretRepository, config.Config{})

				_, error := subject.GenerateSecret("my-secret")

				Expect(error).To(MatchError(cmcli_errors.NewNoTargetUrlError()))
			})

			It("returns an error if the request fails", func() {
				request := client.NewGenerateSecretRequest("pivotal.io", "my-secret")
				expectedError := errors.New("My Special Error")
				secretRepository.SendRequestStub = func(req *http.Request) (client.SecretBody, error) {
					Expect(req).To(Equal(request))
					return client.SecretBody{}, expectedError
				}

				_, err := subject.GenerateSecret("my-secret")

				Expect(err).To(Equal(expectedError))
			})
		})
	})
})
