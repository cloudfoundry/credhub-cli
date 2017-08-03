package auth_test

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

func ExampleUaa() {
	// To retrieve the access token from the CredHub client, use type assertion
	uaa := ch.Auth.(auth.Uaa)

	fmt.Println("Before logging out: ", uaa.AccessToken)
	uaa.Logout()
	fmt.Println("After logging out: ", uaa.AccessToken)
	// Output:
	// Before logging out: some-access-token
	// After logging out:
}
