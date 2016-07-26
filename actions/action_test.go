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
		subject      Action
		repository   repositoriesfakes.FakeRepository
		cfg          config.Config
		expectedBody models.SecretBody
		expectedItem models.Item
	)

	BeforeEach(func() {
		cfg = config.Config{
			ApiURL:  "pivotal.io",
			AuthURL: "example.com",
		}
		subject = NewAction(&repository, cfg)
		expectedBody = models.SecretBody{
			ContentType: "value",
			Credential:  "potatoes",
		}
		expectedItem = models.NewSecret("my-item", expectedBody)
	})

	AfterEach(func() {
		config.RemoveConfig()
	})

	Describe("DoAction", func() {
		It("performs a network request", func() {
			request, _ := http.NewRequest("GET", "my-url", nil)
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

			It("refreshes the access token when server returns unauthorized", func() {
				var authRepository repositoriesfakes.FakeRepository
				authRepository.SendRequestReturns(models.Token{AccessToken: "access_token", RefreshToken: "refresh_token"}, nil)
				subject.AuthRepository = &authRepository

				i := 0
				repository.SendRequestStub = func(req *http.Request, identifier string) (models.Item, error) {
					i = i + 1
					if i == 1 {
						return models.NewItem(), cm_errors.NewUnauthorizedError()
					}
					Expect(req.Header.Get("Authorization")).To(Equal("Bearer access_token"))
					return expectedItem, nil
				}
				req, _ := http.NewRequest("GET", "my-url", nil)

				secret, err := subject.DoAction(req, "my-item")

				Expect(err).ToNot(HaveOccurred())
				Expect(secret).To(Equal(expectedItem))
				Expect(config.ReadConfig().AccessToken).To(Equal("access_token"))
				Expect(config.ReadConfig().RefreshToken).To(Equal("refresh_token"))
			})
		})
	})
})
