package credhub_test

import (
	. "code.cloudfoundry.org/credhub-cli/credhub"

	"bytes"
	"io/ioutil"
	"net/http"

	"code.cloudfoundry.org/credhub-cli/credhub/permissions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Permissions", func() {
	Context("GetPermissions", func() {
		It("returns the permissions", func() {
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

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			actualPermissions, err := ch.GetPermission("1234")
			Expect(err).NotTo(HaveOccurred())

			expectedPermission := permissions.Permission{
				Actor:      "user:A",
				Operations: []string{"read"},
				Path:		"/example-password",
				UUID:		"1234",
			}
			Expect(actualPermissions).To(Equal(&expectedPermission))

			By("calling the right endpoints")
			url := dummyAuth.Request.URL.String()
			Expect(url).To(Equal("https://example.com/api/v2/permissions/1234"))
			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})
	})

	Context("AddPermissions", func() {
		Context("when a credential exists", func() {
			It("can add permissions to it", func() {
				responseString :=
					`{
	"actor":"user:B",
	"operations":["read"],
	"path":"/example-password"
}`
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusCreated,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}
				ch, _ := New("https://example.com", Auth(dummy.Builder()))

				_, err := ch.AddPermission("/example-password", "user:A", []string{"read_acl", "write_acl"})
				Expect(err).NotTo(HaveOccurred())

				By("calling the right endpoints")
				url := dummy.Request.URL.String()
				Expect(url).To(Equal("https://example.com/api/v2/permissions"))
				Expect(dummy.Request.Method).To(Equal(http.MethodPost))
				params, err := ioutil.ReadAll(dummy.Request.Body)
				Expect(err).NotTo(HaveOccurred())

				expectedParams := `{
				"actor": "user:A",
				"operations": ["read_acl", "write_acl"],
				"path": "/example-password"
			}`
				Expect(params).To(MatchJSON(expectedParams))
			})
		})
	})
})
