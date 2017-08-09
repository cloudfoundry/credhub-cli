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

// ClientCredentialGrant requests a token using client credential grant
func (u *Client) ClientCredentialGrant(clientId, clientSecret string) (string, error) {
	values := url.Values{
		"grant_type":    {"client_credentials"},
		"response_type": {"token"},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}

	token, _ := u.TokenGrantRequest(values)

	return token.AccessToken, nil
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

	token, _ := u.TokenGrantRequest(values)

	return token.AccessToken, token.RefreshToken, nil
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

	token, _ := u.TokenGrantRequest(values)

	return token.AccessToken, token.RefreshToken, nil
}

// TokenGrantRequest requests a new token with the given request headers
func (u *Client) TokenGrantRequest(headers url.Values) (Token, error) {
	request, _ := http.NewRequest("POST", u.AuthUrl+"/oauth/token", bytes.NewBufferString(headers.Encode()))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, _ := u.Client.Do(request)

	var token Token
	decoder := json.NewDecoder(response.Body)
	_ = decoder.Decode(&token)

	return token, nil
}

// RevokeToken revokes the given access token
func (u *Client) RevokeToken(accessToken string) (err error) {
	panic("not implemented")
}
