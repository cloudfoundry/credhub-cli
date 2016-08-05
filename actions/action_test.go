package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf/cm-cli/actions"

	"net/http"

	"errors"

	"bytes"

	"github.com/pivotal-cf/cm-cli/config"
	cm_errors "github.com/pivotal-cf/cm-cli/errors"
	"github.com/pivotal-cf/cm-cli/models"
	"github.com/pivotal-cf/cm-cli/repositories/repositoriesfakes"
	"github.com/pivotal-cf/cm-cli/util"
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
			request, _ := http.NewRequest("POST", "my-url", nil)
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
				req, _ := http.NewRequest("POST", "my-url", bytes.NewBufferString("{}"))
				_, error := subject.DoAction(req, "my-item")

				Expect(error).To(MatchError(cm_errors.NewNoTargetUrlError()))
			})

			Context("when repository returns unauthorized", func() {
				It("refreshes the access token", func() {
					var authRepository repositoriesfakes.FakeRepository
					authRepository.SendRequestReturns(models.Token{AccessToken: "access_token", RefreshToken: "refresh_token"}, nil)
					subject.AuthRepository = &authRepository

					repository.SendRequestStub = util.SequentialStub(
						func(req *http.Request, identifier string) (models.Item, error) {
							buf := new(bytes.Buffer)
							buf.ReadFrom(req.Body)
							Expect(buf.String()).To(Equal("{}"))
							return models.NewItem(), cm_errors.NewUnauthorizedError()
						},
						func(req *http.Request, identifier string) (models.Item, error) {
							Expect(req.Header.Get("Authorization")).To(Equal("Bearer access_token"))

							buf := new(bytes.Buffer)
							buf.ReadFrom(req.Body)
							Expect(buf.String()).To(Equal("{}"))

							return expectedItem, nil
						},
					)

					req, _ := http.NewRequest("POST", "my-url", bytes.NewBufferString("{}"))

					secret, err := subject.DoAction(req, "my-item")

					Expect(err).ToNot(HaveOccurred())
					Expect(secret).To(Equal(expectedItem))
					cfg, _ := config.ReadConfig()
					Expect(cfg.AccessToken).To(Equal("access_token"))
					Expect(cfg.RefreshToken).To(Equal("refresh_token"))
				})

				Context("after refreshing the token the repository returns an error", func() {
					It("refreshes the access token and returns repository error", func() {
						var authRepository repositoriesfakes.FakeRepository
						authRepository.SendRequestReturns(models.Token{AccessToken: "access_token", RefreshToken: "refresh_token"}, nil)
						subject.AuthRepository = &authRepository
						expectedError := errors.New("Custom Server Error")

						repository.SendRequestStub = util.SequentialStub(
							func(req *http.Request, identifier string) (models.Item, error) {
								return models.NewItem(), cm_errors.NewUnauthorizedError()
							},
							func(req *http.Request, identifier string) (models.Item, error) {
								Expect(req.Header.Get("Authorization")).To(Equal("Bearer access_token"))
								return models.NewItem(), expectedError
							},
						)

						req, _ := http.NewRequest("POST", "my-url", bytes.NewBufferString("{}"))

						_, err := subject.DoAction(req, "my-item")

						Expect(err).To(HaveOccurred())
						Expect(expectedError).To(Equal(err))
						cfg, _ := config.ReadConfig()
						Expect(cfg.AccessToken).To(Equal("access_token"))
						Expect(cfg.RefreshToken).To(Equal("refresh_token"))
					})
				})
			})
		})
	})
})
