package credhub_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	. "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/permissions"
	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Permissions", func() {
	Context("GetPermissionByUUID", func() {
		Context("when server version is less than 2.0.0", func() {
			It("returns permission using V1 endpoint", func() {
				responseString :=
					`{
	"credential_name":"/test-password",
	"permissions":[{
			"actor":"user:A",
			"operations":["read"]
			}]
	}`

				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()), ServerVersion("1.9.0"))
				actualPermissions, err := ch.GetPermissions("1234")
				Expect(err).NotTo(HaveOccurred())

				expectedPermission := []permissions.V1_Permission{
					{
						Actor:      "user:A",
						Operations: []string{"read"},
					},
				}
				Expect(actualPermissions).To(Equal(expectedPermission))

				By("calling the right endpoints")
				url := dummyAuth.Request.URL.String()
				Expect(url).To(Equal("https://example.com/api/v1/permissions?credential_name=1234"))
				Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
			})
		})

		Context("when server version is greater than or equal to 2.0.0", func() {
			It("returns permission using V2 endpoint", func() {
				responseString :=
					`{
		"actor":"user:A",
		"operations":["read"],
		"path":"/example-password",
		"uuid":"1234"
	}`

				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()), ServerVersion("2.0.0"))
				actualPermissions, err := ch.GetPermissionByUUID("1234")
				Expect(err).NotTo(HaveOccurred())

				expectedPermission := permissions.Permission{
					Actor:      "user:A",
					Operations: []string{"read"},
					Path:       "/example-password",
					UUID:       "1234",
				}
				Expect(actualPermissions).To(Equal(&expectedPermission))

				By("calling the right endpoints")
				url := dummyAuth.Request.URL.String()
				Expect(url).To(Equal("https://example.com/api/v2/permissions/1234"))
				Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
			})
		})
	})

	Context("GetPermissionByPathActor", func() {
		It("correctly formats request", func() {
			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusCreated,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}}
			ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.1.2"))
			_, _ = ch.GetPermissionByPathActor("/path", "some-actor")
			request := dummy.Request.URL
			expectedRequest, _ := url.Parse("https://example.com/api/v2/permissions?actor=some-actor&path=/path")
			Expect(request.Query()).To(Equal(expectedRequest.Query()))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
		})

		Context("if permissions exists", func() {
			It("return permission using V2 endpoint", func() {
				responseString :=
					`{
		"actor":"user:A",
		"operations":["read"],
		"path":"/example-password",
		"uuid":"1234"
	}`

				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}
				ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.1.2"))

				actualPermissions, err := ch.GetPermissionByPathActor("/path", "some-actor")
				Expect(err).NotTo(HaveOccurred())

				expectedPermission := permissions.Permission{
					Actor:      "user:A",
					Operations: []string{"read"},
					Path:       "/example-password",
					UUID:       "1234",
				}

				Expect(actualPermissions).To(Equal(&expectedPermission))
			})
		})

		Context("if permissions does not exist", func() {
			It("returns a not found error", func() {

				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error": "The request could not be completed because the permission does not exist or you do not have sufficient authorization."}`)),
				}}
				ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.1.2"))

				_, err := ch.GetPermissionByPathActor("/path", "some-actor")
				Expect(err).To(BeAssignableToTypeOf(&NotFoundError{}))
			})
		})
		Context("if CredHub returns a 500 error", func() {
			It("returns a CredHub error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error": "some-error"}`)),
				}}
				ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.1.2"))

				_, err := ch.GetPermissionByPathActor("/path", "some-actor")
				Expect(err).To(BeAssignableToTypeOf(&Error{}))
			})
		})
	})

	Context("AddPermission", func() {
		Context("when server version is less than 2.0.0", func() {
			It("can add with V1 endpoint", func() {
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusCreated,
					Body:       ioutil.NopCloser(bytes.NewBufferString("")),
				}}
				ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("1.9.0"))

				_, err := ch.AddPermission("/example-password", "some-actor", []string{"read", "write"})

				Expect(err).NotTo(HaveOccurred())

				By("calling the right endpoints")
				url := dummy.Request.URL.String()
				Expect(url).To(Equal("https://example.com/api/v1/permissions"))
				Expect(dummy.Request.Method).To(Equal(http.MethodPost))
				params, err := ioutil.ReadAll(dummy.Request.Body)
				Expect(err).NotTo(HaveOccurred())

				expectedParams := `{
			  "credential_name": "/example-password",
			  "permissions": [
			  {
				"actor": "some-actor",
				"operations": ["read", "write"]
			  }]
			}`
				Expect(params).To(MatchJSON(expectedParams))
			})
		})

		Context("when server version is greater than or equal to 2.0.0", func() {
			var (
				responseString string
				ch             *CredHub
				dummy          *DummyAuth
			)
			BeforeEach(func() {
				responseString =
					`{
						"actor":"user:B",
						"operations":["read"],
						"path":"/example-password"
					}`
				dummy = &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusCreated,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}
				ch, _ = New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.0.0"))

			})

			It("can add with V2 endpoint", func() {
				_, err := ch.AddPermission("/example-password", "user:A", []string{"read", "read"})
				Expect(err).NotTo(HaveOccurred())

				By("calling the right endpoints")
				url := dummy.Request.URL.String()
				Expect(url).To(Equal("https://example.com/api/v2/permissions"))
				Expect(dummy.Request.Method).To(Equal(http.MethodPost))
				params, err := ioutil.ReadAll(dummy.Request.Body)
				Expect(err).NotTo(HaveOccurred())

				expectedParams := `{
				"actor": "user:A",
				"operations": ["read", "read"],
				"path": "/example-password"
			}`
				Expect(params).To(MatchJSON(expectedParams))

			})

			It("properly returns permissions", func() {
				permission, err := ch.AddPermission("/example-password", "user:B", []string{"read"})
				Expect(err).NotTo(HaveOccurred())
				expectedPermission := permissions.Permission{
					Actor:      "user:B",
					Path:       "/example-password",
					Operations: []string{"read"},
				}

				Expect(*permission).To(Equal(expectedPermission))

			})
		})

		Context("when server version is not specified", func() {
			var server *ghttp.Server

			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusOK, `{}`),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/version"),
						ghttp.RespondWith(http.StatusOK, `{"version": "2.0.0"}`),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/api/v2/permissions"),
						ghttp.VerifyJSON(`{
"actor": "some-actor",
"operations":["read", "write"],
"path":"/example-password"
}`),
						ghttp.RespondWith(http.StatusOK, `{}`),
					),
				)
			})

			AfterEach(func() {
				//shut down the server between tests
				server.Close()
			})

			It("can add permissions to it", func() {
				ch, _ := New(server.URL())
				_, err := ch.AddPermission("/example-password", "some-actor", []string{"read", "write"})
				Expect(err).NotTo(HaveOccurred())
				Expect(server.ReceivedRequests()).To(HaveLen(3))
			})
		})
	})

	Context("UpdatePermission", func() {
		Context("when server version is less than 2.0", func() {
			It("throws error", func() {
				ch, _ := New("https://example.com", ServerVersion("1.0.0"))
				_, err := ch.UpdatePermission("123", "path", "testactor", []string{"read"})
				Expect(err).To(HaveOccurred())
				Eventually(err).Should(Equal(fmt.Errorf("credhub server version <2.0 not supported")))
			})
		})
		Context("when server version is greater than or equal to 2.0", func() {
			responseString :=
				`{
						"actor":"user:B",
						"operations":["read"],
						"path":"/example-password",
				    "uuid":"1234"
					}`
			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusCreated,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
			}}
			ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.0.0"))

			It("properly returns permissions", func() {
				permission, err := ch.UpdatePermission("1234", "/example-password", "user:B", []string{"read"})
				Expect(err).NotTo(HaveOccurred())
				expectedPermission := permissions.Permission{
					Actor:      "user:B",
					Path:       "/example-password",
					Operations: []string{"read"},
					UUID:       "1234",
				}

				Expect(*permission).To(Equal(expectedPermission))

			})
		})
	})
	Context("Delete Permission", func() {
		It("correctly formats request", func() {
			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}}
			ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.1.2"))
			_, _ = ch.DeletePermission("1234")
			request := dummy.Request.URL
			expectedRequest, _ := url.Parse("https://example.com/api/v2/permissions/1234")
			Expect(request.String()).To(Equal(expectedRequest.String()))
			Expect(dummy.Request.Method).To(Equal(http.MethodDelete))
		})
		Context("when server version less than 2.0", func() {
			It("throws error", func() {
				ch, _ := New("https://example.com", ServerVersion("1.0.0"))
				_, err := ch.DeletePermission("123")
				Expect(err).To(HaveOccurred())
				Eventually(err).Should(Equal(fmt.Errorf("credhub server version <2.0 not supported")))
			})
		})
		Context("when server version is greater than or equal 2.0", func() {
			responseString :=
				`{
						"actor":"user:B",
						"operations":["read"],
						"path":"/example-password",
				    "uuid":"1234"
					}`
			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
			}}
			ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.0.0"))
			It("properly returns deleted permission", func() {
				permission, err := ch.DeletePermission("1234")
				Expect(err).NotTo(HaveOccurred())
				expectedPermission := permissions.Permission{
					Actor:      "user:B",
					Path:       "/example-password",
					Operations: []string{"read"},
					UUID:       "1234",
				}

				Expect(*permission).To(Equal(expectedPermission))
			})
		})
	})
})
