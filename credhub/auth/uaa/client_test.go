package uaa_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Context("ClientCredentialGrant()", func() {
		It("should make a token grant request", func() {
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

				w.Write([]byte(`{"access_token": "access-token", "token_type": "bearer"}`))
			}))
			defer uaaServer.Close()

			client := Client{
				AuthUrl: uaaServer.URL,
				Client:  http.DefaultClient,
			}

			accessToken, err := client.ClientCredentialGrant("client-id", "client-secret")

			Expect(err).ToNot(HaveOccurred())
			Expect(accessToken).To(Equal("access-token"))
		})
	})

	Context("PasswordGrant()", func() {
		It("should make a password grant token request", func() {
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

			client := Client{
				AuthUrl: uaaServer.URL,
				Client:  http.DefaultClient,
			}

			accessToken, refreshToken, err := client.PasswordGrant("some-client-id", "some-client-secret", "username", "password")

			Expect(err).To(BeNil())
			Expect(accessToken).To(Equal("some-access-token"))
			Expect(refreshToken).To(Equal("some-refresh-token"))
		})
	})

	Context("RefreshTokenGrant()", func() {
		It("should make a refresh grant token request", func() {
			uaaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.ParseForm()

				Expect(r.Method).To(Equal(http.MethodPost))

				Expect(r.URL.Path).To(Equal("/oauth/token"))

				Expect(r.Header.Get("Accept")).To(Equal("application/json"))
				Expect(r.Header.Get("Content-Type")).To(Equal("application/x-www-form-urlencoded"))

				Expect(r.PostForm.Get("grant_type")).To(Equal("refresh_token"))
				Expect(r.PostForm.Get("response_type")).To(Equal("token"))

				Expect(r.PostForm.Get("client_id")).To(Equal("client-id"))
				Expect(r.PostForm.Get("client_secret")).To(Equal("client-secret"))

				Expect(r.PostForm.Get("refresh_token")).To(Equal("some-refresh-token"))

				w.Write([]byte(`{"access_token": "new-access-token", "refresh_token": "new-refresh-token", "token_type": "bearer"}`))
			}))
			defer uaaServer.Close()

			client := Client{
				AuthUrl: uaaServer.URL,
				Client:  http.DefaultClient,
			}

			accessToken, refreshToken, err := client.RefreshTokenGrant("client-id", "client-secret", "some-refresh-token")

			Expect(err).To(BeNil())
			Expect(accessToken).To(Equal("new-access-token"))
			Expect(refreshToken).To(Equal("new-refresh-token"))
		})
	})

	DescribeTable("unable to complete the request",
		func(performAction func(*Client) error) {
			uaaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Kill the connection without returning a status code
				hj, _ := w.(http.Hijacker)
				conn, _, _ := hj.Hijack()
				conn.Close()
			}))

			defer uaaServer.Close()

			client := &Client{
				AuthUrl: uaaServer.URL,
				Client:  http.DefaultClient,
			}

			err := performAction(client)

			Expect(err).To(HaveOccurred())
		},
		Entry("client credentials", func(c *Client) error {
			_, err := c.ClientCredentialGrant("client-id", "client-secret")
			return err
		}),
		Entry("password grant", func(c *Client) error {
			_, _, err := c.PasswordGrant("some-client-id", "some-client-secret", "username", "password")
			return err
		}),
		Entry("refresh token grant", func(c *Client) error {
			_, _, err := c.RefreshTokenGrant("client-id", "client-secret", "some-refresh-token")
			return err
		}),
	)

	DescribeTable("response body is invalid",
		func(performAction func(*Client) error) {
			uaaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{`))
			}))

			defer uaaServer.Close()

			client := &Client{
				AuthUrl: uaaServer.URL,
				Client:  http.DefaultClient,
			}

			err := performAction(client)
			Expect(err).To(HaveOccurred())
		},

		Entry("client credentials", func(c *Client) error {
			_, err := c.ClientCredentialGrant("client-id", "client-secret")
			return err
		}),
		Entry("password grant", func(c *Client) error {
			_, _, err := c.PasswordGrant("some-client-id", "some-client-secret", "username", "password")
			return err
		}),
		Entry("refresh token grant", func(c *Client) error {
			_, _, err := c.RefreshTokenGrant("client-id", "client-secret", "some-refresh-token")
			return err
		}),
	)
})
