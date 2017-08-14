package auth_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/authfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Uaa", func() {
	var (
		mockUaaClient *dummyUaaClient
	)

	BeforeEach(func() {
		mockUaaClient = &dummyUaaClient{}
	})

	Context("Do()", func() {
		It("should add the bearer token to the request header", func() {
			expectedResponse := &http.Response{StatusCode: 539, Body: ioutil.NopCloser(strings.NewReader(""))}
			expectedError := errors.New("some error")

			dc := &DummyClient{Response: expectedResponse, Error: expectedError}

			uaa := auth.Uaa{
				AccessToken: "some-access-token",
				ApiClient:   dc,
				UaaClient:   mockUaaClient,
			}

			request, _ := http.NewRequest("GET", "https://some-endpoint.com/path/", nil)

			actualResponse, actualError := uaa.Do(request)
			actualRequest := dc.Request

			authHeader := actualRequest.Header.Get("Authorization")
			Expect(authHeader).To(Equal("Bearer some-access-token"))
			Expect(actualRequest.Method).To(Equal("GET"))
			Expect(actualRequest.URL.String()).To(Equal("https://some-endpoint.com/path/"))

			Expect(actualResponse).To(BeIdenticalTo(expectedResponse))
			Expect(actualError).To(BeIdenticalTo(expectedError))
		})

		Context("when there is no access token", func() {
			It("should request an access token", func() {
				mockUaaClient.NewAccessToken = "new-access-token"
				mockUaaClient.NewRefreshToken = "new-refresh-token"

				dc := &DummyClient{Response: &http.Response{}, Error: nil}

				uaa := auth.Uaa{
					UaaClient:    mockUaaClient,
					ApiClient:    dc,
					ClientId:     "client-id",
					ClientSecret: "client-secret",
					Username:     "user-name",
					Password:     "user-password",
				}

				request, _ := http.NewRequest("GET", "https://some-endpoint.com/path/", nil)

				uaa.Do(request)

				Expect(dc.Request.Header.Get("Authorization")).To(Equal("Bearer new-access-token"))
				Expect(uaa.AccessToken).To(Equal("new-access-token"))
				Expect(uaa.RefreshToken).To(Equal("new-refresh-token"))
			})
		})

		Context("when the access token has expired", func() {
			It("should refresh the token and submit the request again", func() {
				fhc := &authfakes.FakeHttpClient{}
				fhc.DoStub = func(req *http.Request) (*http.Response, error) {
					resp := &http.Response{}
					if req.Header.Get("Authorization") != "Bearer new-access-token" {
						resp.StatusCode = 573
						resp.Body = ioutil.NopCloser(strings.NewReader(`{"error": "access_token_expired"}`))
					} else {
						resp.Body = ioutil.NopCloser(strings.NewReader(`Success!`))
					}
					return resp, nil
				}

				mockUaaClient.NewAccessToken = "new-access-token"
				mockUaaClient.NewRefreshToken = "new-refresh-token"

				uaa := auth.Uaa{
					AccessToken:  "old-access-token",
					RefreshToken: "old-refresh-token",
					ClientId:     "client-id",
					ClientSecret: "client-secret",
					ApiClient:    fhc,
					UaaClient:    mockUaaClient,
				}

				request, _ := http.NewRequest("GET", "https://some-endpoint.com/path/", nil)

				response, err := uaa.Do(request)

				Expect(err).ToNot(HaveOccurred())

				Expect(mockUaaClient.ClientId).To(Equal("client-id"))
				Expect(mockUaaClient.ClientSecret).To(Equal("client-secret"))
				Expect(mockUaaClient.RefreshToken).To(Equal("old-refresh-token"))

				body, err := ioutil.ReadAll(response.Body)

				Expect(err).ToNot(HaveOccurred())
				Expect(string(body)).To(Equal("Success!"))
			})
		})

		Context("when a non-auth error has occurred", func() {
			It("should forward the response untouched", func() {
				fhc := &authfakes.FakeHttpClient{}
				fhc.DoStub = func(req *http.Request) (*http.Response, error) {
					resp := &http.Response{}
					resp.StatusCode = 573
					resp.Body = ioutil.NopCloser(strings.NewReader(`{"error": "some other error"}`))
					return resp, nil
				}

				uaa := auth.Uaa{
					AccessToken:  "old-access-token",
					RefreshToken: "old-refresh-token",
					ClientId:     "client-id",
					ClientSecret: "client-secret",
					ApiClient:    fhc,
					UaaClient:    mockUaaClient,
				}

				request, _ := http.NewRequest("GET", "https://some-endpoint.com/path/", nil)

				response, err := uaa.Do(request)

				Expect(err).ToNot(HaveOccurred())

				body, err := ioutil.ReadAll(response.Body)

				Expect(err).ToNot(HaveOccurred())
				Expect(body).To(MatchJSON(`{"error": "some other error"}`))
			})
		})

	})

	Context("Refresh()", func() {
		BeforeEach(func() {
			mockUaaClient.NewAccessToken = "new-access-token"
			mockUaaClient.NewRefreshToken = "new-refresh-token"
		})

		Context("with a refresh token", func() {
			It("should make a refresh grant token request and save the new tokens", func() {
				uaa := auth.Uaa{
					ClientId:     "client-id",
					ClientSecret: "client-secret",
					RefreshToken: "some-refresh-token",
					UaaClient:    mockUaaClient,
				}

				uaa.Refresh()

				Expect(mockUaaClient.ClientId).To(Equal("client-id"))
				Expect(mockUaaClient.ClientSecret).To(Equal("client-secret"))
				Expect(mockUaaClient.RefreshToken).To(Equal("some-refresh-token"))

				Expect(uaa.AccessToken).To(Equal("new-access-token"))
				Expect(uaa.RefreshToken).To(Equal("new-refresh-token"))
			})

			Context("when the refresh token grant fails", func() {
				It("returns an error", func() {
					mockUaaClient.Error = errors.New("refresh token grant failed")

					uaa := auth.Uaa{
						ClientId:     "client-id",
						ClientSecret: "client-secret",
						RefreshToken: "some-refresh-token",
						UaaClient:    mockUaaClient,
					}

					err := uaa.Refresh()

					Expect(err).To(MatchError("refresh token grant failed"))

				})
			})
		})

		Context("without a refresh token", func() {
			Context("with a username and password", func() {
				It("should make a password grant request", func() {
					uaa := auth.Uaa{
						ClientId:     "client-id",
						ClientSecret: "client-secret",
						Username:     "user-name",
						Password:     "user-password",
						UaaClient:    mockUaaClient,
					}

					uaa.Refresh()

					Expect(mockUaaClient.ClientId).To(Equal("client-id"))
					Expect(mockUaaClient.ClientSecret).To(Equal("client-secret"))
					Expect(mockUaaClient.Username).To(Equal("user-name"))
					Expect(mockUaaClient.Password).To(Equal("user-password"))

					Expect(uaa.AccessToken).To(Equal("new-access-token"))
					Expect(uaa.RefreshToken).To(Equal("new-refresh-token"))
				})

				It("when performing the password grant returns an error", func() {
					mockUaaClient.Error = errors.New("password grant error")

					uaa := auth.Uaa{
						ClientId:     "client-id",
						ClientSecret: "client-secret",
						Username:     "user-name",
						Password:     "user-password",
						UaaClient:    mockUaaClient,
					}

					err := uaa.Refresh()
					Expect(err).To(MatchError("password grant error"))
				})
			})

			Context("with client credentials", func() {
				It("should make a client credentials grant request", func() {
					uaa := auth.Uaa{
						ClientId:     "client-id",
						ClientSecret: "client-secret",
						UaaClient:    mockUaaClient,
					}

					uaa.Refresh()

					Expect(mockUaaClient.ClientId).To(Equal("client-id"))
					Expect(mockUaaClient.ClientSecret).To(Equal("client-secret"))

					Expect(uaa.AccessToken).To(Equal("new-access-token"))
					Expect(uaa.RefreshToken).To(BeEmpty())
				})

				Context("when the client credentials grant fails", func() {
					It("returns an error", func() {
						mockUaaClient.Error = errors.New("client credentials grant failed")

						uaa := auth.Uaa{
							ClientId:     "client-id",
							ClientSecret: "client-secret",
							UaaClient:    mockUaaClient,
						}

						err := uaa.Refresh()

						Expect(err).To(MatchError("client credentials grant failed"))

					})
				})

			})
		})
	})

	Context("Login()", func() {
		BeforeEach(func() {
			mockUaaClient.NewAccessToken = "new-access-token"
			mockUaaClient.NewRefreshToken = "new-refresh-token"
		})

		Context("when there is already an access token", func() {
			It("should do nothing", func() {
				uaa := auth.Uaa{
					AccessToken: "some-access-token",
				}

				err := uaa.Login()

				Expect(err).To(BeNil())
			})
		})

		Context("with a username and password", func() {
			It("should make a password grant request", func() {
				uaa := auth.Uaa{
					ClientId:     "client-id",
					ClientSecret: "client-secret",
					Username:     "user-name",
					Password:     "user-password",
					UaaClient:    mockUaaClient,
				}

				uaa.Refresh()

				Expect(mockUaaClient.ClientId).To(Equal("client-id"))
				Expect(mockUaaClient.ClientSecret).To(Equal("client-secret"))
				Expect(mockUaaClient.Username).To(Equal("user-name"))
				Expect(mockUaaClient.Password).To(Equal("user-password"))

				Expect(uaa.AccessToken).To(Equal("new-access-token"))
				Expect(uaa.RefreshToken).To(Equal("new-refresh-token"))
			})

			Context("when the refresh token grant fails", func() {
				It("returns an error", func() {
					mockUaaClient.Error = errors.New("refresh token grant failed")

					uaa := auth.Uaa{
						ClientId:     "client-id",
						ClientSecret: "client-secret",
						RefreshToken: "some-refresh-token",
						UaaClient:    mockUaaClient,
					}

					err := uaa.Refresh()

					Expect(err).To(MatchError("refresh token grant failed"))

				})
			})
		})

		Context("with client credentials", func() {
			It("should make a client credentials grant request", func() {
				uaa := auth.Uaa{
					ClientId:     "client-id",
					ClientSecret: "client-secret",
					UaaClient:    mockUaaClient,
				}

				uaa.Refresh()

				Expect(mockUaaClient.ClientId).To(Equal("client-id"))
				Expect(mockUaaClient.ClientSecret).To(Equal("client-secret"))

				Expect(uaa.AccessToken).To(Equal("new-access-token"))
				Expect(uaa.RefreshToken).To(BeEmpty())
			})

			Context("when the client credentials grant fails", func() {
				It("returns an error", func() {
					mockUaaClient.Error = errors.New("client credentials grant failed")

					uaa := auth.Uaa{
						ClientId:     "client-id",
						ClientSecret: "client-secret",
						UaaClient:    mockUaaClient,
					}

					err := uaa.Refresh()

					Expect(err).To(MatchError("client credentials grant failed"))

				})
			})
		})
	})

})
