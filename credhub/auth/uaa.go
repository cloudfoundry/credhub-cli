package auth

import (
	"encoding/json"
	"net/http"
)

// UAA authentication strategy
//
// Fields will be filled in based on Method used to construct strategy.
//
// When a UAA auth.Method (eg. UaaPasswordGrant()) is provided to credhub.New(),
// CredHub will use this Uaa.Do() to send authenticated requests to CredHub.
type Uaa struct {
	AccessToken  string
	RefreshToken string
	Username     string
	Password     string
	ClientId     string
	ClientSecret string
	ApiClient    HttpClient
	UaaClient    UaaClient
}

type UaaClient interface {
	ClientCredentialGrant(clientId, clientSecret string) (string, error)
	PasswordGrant(clientId, clientSecret, username, password string) (string, string, error)
	RefreshTokenGrant(clientId, clientSecret, refreshToken string) (string, string, error)
	RevokeToken(token string) error
}

// Provides http.Client-like interface to send requests authenticated with UAA
//
// Will automatically refresh the access token and retry the request if the token has expired.
func (a *Uaa) Do(req *http.Request) (*http.Response, error) {
	a.Login()

	req.Header.Set("Authorization", "Bearer "+a.AccessToken)
	resp, err := a.ApiClient.Do(req)

	if err == nil && tokenExpired(resp) {
		a.Refresh()

		req.Header.Set("Authorization", "Bearer "+a.AccessToken)
		resp, err = a.ApiClient.Do(req)
	}

	return resp, err
}

func tokenExpired(resp *http.Response) bool {
	if resp.StatusCode < 400 {
		return false
	}

	var errResp map[string]string

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&errResp)
	if err != nil {
		return false
	}

	return errResp["error"] == "access_token_expired"
}

// Refresh the access token
// If refresh token is available (ie. constructed with UaaPasswordGrant() or UaaSession()),
// a refresh grant will be used.
// Otherwise, client credential grant will be used to retrieve a new access token.
func (a *Uaa) Refresh() error {
	if a.RefreshToken == "" {
		a.AccessToken = ""

		return a.Login()
	}

	accessToken, refreshToken, _ := a.UaaClient.RefreshTokenGrant(a.ClientId, a.ClientSecret, a.RefreshToken)

	a.AccessToken = accessToken
	a.RefreshToken = refreshToken

	return nil
}

// Invalidate the access and refresh tokens on the UAA server
func (a *Uaa) Logout() {
	panic("Not implemented")
}

// Login will make a token grant request to the UAA server
//
// The grant type will be client credentials grant if either ClientID or ClientSecret is not empty,
// otherwise password grant will be used.
//
// On success, the AccessToken and RefreshToken (if given) will be populated.
//
// Login will be a no-op if the AccessToken is not-empty when invoked.
func (a *Uaa) Login() error {
	if a.AccessToken != "" {
		return nil
	}

	var accessToken string
	var refreshToken string

	if a.Username == "" {
		accessToken, _ = a.UaaClient.ClientCredentialGrant(a.ClientId, a.ClientSecret)
	} else {
		accessToken, refreshToken, _ = a.UaaClient.PasswordGrant(a.ClientId, a.ClientSecret, a.Username, a.Password)
	}

	a.AccessToken = accessToken
	a.RefreshToken = refreshToken

	return nil
}
