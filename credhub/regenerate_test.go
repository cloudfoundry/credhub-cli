package credhub_test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
)

var _ = Describe("Regenerate", func() {
	It("regenerates the specified credential", func() {
		dummyAuth := &DummyAuth{Response: &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("")),
		}}

		ch, _ := New("https://example.com", Auth(dummyAuth.Builder()), ServerVersion("4.4.4"))

		ch.Regenerate("/example-password")
		url := dummyAuth.Request.URL.String()
		Expect(url).To(Equal("https://example.com/api/v1/regenerate"))
		Expect(dummyAuth.Request.Method).To(Equal(http.MethodPost))

		var requestBody map[string]interface{}
		body, _ := ioutil.ReadAll(dummyAuth.Request.Body)
		json.Unmarshal(body, &requestBody)

		Expect(requestBody["name"]).To(Equal("/example-password"))
	})

	Context("when successful", func() {
		It("returns the new credential", func() {
			responseString := `{
	  "id": "some-id",
	  "name": "/example-password",
	  "type": "password",
	  "value": "new-password",
	  "version_created_at": "2017-01-05T01:01:01Z"
	}`
			dummyAuth := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()), ServerVersion("4.4.4"))

			cred, err := ch.Regenerate("/example-password")
			Expect(err).To(BeNil())
			Expect(cred.Id).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-password"))
			Expect(cred.Type).To(Equal("password"))
			Expect(cred.Value.(string)).To(Equal("new-password"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
		})
	})

	Context("when response body cannot be unmarshalled", func() {
		It("returns an error", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()), ServerVersion("4.4.4"))
			_, err := ch.Regenerate("/example-password")

			Expect(err).To(HaveOccurred())
		})
	})

	Context("when the server version is older than 1.4.0", func() {
		It("regenerates the specified credential using the /data endpoint", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()), ServerVersion("1.3.0"))

			ch.Regenerate("/example-password")
			url := dummyAuth.Request.URL.String()
			Expect(url).To(Equal("https://example.com/api/v1/data"))
			Expect(dummyAuth.Request.Method).To(Equal(http.MethodPost))

			var requestBody map[string]interface{}
			body, _ := ioutil.ReadAll(dummyAuth.Request.Body)
			json.Unmarshal(body, &requestBody)

			Expect(requestBody["name"]).To(Equal("/example-password"))
			Expect(requestBody["regenerate"]).To(Equal(true))
		})
	})
})
