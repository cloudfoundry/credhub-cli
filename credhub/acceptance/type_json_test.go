package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSON Credential Type", func() {
	Specify("lifecycle", func() {
		name := testCredentialPath("some-json")
		cred := make(map[string]interface{})
		cred["key"] = "value"

		cred2 := make(map[string]interface{})
		cred2["another_key"] = "another value"

		By("setting the json for the first time returns same json")
		json, err := credhubClient.SetJSON(name, cred, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(json.Value).To(Equal(cred))

		By("setting the json again without overwrite returns same json")
		json, err = credhubClient.SetJSON(name, cred2, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(json.Value).To(Equal(cred))

		By("overwriting the json with set")
		json, err = credhubClient.SetJSON(name, cred2, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(json.Value).To(Equal(cred2))

		By("getting the json")
		json, err = credhubClient.GetJSON(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(json.Value).To(Equal(cred2))

		By("deleting the json")
		err = credhubClient.Delete(name)
		Expect(err).ToNot(HaveOccurred())
		_, err = credhubClient.GetJSON(name)
		Expect(err).To(HaveOccurred())
	})
})
