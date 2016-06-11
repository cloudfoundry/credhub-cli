package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/actions"

	"net/http"

	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/config"
	cmcli_errors "github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories/repositoriesfakes"
)

var _ = Describe("Delete", func() {

	var (
		subject          Delete
		secretRepository repositoriesfakes.FakeSecretRepository
	)

	BeforeEach(func() {
		myConfig := config.Config{ApiURL: "pivotal.io"}
		subject = NewDelete(&secretRepository, myConfig)
	})

	Describe("Delete", func() {
		It("deletes a secret from the server", func() {
			request, _ := http.NewRequest("DELETE", "my-url", nil)
			secretRepository.SendRequestStub = func(req *http.Request) (models.SecretBody, error) {
				Expect(req).To(Equal(request))
				return models.SecretBody{}, nil
			}

			err := subject.Delete(request, "my-secret")

			Expect(err).ToNot(HaveOccurred())
		})

		Describe("Errors", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewDelete(&secretRepository, config.Config{})

				req := client.NewDeleteSecretRequest("pivotal.io", "my-secret")
				error := subject.Delete(req, "my-secret")

				Expect(error).To(MatchError(cmcli_errors.NewNoTargetUrlError()))
			})
		})
	})
})
