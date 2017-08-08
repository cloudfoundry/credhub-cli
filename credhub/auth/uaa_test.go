package auth_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Uaa", func() {
	Context("Login()", func() {
		Context("Client Credential grant", func() {
			It("should make a client credential grant token request and save the access token", func() {
				uaaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.ParseForm()

					Expect(r.Method).To(Equal(http.MethodPost))

					Expect(r.URL.Path).To(Equal("/oauth/token"))

					Expect(r.Header.Get("Accept")).To(Equal("application/json"))
					Expect(r.Header.Get("Content-Type")).To(Equal("application/x-www-form-urlencoded"))

					Expect(r.PostForm.Get("grant_type")).To(Equal("client_credentials"))
					Expect(r.PostForm.Get("response_type")).To(Equal("token"))

					Expect(r.PostForm.Get("client_id")).To(Equal("client-id"))
					Expect(r.PostForm.Get("client_secret")).To(Equal("client-secret"))

					w.Write([]byte(`{"access_token": "some-access-token", "token_type": "bearer"}`))
				}))
				defer uaaServer.Close()

				uaa := auth.Uaa{
					ClientId:     "client-id",
					ClientSecret: "client-secret",
					AuthUrl:      uaaServer.URL,
					ApiClient:    http.DefaultClient,
				}

				uaa.Login()

				Expect(uaa.AccessToken).To(Equal("some-access-token"))
			})
		})
		Context("Password grant", func() {
			It("should make a password grant token request and save the access token", func() {
				uaaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.ParseForm()

					Expect(r.Method).To(Equal(http.MethodPost))

					Expect(r.URL.Path).To(Equal("/oauth/token"))

					Expect(r.Header.Get("Accept")).To(Equal("application/json"))
					Expect(r.Header.Get("Content-Type")).To(Equal("application/x-www-form-urlencoded"))

					Expect(r.PostForm.Get("grant_type")).To(Equal("password"))
					Expect(r.PostForm.Get("response_type")).To(Equal("token"))

					Expect(r.PostForm.Get("username")).To(Equal("username"))
					Expect(r.PostForm.Get("password")).To(Equal("password"))

					Expect(r.PostForm.Get("client_id")).To(Equal("some-client-id"))
					Expect(r.PostForm.Get("client_secret")).To(Equal("some-client-secret"))

					w.Write([]byte(`{"access_token": "some-access-token", "refresh_token": "some-refresh-token", "token_type": "bearer"}`))
				}))
				defer uaaServer.Close()

				uaa := auth.Uaa{
					Username:     "username",
					Password:     "password",
					ClientId:     "some-client-id",
					ClientSecret: "some-client-secret",
					AuthUrl:      uaaServer.URL,
					ApiClient:    http.DefaultClient,
				}

				uaa.Login()

				Expect(uaa.AccessToken).To(Equal("some-access-token"))
				Expect(uaa.RefreshToken).To(Equal("some-refresh-token"))
			})
		})
	})
})
