package credhub_test

import (
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
)

func ExampleCredHub() {
	return

	// Use a CredHub server on "https://example.com" using UAA password grant
	server := credhub.Config{
		ApiUrl:             "https://example.com",
		InsecureSkipVerify: true,
	}
	authOption := auth.UaaPasswordGrant("credhub_cli", "", "username", "password")

	ch, err := credhub.New(&server, authOption)

	if err != nil {
		panic("credhub client configured incorrectly: " + err.Error())
	}

	authUrl, err := ch.AuthUrl()
	if err != nil {
		panic("couldn't fetch authurl")
	}

	fmt.Println("CredHub server: ", ch.ApiUrl)
	fmt.Println("Auth server: ", authUrl)

	// Retrieve a password stored at "/my/password"
	password, err := ch.GetPassword("/my/password")
	if err != nil {
		panic("password not found")
	}

	fmt.Println("My password: ", password.Value)

	// Manually refresh the access token
	uaa, ok := ch.Auth.(*auth.Uaa) // This works because we authenticated with auth.UaaPasswordGrant
	if !ok {
		panic("not using uaa")
	}

	fmt.Println("Old access token: ", uaa.AccessToken)

	uaa.Refresh() // For demo purposes only, tokens will be automatically refreshed by auth.Uaa

	fmt.Println("New access token:", uaa.AccessToken)
	// FIXME Output:
	// CredHub server: https://example.com
	// Auth server: https://uaa.example.com
	// My password: random-password
	// Old access token: some-access-token
	// New access token: new-access-token
}

func ExampleNew() {
	return

	server := credhub.Config{
		ApiUrl:             "https://example.com",
		InsecureSkipVerify: true,
	}
	authOption := auth.UaaClientCredentialGrant("client-id", "client-secret")

	ch, _ := credhub.New(&server, authOption)

	fmt.Println("Connected to ", ch.ApiUrl)
}

func ExampleCredHub_Request() {
	return

	ch := credhub.CredHub{}

	// Get encryption key usage
	response, err := ch.Request("POST", "/api/v1/key-usage", nil)
	if err != nil {
		panic("couldn't get key usage")
	}

	var keyUsage map[string]int
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&keyUsage)
	if err != nil {
		panic("couldn't parse response")
	}

	fmt.Println("Active Key: ", keyUsage["active_key"])
	// FIXME Output:
	// Active Key: 1231231
}

func Example() {
	return

	// CredHub server at https://example.com, using UAA Password grant
	ch, _ := credhub.New(
		&credhub.Config{
			ApiUrl:  "https://example.com",
			CaCerts: []string{"--- BEGIN ---\nroot-certificate\n--- END ---"},
		},
		auth.UaaPasswordGrant("credhub_cli", "", "username", "password"),
	)

	// We'll be working with a certificate stored at "/my-certificates/the-cert"
	path := "/my-certificates/"
	name := "the-cert"

	// If the certificate already exists, delete it
	cert, err := ch.GetCertificate(path + name)
	if err == nil {
		ch.Delete(cert.Name)
	}

	// Generate a new certificate
	gen := generate.Certificate{
		CommonName: "pivotal",
		KeyLength:  2048,
	}
	cert, err = ch.GenerateCertificate(path+name, gen, false)
	if err != nil {
		panic("couldn't generate certificate")
	}

	// Use the generated certificate's values to create a new certificate
	dupCert, err := ch.SetCertificate(path+"dup-cert", cert.Value, false)
	if err != nil {
		panic("couldn't create certificate")
	}

	if dupCert.Value.Certificate != cert.Value.Certificate {
		panic("certs don't match")
	}

	// List all credentials in "/my-certificates"
	creds, err := ch.FindByPath(path)
	if err != nil {
		panic("couldn't list certificates")
	}

	fmt.Println("Found the following credentials in " + path + ":")
	for _, cred := range creds {
		fmt.Println(cred.Name)
	}
	// FIXME Output:
	// Found the following credentials in /my-certificates:
	// /my-certificates/dup-cert
	// /my-certificates/the-cert
}
