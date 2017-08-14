package uaa

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

// Client makes requests to the UAA server at AuthUrl
type Client struct {
	AuthUrl string
	Client  *http.Client
}

type token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// ClientCredentialGrant requests a token using client credential grant
func (u *Client) ClientCredentialGrant(clientId, clientSecret string) (string, error) {
	values := url.Values{
		"grant_type":    {"client_credentials"},
		"response_type": {"token"},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}

	token, err := u.tokenGrantRequest(values)

	return token.AccessToken, err
}

// ClientCredentialGrant requests an access token and refresh token using password grant
func (u *Client) PasswordGrant(clientId, clientSecret, username, password string) (string, string, error) {
	values := url.Values{
		"grant_type":    {"password"},
		"response_type": {"token"},
		"username":      {username},
		"password":      {password},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}

	token, err := u.tokenGrantRequest(values)

	return token.AccessToken, token.RefreshToken, err
}

// RefreshTokenGrant requests a new access token and refresh token using refresh token grant
func (u *Client) RefreshTokenGrant(clientId, clientSecret, refreshToken string) (string, string, error) {
	values := url.Values{
		"grant_type":    {"refresh_token"},
		"response_type": {"token"},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"refresh_token": {refreshToken},
	}

	token, err := u.tokenGrantRequest(values)

	return token.AccessToken, token.RefreshToken, err
}

// TokenGrantRequest requests a new token with the given request headers
func (u *Client) tokenGrantRequest(headers url.Values) (token, error) {
	request, _ := http.NewRequest("POST", u.AuthUrl+"/oauth/token", bytes.NewBufferString(headers.Encode()))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := u.Client.Do(request)

	var t token

	if err != nil {
		return t, err
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&t)

	return t, err
}

// RevokeToken revokes the given access token
func (u *Client) RevokeToken(accessToken string) (err error) {
	panic("not implemented")
}
