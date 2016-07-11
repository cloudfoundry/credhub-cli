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

var _ = Describe("Action", func() {

	var (
		subject    Action
		repository repositoriesfakes.FakeRepository
	)

	BeforeEach(func() {
		myConfig := config.Config{ApiURL: "pivotal.io"}
		subject = NewAction(&repository, myConfig)
	})

	Describe("DoAction", func() {
		It("performs a network request", func() {
			request, _ := http.NewRequest("GET", "my-url", nil)
			expectedBody := models.SecretBody{
				ContentType: "value",
				Credential:  "potatoes",
			}
			expectedItem := models.NewSecret("my-item", expectedBody)
			repository.SendRequestStub = func(req *http.Request, identifier string) (models.Item, error) {
				Expect(req).To(Equal(request))
				return expectedItem, nil
			}

			secret, err := subject.DoAction(request, "my-item")

			Expect(err).ToNot(HaveOccurred())
			Expect(secret).To(Equal(expectedItem))
		})

		Describe("Errors", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewAction(&repository, config.Config{})
				req, _ := http.NewRequest("GET", "my-url", nil)
				_, error := subject.DoAction(req, "my-item")

				Expect(error).To(MatchError(cm_errors.NewNoTargetUrlError()))
			})
		})
	})
})
