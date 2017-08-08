package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
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
	ApiClient    *http.Client
	AuthUrl      string
}

// Provides http.Client-like interface to send requests authenticated with UAA
//
// Will automatically refresh the access token and retry the request if the token has expired.
func (a Uaa) Do(req *http.Request) (*http.Response, error) {
	panic("Not implemented yet")
}

// Refresh the access token
// If refresh token is available (ie. constructed with UaaPasswordGrant() or UaaSession()),
// a refresh grant will be used.
// Otherwise, client credential grant will be used to retrieve a new access token.
func (a *Uaa) Refresh() {
	panic("Not implemented yet")
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
	var values url.Values
	if a.Username == "" {
		values = url.Values{
			"grant_type":    {"client_credentials"},
			"response_type": {"token"},
			"client_id":     {a.ClientId},
			"client_secret": {a.ClientSecret},
		}
	} else {
		values = url.Values{
			"grant_type":    {"password"},
			"response_type": {"token"},
			"username":      {a.Username},
			"password":      {a.Password},
			"client_id":     {a.ClientId},
			"client_secret": {a.ClientSecret},
		}
	}

	request, _ := http.NewRequest("POST", a.AuthUrl+"/oauth/token", bytes.NewBufferString(values.Encode()))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := a.ApiClient.Do(request)
	if err != nil {
		panic(err)
	}
	var token struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
	}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&token)
	if err != nil {
		panic(err)
	}
	a.AccessToken = token.AccessToken
	a.RefreshToken = token.RefreshToken
	return nil
}

// Provides a constructor for a UAA authentication strategy using password grant
func UaaPasswordGrant(clientId, clientSecret, username, password string) Method {
	return func(s *server.Server) Auth {
		panic("Not implemented")
	}
}

// UaaClientCredentialGrant provides a constructor for a UAA authentication strategy
// using client credential grant.
func UaaClientCredentialGrant(clientId, clientSecret string) Method {
	return func(s *server.Server) Auth {
		panic("Not implemented")
	}
}

// Provides a constructor for a UAA authentication strategy using existing tokens
//
// Only use this for a sessions created with password grant
// For existing sessions created with a client credential grant, use UaaClient()
func UaaSession(clientId, clientSecret, accessToken, refreshToken string) Method {
	return func(s *server.Server) Auth {
		panic("Not implemented")
	}
}
