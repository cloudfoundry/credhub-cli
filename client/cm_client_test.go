package client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/cm-cli/client"
	"github.com/pivotal-cf/cm-cli/mocks"
)

var _ = Describe("CM Client", func() {
	var (
		c          *client.CMClient
		httpClient *mocks.HttpClient
	)

	BeforeEach(func() {
		httpClient = &mocks.HttpClient{}
		c = &client.CMClient{
			HttpClient: httpClient,
		}
	})

	Describe("Set Secret", func() {
		It("should save the secret and return the key pairs", func() {
			httpClient.PutCall.ResponseJSON = `{ "values": {"key1" : "value1", "key2": "value2" } }`
			secretName := "secretName"
			secretPairs := make(map[string]string)
			secretPairs["key1"] = "value1"
			secretPairs["key2"] = "value2"
			secret := client.Secret{Values: secretPairs}
			response, err := c.SetSecrets(secretName, secretPairs)

			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.PutCall.Args.Route).To(Equal("/api/secret/" + secretName))
			Expect(response).NotTo(BeNil())
			Expect(response).To(Equal(secret))
		})
	})
})
