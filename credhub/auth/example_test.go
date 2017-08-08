package auth_test

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

func ExampleUaa() {
	return

	// To retrieve the access token from the CredHub client, use type assertion
	ch := credhub.CredHub{}
	uaa, ok := ch.Auth.(*auth.Uaa)
	if !ok {
		panic("Not using UAA")
	}

	fmt.Println("Before logging out: ", uaa.AccessToken)
	uaa.Logout()
	fmt.Println("After logging out: ", uaa.AccessToken)
	// Output:
	// Before logging out: some-access-token
	// After logging out:
}
