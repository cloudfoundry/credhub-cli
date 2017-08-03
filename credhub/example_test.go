package credhub_test

import (
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
)

func ExampleCredHub() {
	// Use a CredHub server on "https://example.com" using UAA password grant
	server := server.Server{
		ApiUrl:             "https://example.com",
		InsecureSkipVerify: true,
	}
	authOption := auth.UaaPasswordGrant("username", "password")

	ch := credhub.New(server, authOption)

	fmt.Println("CredHub server: ", ch.ApiUrl)
	fmt.Println("Auth server: ", ch.AuthUrl())

	// Retrieve a password stored at "/my/password"
	password, err := ch.GetPassword("/my/password")
	if err != nil {
		panic("password not found")
	}

	fmt.Println("My password: ", password.Value)

	// Manually refresh the access token
	uaa, ok := ch.Auth.(auth.Uaa) // This works because we authenticated with auth.UaaPasswordGrant
	if !ok {
		panic("not using uaa")
	}

	fmt.Println("Old access token: ", uaa.AccessToken)

	uaa.Refresh() // For demo purposes only, tokens will be automatically refreshed by auth.Uaa

	fmt.Println("New access token:", uaa.AccessToken)
	// Output:
	// CredHub server: https://example.com
	// Auth server: https://uaa.example.com
	// My password: random-password
	// Old access token: some-access-token
	// New access token: new-access-token
}

func ExampleNew() {
	server := server.Server{
		ApiUrl:             "https://example.com",
		InsecureSkipVerify: true,
	}
	authOption := auth.UaaClientCredentialGrant("client-id", "client-secret")

	ch := credhub.New(server, authOption)
}

func ExampleCredHub_Request() {
	// Get encryption key usage
	response, err := ch.Request("POST", "/api/v1/key-usage", nil)
	if err != nil {
		panic("couldn't get key usage")
	}

	var keyUsage map[string]int
	json.Unmarshal(response.Body, &keyUsage)

	fmt.Println("Active Key: ", keyUsage["active_key"])
	// Output:
	// Active Key: 1231231
}

func Example() {
	// CredHub server at https://example.com, using UAA Password grant
	ch := credhub.New(
		server.Server{
			ApiUrl:  "https://example.com",
			CaCerts: []string{"--- BEGIN ---\nroot-certificate\n--- END ---"},
		},
		auth.UaaPasswordGrant("username", "password"),
	)

	// We'll be working with a certificate stored at "/my-certificates/the-cert"
	path := "/my-certificates/"
	name := "the-cert"

	// If the certificate already exists, delete it
	cert, err := ch.GetCertificate(path + name)
	if err == nil {
		ch.Delete(cert)
	}

	// Generate a new certificate
	gen := generate.Certificate{
		CommonName: "pivotal",
		KeyLength:  2048,
	}
	cert, err := ch.GenerateCertificate(path+name, gen, false)
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
	// Output:
	// Found the following credentials in /my-certificates:
	// /my-certificates/dup-cert
	// /my-certificates/the-cert
}
