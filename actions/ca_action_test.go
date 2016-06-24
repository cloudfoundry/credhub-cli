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

var _ = Describe("CaAction", func() {

	var (
		subject      CaAction
		caRepository repositoriesfakes.FakeCaRepository
	)

	BeforeEach(func() {
		myConfig := config.Config{ApiURL: "pivotal.io"}
		subject = NewCaAction(&caRepository, myConfig)
	})

	Describe("DoAction", func() {
		It("performs a network request", func() {
			request, _ := http.NewRequest("GET", "my-url", nil)
			caParams := models.CaParameters{
				Public:  "my-pub",
				Private: "my-priv",
			}

			expectedBody := models.CaBody{
				Ca: &caParams,
			}
			expectedCa := models.NewCa("my-ca", expectedBody)
			caRepository.SendRequestStub = func(req *http.Request) (models.CaBody, error) {
				Expect(req).To(Equal(request))
				return expectedBody, nil
			}

			ca, err := subject.DoCaAction(request, "my-ca")

			Expect(err).ToNot(HaveOccurred())
			Expect(ca).To(Equal(expectedCa))
		})

		Describe("Errors", func() {
			It("returns a invalid target error when no api is set", func() {
				subject = NewCaAction(&caRepository, config.Config{})
				req, _ := http.NewRequest("GET", "my-url", nil)
				_, error := subject.DoCaAction(req, "my-ca")

				Expect(error).To(MatchError(cm_errors.NewNoTargetUrlError()))
			})
		})
	})
})
