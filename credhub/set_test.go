package credhub_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
)

var _ = Describe("Set", func() {

	withMetadata := func(metadata credentials.Metadata) func(s *SetOptions) error {
		return func(s *SetOptions) error {
			s.Metadata = metadata
			return nil
		}
	}

	Describe("SetCredential()", func() {
		It("returns the credential that has been set", func() {
			metadataStr := `{"some":{"json":"metadata"}}`
			var metadata credentials.Metadata
			Expect(json.Unmarshal([]byte(metadataStr), &metadata)).To(Succeed())

			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`{
		  "id": "some-credential",
		  "name": "some-credential",
		  "type": "some-type",
		  "value": "some-value",
          "metadata": {"some":{"json":"metadata"}},
		  "version_created_at": "2017-01-01T04:07:18Z"
		}`)),
			}}

			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).ToNot(HaveOccurred())

			cred, err := ch.SetCredential("some-credential", "some-type", "some-value", withMetadata(metadata))
			Expect(err).ToNot(HaveOccurred())

			requestBody := getBody(dummy.Request.Body)

			Expect(requestBody["name"]).To(Equal("some-credential"))
			Expect(requestBody["type"]).To(Equal("some-type"))
			Expect(requestBody["value"]).To(BeEquivalentTo("some-value"))
			Expect(requestBody["metadata"]).To(BeEquivalentTo(metadata))

			Expect(cred.Name).To(Equal("some-credential"))
			Expect(cred.Type).To(Equal("some-type"))
			Expect(cred.Value).To(BeAssignableToTypeOf(""))
			Expect(cred.Value).To(BeEquivalentTo("some-value"))
			Expect(cred.Metadata).To(Equal(metadata))
		})

		It("sends a request for server version when server version is not provided", func() {
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/info"),
					ghttp.RespondWith(http.StatusOK, `{}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/version"),
					ghttp.RespondWith(http.StatusOK, `{"version": "1.9.0"}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/api/v1/data"),
					ghttp.RespondWith(http.StatusOK, `{}`),
				),
			)
			ch, err := New(server.URL())
			Expect(err).ToNot(HaveOccurred())
			_, err = ch.SetCredential("some-credential", "some-type", "some-value")
			Expect(err).NotTo(HaveOccurred())
		})

		It("sets overwrite mode when server version is older than 2.x", func() {
			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("{}")),
			},
				Error: nil,
			}

			version := fmt.Sprintf("1.%d.0", rand.Intn(10))
			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion(version))
			Expect(err).ToNot(HaveOccurred())
			_, err = ch.SetCredential("some-credential", "some-type", "some-value")
			Expect(err).NotTo(HaveOccurred())

			var requestBody map[string]interface{}
			body, err := io.ReadAll(dummy.Request.Body)
			Expect(err).ToNot(HaveOccurred())
			Expect(json.Unmarshal(body, &requestBody)).To(Succeed())

			Expect(requestBody["mode"]).To(Equal("overwrite"))
		})

		It("returns an error when server version is invalid", func() {
			ch, err := New("https://example.com", ServerVersion("invalid-version"))
			Expect(err).ToNot(HaveOccurred())
			_, err = ch.SetCredential("some-credential", "some-type", "some-value")
			Expect(err).To(MatchError("malformed version: invalid-version"))
		})

		It("returns an error when request fails", func() {
			dummy := &DummyAuth{Error: errors.New("network error occurred")}
			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).ToNot(HaveOccurred())

			_, err = ch.SetCredential("some-credential", "some-type", "some-value")
			Expect(err).To(MatchError("network error occurred"))
		})

		It("returns an error when response body cannot be unmarshalled", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: io.NopCloser(bytes.NewBufferString("something-invalid")),
			}}
			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).ToNot(HaveOccurred())

			_, err = ch.SetCredential("some-credential", "some-type", "some-value")
			Expect(err).To(MatchError(ContainSubstring("invalid character 's'")))
		})
	})

	Describe("SetCertificate()", func() {
		It("returns the credential that has been set", func() {
			metadataStr := `{"some":{"json":"metadata"}}`
			var metadata credentials.Metadata
			Expect(json.Unmarshal([]byte(metadataStr), &metadata)).To(Succeed())

			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`{
		  "id": "some-id",
		  "name": "/example-certificate",
		  "type": "certificate",
		  "value": {
		    "ca": "some-ca",
		    "certificate": "some-certificate",
		    "private_key": "some-private-key"
		  },
          "metadata": {"some":{"json":"metadata"}},
		  "version_created_at": "2017-01-01T04:07:18Z"
		}`)),
			}}

			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())

			certificate := values.Certificate{
				Certificate: "some-certificate",
			}
			cred, err := ch.SetCertificate("/example-certificate", certificate, withMetadata(metadata))
			Expect(err).NotTo(HaveOccurred())

			requestBody := getBody(dummy.Request.Body)

			Expect(requestBody["name"]).To(Equal("/example-certificate"))
			Expect(requestBody["type"]).To(Equal("certificate"))
			Expect(requestBody["value"]).To(HaveKeyWithValue("ca", ""))
			Expect(requestBody["value"]).To(HaveKeyWithValue("certificate", "some-certificate"))
			Expect(requestBody["value"]).To(HaveKeyWithValue("private_key", ""))
			Expect(requestBody["metadata"]).To(BeEquivalentTo(metadata))

			Expect(cred.Name).To(Equal("/example-certificate"))
			Expect(cred.Type).To(Equal("certificate"))
			Expect(cred.Value.Ca).To(Equal("some-ca"))
			Expect(cred.Value.Certificate).To(Equal("some-certificate"))
			Expect(cred.Value.PrivateKey).To(Equal("some-private-key"))
			Expect(cred.Metadata).To(Equal(metadata))
		})
		It("returns an error when request fails", func() {
			dummy := &DummyAuth{Error: errors.New("network error occurred")}
			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())
			certificate := values.Certificate{
				Ca: "some-ca",
			}

			_, err = ch.SetCertificate("/example-certificate", certificate)
			Expect(err).To(MatchError("network error occurred"))
		})
	})

	Describe("SetPassword()", func() {
		It("returns the credential that has been set", func() {
			metadataStr := `{"some":{"json":"metadata"}}`
			var metadata credentials.Metadata
			Expect(json.Unmarshal([]byte(metadataStr), &metadata)).To(Succeed())

			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`{
		  "id": "some-id",
		  "name": "/example-password",
		  "type": "password",
		  "value": "some-password",
		  "metadata": {"some":{"json":"metadata"}},
		  "version_created_at": "2017-01-01T04:07:18Z"
		}`)),
			}}

			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())

			password := values.Password("some-password")

			cred, err := ch.SetPassword("/example-password", password, withMetadata(metadata))
			Expect(err).NotTo(HaveOccurred())

			requestBody := getBody(dummy.Request.Body)

			Expect(requestBody["name"]).To(Equal("/example-password"))
			Expect(requestBody["type"]).To(Equal("password"))
			Expect(requestBody["value"]).To(Equal("some-password"))
			Expect(requestBody["metadata"]).To(BeEquivalentTo(metadata))

			Expect(cred.Name).To(Equal("/example-password"))
			Expect(cred.Type).To(Equal("password"))
			Expect(cred.Value).To(BeEquivalentTo("some-password"))
			Expect(cred.Metadata).To(Equal(metadata))

		})
		It("returns an error when request fails", func() {
			dummy := &DummyAuth{Error: errors.New("network error occurred")}
			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())
			password := values.Password("some-password")

			_, err = ch.SetPassword("/example-password", password)
			Expect(err).To(MatchError("network error occurred"))
		})
	})

	Describe("SetUser()", func() {
		It("returns the credential that has been set", func() {
			metadataStr := `{"some":{"json":"metadata"}}`
			var metadata credentials.Metadata
			Expect(json.Unmarshal([]byte(metadataStr), &metadata)).To(Succeed())

			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
					{
						"id": "67fc3def-bbfb-4953-83f8-4ab0682ad675",
						"name": "/example-user",
						"type": "user",
						"value": {
							"username": "FQnwWoxgSrDuqDLmeLpU",
							"password": "6mRPZB3bAfb8lRpacnXsHfDhlPqFcjH2h9YDvLpL",
							"password_hash": "some-hash"
						},
		  				"metadata": {"some":{"json":"metadata"}},
						"version_created_at": "2017-01-05T01:01:01Z"
					}`)),
			}}

			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())

			user := values.User{Username: "some-username", Password: "some-password"}
			cred, err := ch.SetUser("/example-user", user, withMetadata(metadata))
			Expect(err).NotTo(HaveOccurred())

			requestBody := getBody(dummy.Request.Body)

			Expect(requestBody["name"]).To(Equal("/example-user"))
			Expect(requestBody["type"]).To(Equal("user"))
			Expect(requestBody["value"]).To(HaveKeyWithValue("username", "some-username"))
			Expect(requestBody["value"]).To(HaveKeyWithValue("password", "some-password"))
			Expect(requestBody["metadata"]).To(BeEquivalentTo(metadata))

			Expect(cred.Name).To(Equal("/example-user"))
			Expect(cred.Type).To(Equal("user"))
			alternateUsername := "FQnwWoxgSrDuqDLmeLpU"
			Expect(cred.Value.User).To(Equal(values.User{
				Username: alternateUsername,
				Password: "6mRPZB3bAfb8lRpacnXsHfDhlPqFcjH2h9YDvLpL",
			}))
			Expect(cred.Value.PasswordHash).To(Equal("some-hash"))
			Expect(cred.Metadata).To(Equal(metadata))
		})

		It("returns an error", func() {
			dummy := &DummyAuth{Error: errors.New("network error occurred")}
			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())
			user := values.User{Username: "username", Password: "some-password"}
			_, err = ch.SetUser("/example-user", user)
			Expect(err).To(MatchError("network error occurred"))
		})
	})

	Describe("SetRSA()", func() {
		It("returns the credential that has been set", func() {
			metadataStr := `{"some":{"json":"metadata"}}`
			var metadata credentials.Metadata
			Expect(json.Unmarshal([]byte(metadataStr), &metadata)).To(Succeed())

			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
					{
						"id": "67fc3def-bbfb-4953-83f8-4ab0682ad676",
						"name": "/example-rsa",
						"type": "rsa",
						"value": {
							"public_key": "public-key",
							"private_key": "private-key"
						},
		  				"metadata": {"some":{"json":"metadata"}},
						"version_created_at": "2017-01-01T04:07:18Z"
					}`)),
			}}

			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())

			cred, err := ch.SetRSA("/example-rsa", values.RSA{}, withMetadata(metadata))
			Expect(err).NotTo(HaveOccurred())

			requestBody := getBody(dummy.Request.Body)

			Expect(requestBody["name"]).To(Equal("/example-rsa"))
			Expect(requestBody["type"]).To(Equal("rsa"))
			Expect(requestBody["value"]).To(HaveKeyWithValue("private_key", ""))
			Expect(requestBody["value"]).To(HaveKeyWithValue("public_key", ""))
			Expect(requestBody["metadata"]).To(BeEquivalentTo(metadata))

			Expect(cred.Name).To(Equal("/example-rsa"))
			Expect(cred.Type).To(Equal("rsa"))
			Expect(cred.Value).To(Equal(values.RSA{
				PrivateKey: "private-key",
				PublicKey:  "public-key",
			}))
			Expect(cred.Metadata).To(Equal(metadata))
		})

		It("returns an error", func() {
			dummy := &DummyAuth{Error: errors.New("network error occurred")}
			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())
			_, err = ch.SetRSA("/example-rsa", values.RSA{})
			Expect(err).To(MatchError("network error occurred"))
		})
	})

	Describe("SetSSH()", func() {
		It("returns the credential that has been set", func() {
			metadataStr := `{"some":{"json":"metadata"}}`
			var metadata credentials.Metadata
			Expect(json.Unmarshal([]byte(metadataStr), &metadata)).To(Succeed())

			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
					{
						"id": "67fc3def-bbfb-4953-83f8-4ab0682ad676",
						"name": "/example-ssh",
						"type": "ssh",
						"value": {
							"public_key": "public-key",
							"private_key": "private-key"
						},
		  				"metadata": {"some":{"json":"metadata"}},
						"version_created_at": "2017-01-01T04:07:18Z"
					}`)),
			}}

			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())

			cred, err := ch.SetSSH("/example-ssh", values.SSH{}, withMetadata(metadata))
			Expect(err).NotTo(HaveOccurred())

			requestBody := getBody(dummy.Request.Body)

			Expect(requestBody["name"]).To(Equal("/example-ssh"))
			Expect(requestBody["type"]).To(Equal("ssh"))
			Expect(requestBody["value"]).To(HaveKeyWithValue("private_key", ""))
			Expect(requestBody["value"]).To(HaveKeyWithValue("public_key", ""))
			Expect(requestBody["metadata"]).To(BeEquivalentTo(metadata))

			Expect(cred.Name).To(Equal("/example-ssh"))
			Expect(cred.Type).To(Equal("ssh"))
			Expect(cred.Value.SSH).To(Equal(values.SSH{
				PrivateKey: "private-key",
				PublicKey:  "public-key",
			}))
			Expect(cred.Metadata).To(Equal(metadata))
		})

		It("returns an error", func() {
			dummy := &DummyAuth{Error: errors.New("network error occurred")}
			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())
			_, err = ch.SetSSH("/example-ssh", values.SSH{})
			Expect(err).To(MatchError("network error occurred"))
		})
	})

	Describe("SetJSON()", func() {
		It("returns the credential that has been set", func() {
			metadataStr := `{"some":{"json":"metadata"}}`
			var metadata credentials.Metadata
			Expect(json.Unmarshal([]byte(metadataStr), &metadata)).To(Succeed())

			JSONValue := `{
					"key": 123,
					"key_list": [
					  "val1",
					  "val2"
					],
					"is_true": true
				}`
			var unmarshalledJSONValue values.JSON
			Expect(json.Unmarshal([]byte(JSONValue), &unmarshalledJSONValue)).To(Succeed())

			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
					{
						"id": "some-id",
						"name": "/example-json",
						"type": "json",
						"value": %s,
					 	"metadata": {"some":{"json":"metadata"}},
						"version_created_at": "2017-01-01T04:07:18Z"
					}`, JSONValue))),
			}}

			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())

			cred, err := ch.SetJSON("/example-json", unmarshalledJSONValue, withMetadata(metadata))
			Expect(err).NotTo(HaveOccurred())

			requestBody := getBody(dummy.Request.Body)

			Expect(requestBody["name"]).To(Equal("/example-json"))
			Expect(requestBody["type"]).To(Equal("json"))
			Expect(requestBody["value"]).To(BeEquivalentTo(unmarshalledJSONValue))
			Expect(requestBody["metadata"]).To(BeEquivalentTo(metadata))

			Expect(cred.Name).To(Equal("/example-json"))
			Expect(cred.Type).To(Equal("json"))
			Expect(cred.Value).To(Equal(unmarshalledJSONValue))
			Expect(cred.Metadata).To(Equal(metadata))
		})

		It("returns an error when request fails", func() {
			dummy := &DummyAuth{Error: errors.New("network error occurred")}
			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())
			_, err = ch.SetJSON("/example-json", nil)
			Expect(err).To(MatchError("network error occurred"))
		})
	})

	Describe("SetValue()", func() {
		It("returns the credential that has been set", func() {
			metadataStr := `{"some":{"json":"metadata"}}`
			var metadata credentials.Metadata
			Expect(json.Unmarshal([]byte(metadataStr), &metadata)).To(Succeed())
			dummy := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
					{
						"id": "some-id",
						"name": "/example-value",
						"type": "value",
						"value": "some string value",
						"metadata": {"some":{"json":"metadata"}},
						"version_created_at": "2017-01-01T04:07:18Z"
					}`)),
			}}

			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())

			cred, err := ch.SetValue("/example-value", values.Value("some string value"), withMetadata(metadata))
			Expect(err).NotTo(HaveOccurred())

			requestBody := getBody(dummy.Request.Body)

			Expect(requestBody["name"]).To(Equal("/example-value"))
			Expect(requestBody["type"]).To(Equal("value"))
			Expect(requestBody["value"]).To(BeEquivalentTo("some string value"))
			Expect(requestBody["metadata"]).To(BeEquivalentTo(metadata))

			Expect(cred.Name).To(Equal("/example-value"))
			Expect(cred.Type).To(Equal("value"))
			Expect(cred.Value).To(Equal(values.Value("some string value")))
			Expect(cred.Metadata).To(Equal(metadata))
		})

		It("returns an error when request fails", func() {
			dummy := &DummyAuth{Error: errors.New("network error occurred")}
			ch, err := New("https://example.com", Auth(dummy.Builder()), ServerVersion("2.6.0"))
			Expect(err).NotTo(HaveOccurred())
			_, err = ch.SetValue("/example-value", values.Value(""))
			Expect(err).To(MatchError("network error occurred"))
		})
	})
})

func getBody(body io.ReadCloser) map[string]interface{} {
	var requestBody map[string]interface{}
	bodyBytes, err := io.ReadAll(body)
	Expect(err).ToNot(HaveOccurred())
	Expect(json.Unmarshal(bodyBytes, &requestBody)).To(Succeed())
	return requestBody
}
