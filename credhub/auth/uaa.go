package auth

import "net/http"

type Uaa struct {
	AccessToken  string
	RefreshToken string
	Username     string
	Password     string
	ClientId     string
	ClientKey    string
}

// Does a Do method on a Config seem weird?
// It will:
//   Clone the request body for reuse later
//   Wrap the request headers with the appropriate tokens
//   Use its ApiClient to complete the request
//      on failure
//        refresh token
//        retry the request with cloned request body
//   And returns the api response
func (c *Uaa) Client() http.Client {
	panic("Not implemented")
}

// It will:
//   Use its ApiClient to complete a refresh token api request
//   And update the config's access/refresh tokens with the response
func (c *Uaa) Refresh() {
	panic("Not implemented")
}

// It will:
//   Use its ApiClient to complete a logout token api request
//   And invalidate the config's access/refresh tokens
func (c *Uaa) Logout() {
	panic("Not implemented")
}

func (c *Uaa) Login() {
	panic("Not implemented")
}

// Constructs a func that will produce a UaaConfig using password authentication
func UaaPassword(username, password string) AuthOption {
	panic("Not implemented")
}

// Constructs a func that will produce a UaaConfig using client credentials authentication
func UaaClient(clientId, clientKey string) AuthOption {
	panic("Not implemented")
}

// Constructs a func that will produce a UaaConfig using existing tokens
func UaaSession(accessToken, refreshToken string) AuthOption {
	panic("Not implemented")
}
