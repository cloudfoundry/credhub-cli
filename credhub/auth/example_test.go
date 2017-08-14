package auth_test

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

func ExampleOAuth() {
	return

	// To retrieve the access token from the CredHub client, use type assertion
	ch := credhub.CredHub{}
	oauth, ok := ch.Auth.(*auth.OAuthStrategy)
	if !ok {
		panic("Not using UAA")
	}

	fmt.Println("Before logging out: ", oauth.AccessToken())
	oauth.Logout()
	fmt.Println("After logging out: ", oauth.AccessToken())
	// FIXME Output:
	// Before logging out: some-access-token
	// After logging out:
}
