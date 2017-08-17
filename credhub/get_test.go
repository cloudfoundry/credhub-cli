package credhub_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("Get", func() {

	Describe("Get()", func() {
		It("requests the credential by name", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))

			ch.Get("/example-password")
			url := dummy.Request.URL.String()
			Expect(url).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-password"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
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
      "version_created_at": "2017-01-05T01:01:01Z"
    }
    ]}`
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))

				cred, err := ch.Get("/example-password")
				Expect(err).To(BeNil())
				Expect(cred.Id).To(Equal("some-id"))
				Expect(cred.Name).To(Equal("/example-password"))
				Expect(cred.Type).To(Equal("password"))
				Expect(cred.Value.(string)).To(Equal("some-password"))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}

				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.Get("/example-password")

				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the response body contains an empty list", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":[]}`)),
				}}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.Get("/example-password")

				Expect(err).To(MatchError("response did not contain any credentials"))
			})
		})
	})

	Describe("GetPassword()", func() {
		It("requests the credential by name", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))
			ch.GetPassword("/example-password")
			url := dummy.Request.URL
			Expect(url.String()).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-password"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
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
      "version_created_at": "2017-01-05T01:01:01Z"
    }]}`
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))
				cred, err := ch.GetPassword("/example-password")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value).To(BeEquivalentTo("some-password"))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.GetPassword("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetCertificate()", func() {
		It("requests the credential by name", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))
			ch.GetCertificate("/example-certificate")
			url := dummy.Request.URL
			Expect(url.String()).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-certificate"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
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
	"version_created_at": "2017-01-01T04:07:18Z"
}]}`
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))

				cred, err := ch.GetCertificate("/example-certificate")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value.Ca).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
				Expect(cred.Value.Certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
				Expect(cred.Value.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))
				Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.GetCertificate("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetUser()", func() {
		It("requests the credential by name", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))
			ch.GetUser("/example-user")
			url := dummy.Request.URL
			Expect(url.String()).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-user"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
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
					  "version_created_at": "2017-01-05T01:01:01Z"
					}
				  ]
				}`
				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))
				cred, err := ch.GetUser("/example-user")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value.PasswordHash).To(Equal("some-hash"))
				Expect(cred.Value.User).To(Equal(values.User{
					Username: "some-username",
					Password: "some-password",
				}))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.GetUser("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetRSA()", func() {
		It("requests the credential by name", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))
			ch.GetRSA("/example-rsa")
			url := dummy.Request.URL
			Expect(url.String()).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-rsa"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
		})

		Context("when successful", func() {
			It("returns a rsa credential", func() {
				responseString := `{
				  "data": [
					{
					  "id": "67fc3def-bbfb-4953-83f8-4ab0682ad677",
					  "name": "/example-rsa",
					  "type": "rsa",
					  "value": {
						"public_key": "public-key",
						"private_key": "private-key"
					  },
					  "version_created_at": "2017-01-01T04:07:18Z"
					}
				  ]
				}`

				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))
				cred, err := ch.GetRSA("/example-rsa")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value).To(Equal(values.RSA{
					PublicKey:  "public-key",
					PrivateKey: "private-key",
				}))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.GetRSA("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetSSH()", func() {
		It("requests the credential by name", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))
			ch.GetSSH("/example-ssh")
			url := dummy.Request.URL
			Expect(url.String()).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-ssh"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
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
					  "version_created_at": "2017-01-01T04:07:18Z"
					}
				  ]
				}`

				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))
				cred, err := ch.GetSSH("/example-ssh")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value.PublicKeyFingerprint).To(Equal("public-key-fingerprint"))
				Expect(cred.Value.SSH).To(Equal(values.SSH{
					PublicKey:  "public-key",
					PrivateKey: "private-key",
				}))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.GetSSH("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetJSON()", func() {
		It("requests the credential by name", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))
			ch.GetJSON("/example-json")
			url := dummy.Request.URL
			Expect(url.String()).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-json"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
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
					  "version_created_at": "2017-01-01T04:07:18Z"
					}
				  ]
				}`

				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))
				cred, err := ch.GetJSON("/example-json")
				Expect(err).ToNot(HaveOccurred())
				Expect([]byte(cred.Value)).To(MatchJSON(`{
						"key": 123,
						"key_list": [
						  "val1",
						  "val2"
						],
						"is_true": true
					}`))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.GetJSON("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetValue()", func() {
		It("requests the credential by name", func() {
			dummy := &DummyAuth{Response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}}

			ch, _ := New("https://example.com", Auth(dummy))
			ch.GetValue("/example-value")
			url := dummy.Request.URL
			Expect(url.String()).To(Equal("https://example.com/api/v1/data?current=true&name=%2Fexample-value"))
			Expect(dummy.Request.Method).To(Equal(http.MethodGet))
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
					  "version_created_at": "2017-01-05T01:01:01Z"
				}]}`

				dummy := &DummyAuth{Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(responseString)),
				}}

				ch, _ := New("https://example.com", Auth(dummy))
				cred, err := ch.GetValue("/example-value")
				Expect(err).ToNot(HaveOccurred())
				Expect(cred.Value).To(Equal(values.Value("some-value")))
			})
		})

		Context("when response body cannot be unmarshalled", func() {
			It("returns an error", func() {
				dummy := &DummyAuth{Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("something-invalid")),
				}}
				ch, _ := New("https://example.com", Auth(dummy))
				_, err := ch.GetValue("/example-cred")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	DescribeTable("request fails due to network error",
		func(performAction func(*CredHub) error) {
			networkError := errors.New("Network error occurred")
			dummy := &DummyAuth{Error: networkError}
			ch, _ := New("https://example.com", Auth(dummy))

			err := performAction(ch)

			Expect(err).To(Equal(networkError))
		},

		Entry("Get", func(ch *CredHub) error {
			_, err := ch.Get("/example-password")
			return err
		}),
		Entry("GetPassword", func(ch *CredHub) error {
			_, err := ch.GetPassword("/example-password")
			return err
		}),
		Entry("GetCertificate", func(ch *CredHub) error {
			_, err := ch.GetCertificate("/example-certificate")
			return err
		}),
		Entry("GetUser", func(ch *CredHub) error {
			_, err := ch.GetUser("/example-password")
			return err
		}),
		Entry("GetRSA", func(ch *CredHub) error {
			_, err := ch.GetRSA("/example-password")
			return err
		}),
		Entry("GetSSH", func(ch *CredHub) error {
			_, err := ch.GetSSH("/example-password")
			return err
		}),
		Entry("GetJSON", func(ch *CredHub) error {
			_, err := ch.GetJSON("/example-password")
			return err
		}),
		Entry("GetValue", func(ch *CredHub) error {
			_, err := ch.GetValue("/example-password")
			return err
		}),
	)
})
