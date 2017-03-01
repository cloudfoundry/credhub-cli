package actions_test

import (
	. "github.com/cloudfoundry-incubator/credhub-cli/actions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"

	"errors"

	"bytes"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	cm_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories/repositoriesfakes"
)

var _ = Describe("Action", func() {

	var (
		subject      Action
		repository   repositoriesfakes.FakeRepository
		cfg          config.Config
		expectedBody models.SecretBody
		expectedItem models.Printable
	)

	BeforeEach(func() {
		cfg = config.Config{
			ApiURL:  "pivotal.io",
			AuthURL: "example.com",
		}
		subject = NewAction(&repository, cfg)
		expectedBody = models.SecretBody{
			Name:       "my-item",
			SecretType: "value",
			Value:      "potatoes",
		}
		expectedItem = models.Secret{
			SecretBody: expectedBody,
		}
	})

	AfterEach(func() {
		config.RemoveConfig()
	})

	Describe("DoAction", func() {
		It("performs a network request", func() {
			request, _ := http.NewRequest("POST", "my-url", nil)
			repository.SendRequestStub = func(req *http.Request, identifier string) (models.Printable, error) {
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

					repository.SendRequestStub = SequentialStub(
						func(req *http.Request, identifier string) (models.Printable, error) {
							buf := new(bytes.Buffer)
							buf.ReadFrom(req.Body)
							Expect(buf.String()).To(Equal("{}"))
							return nil, cm_errors.NewAccessTokenExpiredError()
						},
						func(req *http.Request, identifier string) (models.Printable, error) {
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
					cfg := config.ReadConfig()
					Expect(cfg.AccessToken).To(Equal("access_token"))
					Expect(cfg.RefreshToken).To(Equal("refresh_token"))
				})

				Context("after refreshing the token the repository returns an error", func() {
					It("refreshes the access token and returns repository error", func() {
						var authRepository repositoriesfakes.FakeRepository
						authRepository.SendRequestReturns(models.Token{AccessToken: "access_token", RefreshToken: "refresh_token"}, nil)
						subject.AuthRepository = &authRepository
						expectedError := errors.New("Custom Server Error")

						repository.SendRequestStub = SequentialStub(
							func(req *http.Request, identifier string) (models.Printable, error) {
								return nil, cm_errors.NewAccessTokenExpiredError()
							},
							func(req *http.Request, identifier string) (models.Printable, error) {
								Expect(req.Header.Get("Authorization")).To(Equal("Bearer access_token"))
								return nil, expectedError
							},
						)

						req, _ := http.NewRequest("POST", "my-url", bytes.NewBufferString("{}"))

						_, err := subject.DoAction(req, "my-item")

						Expect(err).To(HaveOccurred())
						Expect(expectedError).To(Equal(err))
						cfg := config.ReadConfig()
						Expect(cfg.AccessToken).To(Equal("access_token"))
						Expect(cfg.RefreshToken).To(Equal("refresh_token"))
					})
				})
			})
		})
	})
})

type RepositoryStub func(req *http.Request, identifier string) (models.Printable, error)

func SequentialStub(stubs ...RepositoryStub) RepositoryStub {
	return func(req *http.Request, identifier string) (models.Printable, error) {
		var s RepositoryStub
		s, stubs = stubs[0], stubs[1:]
		return s(req, identifier)
	}
}
