package credhub_test

import (
	"errors"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"

	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/permissions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddPermissions", func() {
	Context("when a credential exists", func() {
		It("can add permissions to it", func() {
			responseJSON := `{
				  "credential_name": "/example-password",
				  "permissions": [
					{
					  "actor": "some-actor",
					  "operations": [
						"operation-1",
						"operation-2"
					  ]
					},{
					  "actor": "already-existing-actor",
					  "operations": ["existing-operation"]
					}
				  ]
				}`

			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseJSON)),
			}}
			ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("1.2.3"))

			actualPermissions, err := ch.AddPermissions("/example-password", []permissions.Permission{
				{
					Actor:      "some-actor",
					Operations: []string{"operation-1", "operation-2"},
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(actualPermissions).To(HaveLen(2))
			Expect(actualPermissions[0].Actor).To(Equal("some-actor"))
			Expect(actualPermissions[0].Operations).To(Equal([]string{"operation-1", "operation-2"}))
			Expect(actualPermissions[1].Actor).To(Equal("already-existing-actor"))
			Expect(actualPermissions[1].Operations).To(Equal([]string{"existing-operation"}))

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
				"operations": ["operation-1", "operation-2"]
			  }]
			}`
			Expect(params).To(MatchJSON(expectedParams))
		})
	})

	Context("when a credential doesn't exist", func() {
		It("cannot add permissions to it", func() {
			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"The request could not be completed because the credential does not exist or you do not have sufficient authorization."}`)),
			}}
			ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("1.2.3"))

			_, err := ch.AddPermissions("/example-password", []permissions.Permission{
				{
					Actor:      "some-actor",
					Operations: []string{"operation-1", "operation-2"},
				},
			})

			Expect(err).To(MatchError(ContainSubstring("The request could not be completed because the credential does not exist or you do not have sufficient authorization.")))
		})
	})

	Context("when we can't read the response body", func() {
		It("wraps the reader error", func() {
			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(DeadReader{}),
			}}
			ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("1.2.3"))

			_, err := ch.AddPermissions("/example-password", []permissions.Permission{
				{
					Actor:      "some-actor",
					Operations: []string{"operation-1", "operation-2"},
				},
			})

			Expect(err).To(MatchError(ContainSubstring("cannot read response body in AddPermissions:")))
		})
	})

	Context("when the json response is invalid", func() {
		It("wraps the json error", func() {
			responseJSON := `{`

			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseJSON)),
			}}
			ch, _ := New("https://example.com", Auth(dummy.Builder()), ServerVersion("1.2.3"))

			_, err := ch.AddPermissions("/example-password", []permissions.Permission{
				{
					Actor:      "some-actor",
					Operations: []string{"operation-1", "operation-2"},
				},
			})

			Expect(err).To(MatchError(ContainSubstring("cannot unmarshal JSON in AddPermissions:")))
		})
	})
})

type DeadReader struct{}

func (_ DeadReader) Read(_ []byte) (int, error) {
	return 0, errors.New("Injected error")
}
