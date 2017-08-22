package credhub

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Request sends an authenticated request to the CredHub server.
//
// The pathStr should include the full path (eg. /api/v1/data) and any query parameters.
// The request body should be marshallable to JSON, but can be left nil for GET requests.
//
// Request() is used by other CredHub client methods to send authenticated requests to the CredHub server.
//
// Use Request() directly to access the CredHub server if an appropriate helper method is not available.
// For unauthenticated requests (eg. /health), use Config.Client() instead.
func (c *CredHub) Request(method string, pathStr string, query url.Values, body interface{}) (*http.Response, error) {
	u := *c.baseURL // clone
	u.Path = pathStr
	u.RawQuery = query.Encode()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), bytes.NewReader(jsonBody))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Auth.Do(req)

	if err != nil {
		return resp, err
	}

	if err := c.checkForServerError(resp); err != nil {
		return nil, err
	}

	return resp, err
}

func (c *CredHub) request(method string, path string, body io.Reader) (*http.Response, error) {
	client := c.Client()

	url := *c.baseURL // clone
	url.Path = path

	request, _ := http.NewRequest(method, url.String(), body)

	resp, err := client.Do(request)

	if err != nil {
		return resp, err
	}

	if err := c.checkForServerError(resp); err != nil {
		return nil, err
	}

	return resp, err
}

func (c *CredHub) checkForServerError(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)

		respErr := &Error{}

		if err := dec.Decode(respErr); err != nil {
			return err
		}

		return respErr
	}

	return nil
}
