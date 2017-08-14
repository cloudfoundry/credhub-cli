package auth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

// OAuth authentication strategy
//
// Fields will be filled in based on Builder used to construct strategy.
type OAuthStrategy struct {
	accessToken  string
	refreshToken string

	mu sync.RWMutex // guards AccessToken & Refresh Token

	Username     string
	Password     string
	ClientId     string
	ClientSecret string
	ApiClient    HttpClient
	OAuthClient  OAuthClient
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type OAuthClient interface {
	ClientCredentialGrant(clientId, clientSecret string) (string, error)
	PasswordGrant(clientId, clientSecret, username, password string) (string, string, error)
	RefreshTokenGrant(clientId, clientSecret, refreshToken string) (string, string, error)
	RevokeToken(token string) error
}

// Provides http.Client-like interface to send requests authenticated with OAuth
//
// Will automatically refresh the access token and retry the request if the token has expired.
func (a *OAuthStrategy) Do(req *http.Request) (*http.Response, error) {
	if err := a.Login(); err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+a.AccessToken())
	resp, err := a.ApiClient.Do(req)

	if err != nil {
		return resp, err
	}

	expired, err := tokenExpired(resp)

	if err != nil || !expired {
		return resp, err
	}

	if err := a.Refresh(); err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+a.AccessToken())
	return a.ApiClient.Do(req)
}

// Refresh the access token
// If refresh token is available (ie. constructed with UaaPasswordGrant() or UaaSession()),
// a refresh grant will be used.
// Otherwise, client_credentials grant type will be used to retrieve a new access token.
func (a *OAuthStrategy) Refresh() error {
	refreshToken := a.RefreshToken()
	if refreshToken == "" {
		return a.requestToken()
	}

	accessToken, refreshToken, err := a.OAuthClient.RefreshTokenGrant(a.ClientId, a.ClientSecret, refreshToken)

	if err != nil {
		return err
	}

	a.SetTokens(accessToken, refreshToken)

	return nil
}

// Invalidate the access and refresh tokens on the OAuth server
func (a *OAuthStrategy) Logout() {
	panic("Not implemented")
}

// Login will make a token grant request to the OAuth server
//
// The grant type will be client_credentials if either ClientID or ClientSecret is not empty,
// otherwise password grant type will be used.
//
// On success, the AccessToken and RefreshToken (if given) will be populated.
//
// Login will be a no-op if the AccessToken is not-empty when invoked.
func (a *OAuthStrategy) Login() error {
	if a.AccessToken() != "" {
		return nil
	}

	return a.requestToken()
}

func (a *OAuthStrategy) requestToken() error {
	var accessToken string
	var refreshToken string
	var err error

	if a.Username == "" {
		accessToken, err = a.OAuthClient.ClientCredentialGrant(a.ClientId, a.ClientSecret)
	} else {
		accessToken, refreshToken, err = a.OAuthClient.PasswordGrant(a.ClientId, a.ClientSecret, a.Username, a.Password)
	}

	if err != nil {
		return err
	}

	a.SetTokens(accessToken, refreshToken)

	return nil
}

func (a *OAuthStrategy) AccessToken() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.accessToken
}

func (a *OAuthStrategy) RefreshToken() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.refreshToken
}

func (a *OAuthStrategy) SetTokens(access, refresh string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.accessToken = access
	a.refreshToken = refresh
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

var _ Strategy = new(OAuthStrategy)
