package auth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// UAA authentication strategy
//
// Fields will be filled in based on Builder used to construct strategy.
//
// When a UAA auth.Builder (eg. UaaPasswordGrant()) is provided to credhub.New(),
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

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
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

	if err != nil {
		return resp, err
	}

	expired, err := tokenExpired(resp)

	if err == nil && expired {
		a.Refresh()
		req.Header.Set("Authorization", "Bearer "+a.AccessToken)
		resp, err = a.ApiClient.Do(req)
	}

	return resp, err
}

func tokenExpired(resp *http.Response) (bool, error) {
	if resp.StatusCode < 400 {
		return false, nil
	}

	var errResp map[string]string
	buf, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	decoder := json.NewDecoder(bytes.NewBuffer(buf))
	err = decoder.Decode(&errResp)

	if err != nil {
		// Since we fail to decode the error response
		// we cannot ensure that the token is invalid
		return false, nil
	}

	return errResp["error"] == "access_token_expired", nil
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

	accessToken, refreshToken, err := a.UaaClient.RefreshTokenGrant(a.ClientId, a.ClientSecret, a.RefreshToken)

	if err != nil {
		return err
	}

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
	var err error

	if a.Username == "" {
		accessToken, err = a.UaaClient.ClientCredentialGrant(a.ClientId, a.ClientSecret)
	} else {
		accessToken, refreshToken, err = a.UaaClient.PasswordGrant(a.ClientId, a.ClientSecret, a.Username, a.Password)
	}

	if err != nil {
		return err
	}

	a.AccessToken = accessToken
	a.RefreshToken = refreshToken

	return nil
}
