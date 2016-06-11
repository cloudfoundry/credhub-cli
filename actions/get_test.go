package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/actions"

	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	cm_errors "github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories/repositoriesfakes"
)

var _ = Describe("Get", func() {

	var (
		subject          Get
		secretRepository repositoriesfakes.FakeSecretRepository
	)

	BeforeEach(func() {
		myConfig := config.Config{ApiURL: "pivotal.io"}
		subject = NewGet(&secretRepository, myConfig)
	})

	Describe("GetSecret", func() {
		It("gets and returns a secret from the server", func() {
			request, _ := http.NewRequest("GET", "my-url", nil)
			expectedBody := models.SecretBody{
				ContentType: "value",
				Value:       "potatoes",
			}
			expectedSecret := models.NewSecret("my-secret", expectedBody)
			secretRepository.SendRequestStub = func(req *http.Request) (models.SecretBody, error) {
				Expect(req).To(Equal(request))
				return expectedBody, nil
			}

			secret, err := subject.GetSecret(request, "my-secret")

			Expect(err).ToNot(HaveOccurred())
			Expect(secret).To(Equal(expectedSecret))
		})

		Describe("Errors", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewGet(&secretRepository, config.Config{})
				req := client.NewGetSecretRequest("pivotal.io", "my-secret")
				_, error := subject.GetSecret(req, "my-secret")

				Expect(error).To(MatchError(cm_errors.NewNoTargetUrlError()))
			})
		})
	})
})
