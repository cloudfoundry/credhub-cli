package credentials_test

import (
	"encoding/json"

	"gopkg.in/yaml.v2"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	Context("Certificate", func() {
		Specify("when decoding and encoding", func() {
			var cred Certificate

			credJson := `{
	"id": "67fc3def-bbfb-4953-83f8-4ab0682ad676",
	"name": "/example-certificate",
	"type": "certificate",
	"value": {
		"ca": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"certificate": "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
	},
	"version_created_at": "2017-01-01T04:07:18Z"
}`
			credYaml := `id: 67fc3def-bbfb-4953-83f8-4ab0682ad676
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
version_created_at: 2017-01-01T04:07:18Z`

			err := json.Unmarshal([]byte(credJson), &cred)

			Expect(err).To(BeNil())

			Expect(cred.Id).To(Equal("67fc3def-bbfb-4953-83f8-4ab0682ad676"))
			Expect(cred.Name).To(Equal("/example-certificate"))
			Expect(cred.Type).To(Equal("certificate"))
			Expect(cred.Value.Ca).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
			Expect(cred.Value.Certificate).To(Equal("-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"))
			Expect(cred.Value.PrivateKey).To(Equal("-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"))
			Expect(cred.VersionCreatedAt).To(Equal("2017-01-01T04:07:18Z"))

			jsonOutput, err := json.Marshal(cred)

			Expect(jsonOutput).To(MatchJSON(credJson))

			yamlOutput, err := yaml.Marshal(cred)

			Expect(yamlOutput).To(MatchYAML(credYaml))
		})
	})
})
