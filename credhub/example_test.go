package credhub_test

import (
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth/uaa"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generate"
)

func ExampleCredHub() {
	return

	// Use a CredHub server on "https://example.com" using UAA password grant
	ch, err := credhub.New("https://example.com",
		credhub.SkipTLSValidation(),
		credhub.AuthBuilder(uaa.PasswordGrantBuilder("credhub_cli", "", "username", "password")))

	if err != nil {
		panic("credhub client configured incorrectly: " + err.Error())
	}

	authUrl, err := ch.AuthURL()
	if err != nil {
		panic("couldn't fetch authurl")
	}

	fmt.Println("CredHub server: ", ch.ApiURL)
	fmt.Println("Auth server: ", authUrl)

	// Retrieve a password stored at "/my/password"
	password, err := ch.GetPassword("/my/password")
	if err != nil {
		panic("password not found")
	}

	fmt.Println("My password: ", password.Value)

	// Manually refresh the access token
	uaa, ok := ch.Auth.(*auth.OAuthStrategy) // This works because we authenticated with auth.UaaPasswordGrant
	if !ok {
		panic("not using uaa")
	}

	fmt.Println("Old access token: ", uaa.AccessToken())

	uaa.Refresh() // For demo purposes only, tokens will be automatically refreshed by auth.OAuthStrategy

	fmt.Println("New access token:", uaa.AccessToken())
	// FIXME Output:
	// CredHub server: https://example.com
	// Auth server: https://uaa.example.com
	// My password: random-password
	// Old access token: some-access-token
	// New access token: new-access-token
}

func ExampleNew() {
	return

	ch, _ := credhub.New("https://example.com",
		credhub.SkipTLSValidation(),
		credhub.AuthBuilder(uaa.ClientCredentialsGrantBuilder("client-id", "client-secret")))

	fmt.Println("Connected to ", ch.ApiURL)
}

func ExampleCredHub_Request() {
	return

	ch, _ := credhub.New("https://example.com")

	// Get encryption key usage
	response, err := ch.Request("POST", "/api/v1/key-usage", nil, nil)
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
	ch, err := credhub.New("https://example.com",
		credhub.CaCerts(string("--- BEGIN ---\nroot-certificate\n--- END ---")),
		credhub.AuthBuilder(uaa.PasswordGrantBuilder("credhub_cli", "", "username", "password")))

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
