package credhub_test

import (
	"bytes"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"encoding/json"

	. "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
)

var _ = Describe("Get", func() {
	var nilMetadata credentials.Metadata = nil

	Describe("GetLatestVersion()", func() {
		It("requests the credential by name using the 'current' query parameter", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

			ch.GetLatestVersion("/example-password")
			url := dummyAuth.Request.URL.String()
			Expect(url).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-password"))
			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a credential by name", func() {
				responseString := `{
	"data": [
	{
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-password",
      "metadata": {"some":"metadata"},
      "version_created_at": "2017-01-05T01:01:01Z"
    }
    ]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

				cred, err := ch.GetLatestVersion("/example-password")
				Expect(err).To(BeNil())
				Expect(cred.Id).To(Equal("some-id"))
				Expect(cred.Name).To(Equal("/example-password"))
				Expect(cred.Type).To(Equal("password"))
				Expect(cred.Value.(string)).To(Equal("some-password"))
				Expect(cred.Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})
		})

		Context("when the response body contains an empty list", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":[]}`)),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetLatestVersion("/example-password")

				Expect(err).To(MatchError("response did not contain any credentials"))
			})
		})
	})

	Describe("GetById()", func() {
		It("requests the credential by id", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

			ch.GetById("0239482304958")
			url := dummyAuth.Request.URL.String()
			Expect(url).To(Equal("https://example.com/api/v1/data/0239482304958"))
			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a credential by name", func() {
				responseString := `{
      "id": "0239482304958",
      "name": "/reasonable-password",
      "type": "password",
      "value": "some-password",
	  "metadata": {"some":"metadata"},
      "version_created_at": "2017-01-05T01:01:01Z"
    }`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

				cred, err := ch.GetById("0239482304958")
				Expect(err).To(BeNil())
				Expect(cred.Id).To(Equal("0239482304958"))
				Expect(cred.Name).To(Equal("/reasonable-password"))
				Expect(cred.Type).To(Equal("password"))
				Expect(cred.Value.(string)).To(Equal("some-password"))
				Expect(cred.Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})
		})
	})

	Describe("GetAllVersions()", func() {
		It("makes a request for all versions of a credential", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

			ch.GetAllVersions("/example-password")
			url := dummyAuth.Request.URL.String()
			Expect(url).To(Equal("https://example.com/api/v1/data?name=%2Fexample-password"))
			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a list of all passwords", func() {
				responseString := `{
	"data": [
	{
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-password",
	  "metadata": {"some":"metadata"},
      "version_created_at": "2017-01-05T01:01:01Z"
    },
	{
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-other-password",
	  "metadata": null,
      "version_created_at": "2017-01-05T01:01:01Z"
    }
    ]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

				creds, err := ch.GetAllVersions("/example-password")
				Expect(err).To(BeNil())
				Expect(creds[0].Id).To(Equal("some-id"))
				Expect(creds[0].Name).To(Equal("/example-password"))
				Expect(creds[0].Type).To(Equal("password"))
				Expect(creds[0].Value.(string)).To(Equal("some-password"))
				Expect(creds[0].Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(creds[0].VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))

				Expect(creds[1].Id).To(Equal("some-id"))
				Expect(creds[1].Name).To(Equal("/example-password"))
				Expect(creds[1].Type).To(Equal("password"))
				Expect(creds[1].Value.(string)).To(Equal("some-other-password"))
				Expect(creds[1].Metadata).To(Equal(nilMetadata))
				Expect(creds[1].VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})

			It("returns a list of all users", func() {
				responseString := `{
	"data": [
	{
      "id": "some-id",
      "name": "/example-user",
      "type": "user",
      "value": {
      	"username": "first-username",
      	"password": "dummy_password",
      	"password_hash": "$6$kjhlkjh$lkjhasdflkjhasdflkjh"
      },
      "metadata": null,
      "version_created_at": "2017-01-05T01:01:01Z"
    },
	{
      "id": "some-id",
      "name": "/example-user",
      "type": "user",
      "value": {
      	"username": "second-username",
      	"password": "another_random_dummy_password",
      	"password_hash": "$6$kjhlkjh$lkjhasdflkjhasdflkjh"
      },
      "metadata": {"some":"metadata"},
      "version_created_at": "2017-01-05T01:01:01Z"
    }
    ]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

				creds, err := ch.GetAllVersions("/example-user")
				Expect(err).To(BeNil())
				Expect(creds[0].Id).To(Equal("some-id"))
				Expect(creds[0].Name).To(Equal("/example-user"))
				Expect(creds[0].Type).To(Equal("user"))
				firstCredValue := creds[0].Value.(map[string]interface{})
				Expect(firstCredValue["username"]).To(Equal("first-username"))
				Expect(firstCredValue["password"]).To(Equal("dummy_password"))
				Expect(firstCredValue["password_hash"]).To(Equal("$6$kjhlkjh$lkjhasdflkjhasdflkjh"))
				Expect(creds[0].Metadata).To(Equal(nilMetadata))
				Expect(creds[0].VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))

				Expect(creds[1].Id).To(Equal("some-id"))
				Expect(creds[1].Name).To(Equal("/example-user"))
				Expect(creds[1].Type).To(Equal("user"))
				secondCredValue := creds[1].Value.(map[string]interface{})
				Expect(secondCredValue["username"]).To(Equal("second-username"))
				Expect(secondCredValue["password"]).To(Equal("another_random_dummy_password"))
				Expect(secondCredValue["password_hash"]).To(Equal("$6$kjhlkjh$lkjhasdflkjhasdflkjh"))
				Expect(creds[1].Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(creds[1].VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})
		})

		Context("when the response body contains an empty list", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":[]}`)),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetAllVersions("/example-password")

				Expect(err).To(MatchError("response did not contain any credentials"))
			})
		})
	})

	Describe("GetNVersions()", func() {
		It("makes a request for N versions of a credential", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

			ch.GetNVersions("/example-password", 3)
			url := dummyAuth.Request.URL.String()
			Expect(url).To(Equal("https://example.com/api/v1/data?name=%2Fexample-password&versions=3"))
			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a list of N passwords", func() {
				responseString := `{
	"data": [
	{
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-password",
	  "metadata": {"some":"metadata"},
      "version_created_at": "2017-01-05T01:01:01Z"
    },
	{
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-other-password",
	  "metadata": null,
      "version_created_at": "2017-01-05T01:01:01Z"
    }
    ]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

				creds, err := ch.GetNVersions("/example-password", 2)
				Expect(err).To(BeNil())
				Expect(creds[0].Id).To(Equal("some-id"))
				Expect(creds[0].Name).To(Equal("/example-password"))
				Expect(creds[0].Type).To(Equal("password"))
				Expect(creds[0].Value.(string)).To(Equal("some-password"))
				Expect(creds[0].Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(creds[0].VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))

				Expect(creds[1].Id).To(Equal("some-id"))
				Expect(creds[1].Name).To(Equal("/example-password"))
				Expect(creds[1].Type).To(Equal("password"))
				Expect(creds[1].Value.(string)).To(Equal("some-other-password"))
				Expect(creds[1].Metadata).To(Equal(nilMetadata))
				Expect(creds[1].VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})

			It("returns a list of N users", func() {
				responseString := `{
	"data": [
	{
      "id": "some-id",
      "name": "/example-user",
      "type": "user",
      "value": {
      	"username": "first-username",
      	"password": "dummy_password",
      	"password_hash": "$6$kjhlkjh$lkjhasdflkjhasdflkjh"
      },
	  "metadata": null,
      "version_created_at": "2017-01-05T01:01:01Z"
    },
	{
      "id": "some-id",
      "name": "/example-user",
      "type": "user",
      "value": {
      	"username": "second-username",
      	"password": "another_random_dummy_password",
      	"password_hash": "$6$kjhlkjh$lkjhasdflkjhasdflkjh"
      },
	  "metadata": {"some":"metadata"},
      "version_created_at": "2017-01-05T01:01:01Z"
    }
    ]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

				creds, err := ch.GetNVersions("/example-user", 2)
				Expect(err).To(BeNil())
				Expect(creds[0].Id).To(Equal("some-id"))
				Expect(creds[0].Name).To(Equal("/example-user"))
				Expect(creds[0].Type).To(Equal("user"))
				firstCredValue := creds[0].Value.(map[string]interface{})
				Expect(firstCredValue["username"]).To(Equal("first-username"))
				Expect(firstCredValue["password"]).To(Equal("dummy_password"))
				Expect(firstCredValue["password_hash"]).To(Equal("$6$kjhlkjh$lkjhasdflkjhasdflkjh"))
				Expect(creds[0].Metadata).To(Equal(nilMetadata))
				Expect(creds[0].VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))

				Expect(creds[1].Id).To(Equal("some-id"))
				Expect(creds[1].Name).To(Equal("/example-user"))
				Expect(creds[1].Type).To(Equal("user"))
				secondCredValue := creds[1].Value.(map[string]interface{})
				Expect(secondCredValue["username"]).To(Equal("second-username"))
				Expect(secondCredValue["password"]).To(Equal("another_random_dummy_password"))
				Expect(secondCredValue["password_hash"]).To(Equal("$6$kjhlkjh$lkjhasdflkjhasdflkjh"))
				Expect(creds[1].Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(creds[1].VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})
		})

		Context("when the response body contains an empty list", func() {
			It("returns an error", func() {
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":[]}`)),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetNVersions("/example-password", 2)

				Expect(err).To(MatchError("response did not contain any credentials"))
			})
		})

		Context("when the server returns a 200 but the cli can not parse the response", func() {
			It("returns an appropriate error", func() {
				responseString := `some-invalid-json`

				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}
				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				_, err := ch.GetNVersions("/example-password", 2)

				Expect(err.Error()).To(ContainSubstring("The response body could not be decoded:"))

			})

		})
	})

	Describe("GetLatestPassword()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetLatestPassword("/example-password")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-password"))
			Expect(url.String()).To(ContainSubstring("current=true"))
			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a password credential", func() {
				responseString := `{
  "data": [
    {
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-password",
	  "metadata": {"some":"metadata"},
      "version_created_at": "2017-01-05T01:01:01Z"
    }]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetLatestPassword("/example-password")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Id).To(Equal("some-id"))
				Expect(cred.Name).To(Equal("/example-password"))
				Expect(cred.Type).To(Equal("password"))
				Expect(cred.Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(cred.Value).To(BeEquivalentTo("some-password"))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})
		})
	})

	Describe("GetLatestCertificate()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetLatestCertificate("/example-certificate")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-certificate"))
			Expect(url.String()).To(ContainSubstring("current=true"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a certificate credential", func() {
				responseString := `{
				  "data": [{
	"id": "some-id",
	"name": "/example-certificate",
	"type": "certificate",
	"value": {
		"ca": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"certificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
	},
    "metadata": {"some":"metadata"},
	"version_created_at": "2017-01-01T04:07:18Z"
}]}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))

				cred, err := ch.GetLatestCertificate("/example-certificate")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Id).To(Equal("some-id"))
				Expect(cred.Name).To(Equal("/example-certificate"))
				Expect(cred.Type).To(Equal("certificate"))
				Expect(cred.Value.Ca).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
				Expect(cred.Value.Certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
				Expect(cred.Value.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))
				Expect(cred.Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))
			})
		})
	})

	Describe("GetLatestUser()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetLatestUser("/example-user")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-user"))
			Expect(url.String()).To(ContainSubstring("current=true"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a user credential", func() {
				responseString := `{
				  "data": [
					{
					  "id": "some-id",
					  "name": "/example-user",
					  "type": "user",
					  "value": {
						"username": "some-username",
						"password": "some-password",
						"password_hash": "some-hash"
					  },
					  "metadata": {"some":"metadata"},
       				  "version_created_at": "2017-01-05T01:01:01Z"
					}
				  ]
				}`
				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetLatestUser("/example-user")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Id).To(Equal("some-id"))
				Expect(cred.Name).To(Equal("/example-user"))
				Expect(cred.Type).To(Equal("user"))
				Expect(cred.Value.PasswordHash).To(Equal("some-hash"))
				Expect(cred.Value.User).To(Equal(values.User{
					Username: "some-username",
					Password: "some-password",
				}))
				Expect(cred.Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})
		})
	})

	Describe("GetLatestRSA()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetLatestRSA("/example-rsa")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-rsa"))
			Expect(url.String()).To(ContainSubstring("current=true"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a rsa credential", func() {
				responseString := `{
				  "data": [
					{
					  "id": "some-id",
					  "name": "/example-rsa",
					  "type": "rsa",
					  "value": {
						"public_key": "public-key",
						"private_key": "private-key"
					  },
					  "metadata": {"some":"metadata"},
					  "version_created_at": "2017-01-01T04:07:18Z"
					}
				  ]
				}`

				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetLatestRSA("/example-rsa")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Id).To(Equal("some-id"))
				Expect(cred.Name).To(Equal("/example-rsa"))
				Expect(cred.Type).To(Equal("rsa"))
				Expect(cred.Value).To(Equal(values.RSA{
					PublicKey:  "public-key",
					PrivateKey: "private-key",
				}))
				Expect(cred.Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))
			})
		})
	})

	Describe("GetLatestSSH()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetLatestSSH("/example-ssh")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-ssh"))
			Expect(url.String()).To(ContainSubstring("current=true"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a ssh credential", func() {
				responseString := `{
				  "data": [
					{
					  "id": "some-id",
					  "name": "/example-ssh",
					  "type": "ssh",
					  "value": {
						"public_key": "public-key",
						"private_key": "private-key",
						"public_key_fingerprint": "public-key-fingerprint"
					  },
					  "metadata": {"some":"metadata"},
					  "version_created_at": "2017-01-01T04:07:18Z"
					}
				  ]
				}`

				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetLatestSSH("/example-ssh")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Id).To(Equal("some-id"))
				Expect(cred.Name).To(Equal("/example-ssh"))
				Expect(cred.Type).To(Equal("ssh"))
				Expect(cred.Value.PublicKeyFingerprint).To(Equal("public-key-fingerprint"))
				Expect(cred.Value.SSH).To(Equal(values.SSH{
					PublicKey:  "public-key",
					PrivateKey: "private-key",
				}))
				Expect(cred.Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))
			})
		})
	})

	Describe("GetLatestJSON()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetLatestJSON("/example-json")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-json"))
			Expect(url.String()).To(ContainSubstring("current=true"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a json credential", func() {
				responseString := `{
				  "data": [
					{
					  "id": "some-id",
					  "name": "/example-json",
					  "type": "json",
					  "value": {
						"key": 123,
						"key_list": [
						  "val1",
						  "val2"
						],
						"is_true": true
					  },
					  "metadata": {"some":"metadata"},
					  "version_created_at": "2017-01-01T04:07:18Z"
					}
				  ]
				}`

				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetLatestJSON("/example-json")

				JSONResult := `{
						"key": 123,
						"key_list": [
						  "val1",
						  "val2"
						],
						"is_true": true
					}`
				var unmarshalled values.JSON
				json.Unmarshal([]byte(JSONResult), &unmarshalled)
				Expect(err).ToNot(HaveOccurred())

				Expect(cred.Id).To(Equal("some-id"))
				Expect(cred.Name).To(Equal("/example-json"))
				Expect(cred.Type).To(Equal("json"))
				Expect(cred.Value).To(Equal(unmarshalled))
				Expect(cred.Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))
			})
		})

	})

	Describe("GetLatestValue()", func() {
		It("requests the credential by name", func() {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			ch.GetLatestValue("/example-value")
			url := dummyAuth.Request.URL
			Expect(url.String()).To(ContainSubstring("https://example.com/api/v1/data"))
			Expect(url.String()).To(ContainSubstring("name=%2Fexample-value"))
			Expect(url.String()).To(ContainSubstring("current=true"))

			Expect(dummyAuth.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a value credential", func() {
				responseString := `{
				  "data": [
					{
					  "id": "some-id",
					  "name": "/example-value",
					  "type": "value",
					  "value": "some-value",
					  "metadata": {"some":"metadata"},
					  "version_created_at": "2017-01-05T01:01:01Z"
				}]}`

				dummyAuth := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
				cred, err := ch.GetLatestValue("/example-value")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Id).To(Equal("some-id"))
				Expect(cred.Name).To(Equal("/example-value"))
				Expect(cred.Type).To(Equal("value"))
				Expect(cred.Value).To(Equal(values.Value("some-value")))
				Expect(cred.Metadata).To(Equal(credentials.Metadata{"some": "metadata"}))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})
		})
	})

	var listOfActions = []TableEntry{
		Entry("GetAllVersions", func(ch *CredHub) error {
			_, err := ch.GetAllVersions("/example-password")
			return err
		}),
		Entry("GetNVersions", func(ch *CredHub) error {
			_, err := ch.GetNVersions("/example-password", 47)
			return err
		}),
		Entry("GetLatestVersion", func(ch *CredHub) error {
			_, err := ch.GetLatestVersion("/example-password")
			return err
		}),
		Entry("GetLatestPassword", func(ch *CredHub) error {
			_, err := ch.GetLatestPassword("/example-password")
			return err
		}),
		Entry("GetLatestCertificate", func(ch *CredHub) error {
			_, err := ch.GetLatestCertificate("/example-certificate")
			return err
		}),
		Entry("GetLatestUser", func(ch *CredHub) error {
			_, err := ch.GetLatestUser("/example-password")
			return err
		}),
		Entry("GetLatestRSA", func(ch *CredHub) error {
			_, err := ch.GetLatestRSA("/example-password")
			return err
		}),
		Entry("GetById", func(ch *CredHub) error {
			_, err := ch.GetById("0239482304958")
			return err
		}),
		Entry("GetLatestSSH", func(ch *CredHub) error {
			_, err := ch.GetLatestSSH("/example-password")
			return err
		}),
		Entry("GetLatestJSON", func(ch *CredHub) error {
			_, err := ch.GetLatestJSON("/example-password")
			return err
		}),
		Entry("GetLatestValue", func(ch *CredHub) error {
			_, err := ch.GetLatestValue("/example-password")
			return err
		}),
	}

	DescribeTable("errors when response body is not json and/or cannot be decoded",
		func(performAction func(*CredHub) error) {
			dummyAuth := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("<html>")),
			}}
			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			err := performAction(ch)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("The response body could not be decoded:"))
		},

		listOfActions,
	)

	DescribeTable("returns credhub error when the cred does not exist",
		func(performAction func(*CredHub) error) {
			dummyAuth := &DummyAuth{Response: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"The request could not be completed because the credential does not exist or you do not have sufficient authorization."}`)),
			}}
			ch, _ := New("https://example.com", Auth(dummyAuth.Builder()))
			err := performAction(ch)

			Expect(err).To(MatchError("The request could not be completed because the credential does not exist or you do not have sufficient authorization."))
		},

		listOfActions,
	)

	DescribeTable("request fails due to network error",
		func(performAction func(*CredHub) error) {
			networkError := errors.New("Network error occurred making an http request")
			dummyAuth := &DummyAuth{Error: networkError}
			ch, err := New("https://example.com", Auth(dummyAuth.Builder()))
			Expect(err).NotTo(HaveOccurred())

			err = performAction(ch)

			Expect(err).To(Equal(networkError))
		},

		listOfActions,
	)
})
