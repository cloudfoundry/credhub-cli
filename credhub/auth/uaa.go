package auth

import "net/http"

// UAA authentication strategy
//
// Fields will be filled in based on Method used to construct strategy.
//
// When a UAA auth.Method (eg. UaaPasswordGrant()) is provided to credhub.New(),
// CredHub will use this Uaa.Do() to send authenticated requests to CredHub.
// The UAA server will be retrieved from the server.Server.
type Uaa struct {
	AccessToken  string
	RefreshToken string
	Username     string
	Password     string
	ClientId     string
	ClientKey    string
}

// Provides http.Client-like interface to send requests authenticated with UAA
//
// Will automatically refresh the access token and retry the request if the token has expired.
func (a *Uaa) Do(*http.Request) (*http.Response, error) {
	panic("Not implemented")
}

// Refresh the access token
// If refresh token is available (ie. constructed with UaaPasswordGrant() or UaaSession()),
// a refresh grant will be used.
// Otherwise, client credential grant will be used to retrieve a new access token.
func (c *Uaa) Refresh() {
	panic("Not implemented")
}

// Invalidate the access and refresh tokens on the UAA server
func (c *Uaa) Logout() {
	panic("Not implemented")
}

// Retrieve access token and refresh token using password grant
func (c *Uaa) Login() {
	panic("Not implemented")
}

// Provides a constructor for a UAA authentication strategy using password grant
func UaaPasswordGrant(username, password string) Method {
	panic("Not implemented")
}

// Provides a constructor for a UAA authentication strategy using client credential grant
func UaaClientCredentialGrant(clientId, clientKey string) Method {
	panic("Not implemented")
}

// Provides a constructor for a UAA authentication strategy using existing tokens
//
// Only use this for a sessions created with password grant
// For existing sessions created with a client credential grant, use UaaClient()
func UaaSession(accessToken, refreshToken string) Method {
	panic("Not implemented")
}
