package credentials_test

import (
	"encoding/json"

	"gopkg.in/yaml.v2"

	. "code.cloudfoundry.org/credhub-cli/credhub/credentials"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	Describe("Certificate", func() {
		Specify("when decoding and encoding with duration_overridden and duration_used in the output", func() {
			var cred Certificate

			credJSON := `{
	"id": "some-id",
	"name": "/example-certificate",
	"type": "certificate",
	"duration_overridden": true,
	"duration_used": 1234,
	"value": {
		"ca": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"certificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
	},
	"metadata": {"some":"metadata"},
	"version_created_at": "2017-01-01T04:07:18Z"
}`
			credYaml := `id: some-id
name: /example-certificate
type: certificate
duration_overridden: true
duration_used: 1234
value:
  ca: |-
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  certificate: |-
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  private_key: |-
    -----BEGIN RSA PRIVATE KEY-----
    ...
    -----END RSA PRIVATE KEY-----
metadata:
  some: metadata
version_created_at: '2017-01-01T04:07:18Z'`

			err := json.Unmarshal([]byte(credJSON), &cred)

			Expect(err).To(BeNil())

			Expect(cred.Id).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-certificate"))
			Expect(cred.Type).To(Equal("certificate"))
			Expect(cred.DurationOverridden).To(Equal(true))
			Expect(cred.DurationUsed).To(Equal(1234))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.Value.Ca).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
			Expect(cred.Value.Certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
			Expect(cred.Value.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))

			yamlOutput, err := yaml.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(MatchYAML(credYaml))
		})

		Specify("when decoding and encoding with NO duration_overridden in the output", func() {
			var cred Certificate

			credJSON := `{
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
}`
			credYaml := `id: some-id
name: /example-certificate
type: certificate
value:
  ca: |-
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  certificate: |-
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  private_key: |-
    -----BEGIN RSA PRIVATE KEY-----
    ...
    -----END RSA PRIVATE KEY-----
metadata:
  some: metadata
version_created_at: '2017-01-01T04:07:18Z'`

			err := json.Unmarshal([]byte(credJSON), &cred)

			Expect(err).To(BeNil())

			Expect(cred.Id).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-certificate"))
			Expect(cred.Type).To(Equal("certificate"))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.Value.Ca).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
			Expect(cred.Value.Certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
			Expect(cred.Value.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))

			yamlOutput, err := yaml.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(MatchYAML(credYaml))
		})
	})

	Describe("User", func() {
		Specify("when decoding and encoding", func() {
			var cred User
			credJSON := `{
      "id": "some-id",
      "name": "/example-user",
      "type": "user",
      "value": {
        "username": "some-username",
        "password": "some-password",
        "password_hash": "some-password-hash"
      },
      "metadata": {"some":"metadata"},
      "version_created_at": "2017-01-05T01:01:01Z"
}`

			credYaml := `
id: some-id
name: "/example-user"
type: user
value:
  username: some-username
  password: some-password
  password_hash: some-password-hash
metadata:
  some: metadata
version_created_at: '2017-01-05T01:01:01Z'`

			err := json.Unmarshal([]byte(credJSON), &cred)

			Expect(err).To(BeNil())

			Expect(cred.Id).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-user"))
			Expect(cred.Type).To(Equal("user"))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.Value.Username).To(Equal("some-username"))
			Expect(cred.Value.Password).To(Equal("some-password"))
			Expect(cred.Value.PasswordHash).To(Equal("some-password-hash"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))

			yamlOutput, err := yaml.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(MatchYAML(credYaml))
		})

	})

	Describe("Password", func() {
		Specify("when decoding and encoding", func() {
			var cred Password

			credJSON := ` {
      "id": "some-id",
      "name": "/example-password",
      "type": "password",
      "value": "some-password",
	  "metadata": {"some":"metadata"},
      "version_created_at": "2017-01-05T01:01:01Z"
    }`

			credYaml := `
id: some-id
name: "/example-password"
type: password
value: some-password
metadata:
  some: metadata
version_created_at: '2017-01-05T01:01:01Z'
`

			err := json.Unmarshal([]byte(credJSON), &cred)
			Expect(err).To(BeNil())

			Expect(cred.Id).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-password"))
			Expect(cred.Type).To(Equal("password"))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.Value).To(BeEquivalentTo("some-password"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))

			yamlOutput, err := yaml.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(MatchYAML(credYaml))
		})
	})

	Describe("JSON", func() {
		Specify("when decoding and encoding", func() {
			var cred JSON

			credJSON := ` {
      "id": "some-id",
      "name": "/example-json",
      "type": "json",
	  "metadata": {"some":"metadata"},
      "value": {
        "key": 123,
        "key_list": [
          "val1",
          "val2"
        ],
        "is_true": true
      },
      "version_created_at": "2017-01-01T04:07:18Z"
    }`

			credYaml := `
id: some-id
name: "/example-json"
type: json
value:
  key: 123
  key_list:
  - val1
  - val2
  is_true: true
metadata:
  some: metadata
version_created_at: '2017-01-01T04:07:18Z'
`

			err := json.Unmarshal([]byte(credJSON), &cred)

			Expect(err).To(BeNil())

			jsonValueString := `{
        "key": 123,
        "key_list": [
          "val1",
          "val2"
        ],
        "is_true": true
      }`

			var unmarshalled values.JSON
			Expect(json.Unmarshal([]byte(jsonValueString), &unmarshalled)).To(Succeed())

			Expect(cred.Id).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-json"))
			Expect(cred.Type).To(Equal("json"))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.Value).To(Equal(unmarshalled))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))

			yamlOutput, err := yaml.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(MatchYAML(credYaml))
		})
	})

	Describe("Value", func() {
		Specify("when decoding and encoding", func() {
			var cred Value

			credJSON := ` {
      "id": "some-id",
      "name": "/example-value",
      "type": "value",
	  "metadata": {"some":"metadata"},
      "value": "some-value",
      "version_created_at": "2017-01-05T01:01:01Z"
    }`

			credYaml := `
id: some-id
name: "/example-value"
type: value
value: some-value
metadata:
  some: metadata
version_created_at: '2017-01-05T01:01:01Z'
`

			Expect(json.Unmarshal([]byte(credJSON), &cred)).To(Succeed())

			Expect(cred.Id).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-value"))
			Expect(cred.Type).To(Equal("value"))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.Value).To(BeEquivalentTo("some-value"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))

			yamlOutput, err := yaml.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(MatchYAML(credYaml))
		})
	})

	Describe("RSA", func() {
		Specify("when decoding and encoding", func() {
			var cred RSA
			credJSON := `{
      "id": "some-id",
      "name": "/example-rsa",
      "type": "rsa",
      "value": {
        "public_key": "some-public-key",
        "private_key": "some-private-key"
      },
	  "metadata": {"some":"metadata"},
      "version_created_at": "2017-01-05T01:01:01Z"
}`

			credYaml := `
id: some-id
name: "/example-rsa"
type: rsa
value:
  public_key: some-public-key
  private_key: some-private-key
metadata:
  some: metadata
version_created_at: '2017-01-05T01:01:01Z'`

			Expect(json.Unmarshal([]byte(credJSON), &cred)).To(Succeed())

			Expect(cred.Id).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-rsa"))
			Expect(cred.Type).To(Equal("rsa"))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.Value.PublicKey).To(Equal("some-public-key"))
			Expect(cred.Value.PrivateKey).To(Equal("some-private-key"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-05T01:01:01Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))

			yamlOutput, err := yaml.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(MatchYAML(credYaml))
		})
	})

	Describe("SSH", func() {
		Specify("when decoding and encoding", func() {
			var cred SSH
			credJSON := `{
      "id": "some-id",
      "name": "/example-ssh",
      "type": "ssh",
      "value": {
        "public_key": "some-public-key",
        "private_key": "some-private-key",
        "public_key_fingerprint": "some-public-key-fingerprint"
      },
	  "metadata": {"some":"metadata"},
      "version_created_at": "2017-01-01T04:07:18Z"
    }`

			credYaml := `
id: some-id
name: "/example-ssh"
type: ssh
value:
  public_key: some-public-key
  private_key: some-private-key
  public_key_fingerprint: some-public-key-fingerprint
metadata:
  some: metadata
version_created_at: '2017-01-01T04:07:18Z'`

			Expect(json.Unmarshal([]byte(credJSON), &cred)).To(Succeed())

			Expect(cred.Id).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-ssh"))
			Expect(cred.Type).To(Equal("ssh"))
			Expect(cred.Metadata).To(Equal(Metadata{"some": "metadata"}))
			Expect(cred.Value.PublicKey).To(Equal("some-public-key"))
			Expect(cred.Value.PublicKeyFingerprint).To(Equal("some-public-key-fingerprint"))
			Expect(cred.Value.PrivateKey).To(Equal("some-private-key"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))

			yamlOutput, err := yaml.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(MatchYAML(credYaml))
		})
	})

	Describe("Credential", func() {
		Specify("when decoding and encoding", func() {
			var cred Credential
			credJSON := `{
      "id": "some-id",
      "name": "/example-ssh",
      "type": "ssh",
      "value": {
        "public_key": "some-public-key",
        "private_key": "some-private-key",
        "public_key_fingerprint": "some-public-key-fingerprint"
      },
	  "metadata": {"some":{"example":"metadata"}},
      "version_created_at": "2017-01-01T04:07:18Z"
    }`

			credYaml := `id: some-id
name: /example-ssh
type: ssh
value:
  private_key: some-private-key
  public_key: some-public-key
  public_key_fingerprint: some-public-key-fingerprint
metadata:
  some:
    example: metadata
version_created_at: "2017-01-01T04:07:18Z"
`

			Expect(json.Unmarshal([]byte(credJSON), &cred)).To(Succeed())

			jsonValueString := `{
        "public_key": "some-public-key",
        "private_key": "some-private-key",
        "public_key_fingerprint": "some-public-key-fingerprint"
      }`
			var jsonValue map[string]interface{}
			Expect(json.Unmarshal([]byte(jsonValueString), &jsonValue)).To(Succeed())

			Expect(cred.Id).To(Equal("some-id"))
			Expect(cred.Name).To(Equal("/example-ssh"))
			Expect(cred.Type).To(Equal("ssh"))
			Expect(cred.Value).To(Equal(jsonValue))
			Expect(cred.Metadata).To(Equal(Metadata{"some": map[string]interface{}{"example": "metadata"}}))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(credJSON))

			yamlOutput, err := yaml.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(yamlOutput)).To(Equal(credYaml)) // Order matters for yaml output for "default" output
		})

		It("does not Marshal the metadata key when metadata is empty or nil for yaml", func() {
			cred := Credential{
				Value: "some-value",
				Base: Base{
					Id:               "some-id",
					Name:             "some-name",
					Type:             "some-type",
					VersionCreatedAt: "some-time",
					Metadata:         Metadata{},
				},
			}

			nilMetadatacred := Credential{
				Value: "some-value",
				Base: Base{
					Id:               "some-id",
					Name:             "some-name",
					Type:             "some-type",
					VersionCreatedAt: "some-time",
					Metadata:         nil,
				},
			}

			expectedYAMLOutput := `id: some-id
name: some-name
type: some-type
value: some-value
version_created_at: some-time
`

			yamlOutput, err := yaml.Marshal(cred)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(yamlOutput)).To(Equal(expectedYAMLOutput)) // Order matters for yaml output for "default" output
			yamlOutput, err = yaml.Marshal(nilMetadatacred)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(yamlOutput)).To(Equal(expectedYAMLOutput)) // Order matters for yaml output for "default" output
		})
	})

	Describe("Certificate Metadata", func() {
		Specify("when decoding and encoding", func() {
			metadataJSON := `{
      "id": "some-id",
      "name": "/some-cert",
      "signed_by": "/some-cert",
      "signs": ["/another-cert"],
      "versions": [
        {
          "expiry_date": "2020-05-29T12:33:50Z",
          "id": "some-other-id",
          "transitional": false,
		  "certificate_authority": true,
		  "self_signed": true
        }
      ]
    }`

			metadataYaml := `
id: some-id
name: "/some-cert"
signed_by: "/some-cert"
signs:
- "/another-cert"
versions:
- expiry_date: '2020-05-29T12:33:50Z'
  id: some-other-id
  transitional: false
  certificate_authority: true
  self_signed: true
`
			var certMetadata CertificateMetadata
			Expect(json.Unmarshal([]byte(metadataJSON), &certMetadata)).To(Succeed())

			jsonVersionString := `{
          "expiry_date": "2020-05-29T12:33:50Z",
          "id": "some-other-id",
          "transitional": false,
          "certificate_authority": true,
		  "self_signed": true
        }`
			var jsonVersion CertificateMetadataVersion
			Expect(json.Unmarshal([]byte(jsonVersionString), &jsonVersion)).To(Succeed())

			Expect(certMetadata.Id).To(Equal("some-id"))
			Expect(certMetadata.Name).To(Equal("/some-cert"))
			Expect(certMetadata.Signs[0]).To(Equal("/another-cert"))
			Expect(certMetadata.SignedBy).To(Equal("/some-cert"))
			Expect(certMetadata.Versions[0]).To(Equal(jsonVersion))

			jsonOutput, err := json.Marshal(certMetadata)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonOutput).To(MatchJSON(metadataJSON))

			yamlOutput, err := yaml.Marshal(certMetadata)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(MatchYAML(metadataYaml))
		})
	})
})
