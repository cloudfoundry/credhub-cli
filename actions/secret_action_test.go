package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/actions"

	"net/http"

	"github.com/pivotal-cf/cm-cli/config"
	cm_errors "github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories/repositoriesfakes"
)

var _ = Describe("SecretAction", func() {

	var (
		subject          SecretAction
		secretRepository repositoriesfakes.FakeSecretRepository
	)

	BeforeEach(func() {
		myConfig := config.Config{ApiURL: "pivotal.io"}
		subject = NewSecretAction(&secretRepository, myConfig)
	})

	Describe("DoAction", func() {
		It("performs a network request", func() {
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

			secret, err := subject.DoSecretAction(request, "my-secret")

			Expect(err).ToNot(HaveOccurred())
			Expect(secret).To(Equal(expectedSecret))
		})

		Describe("Errors", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewSecretAction(&secretRepository, config.Config{})
				req, _ := http.NewRequest("GET", "my-url", nil)
				_, error := subject.DoSecretAction(req, "my-secret")

				Expect(error).To(MatchError(cm_errors.NewNoTargetUrlError()))
			})
		})
	})
})
