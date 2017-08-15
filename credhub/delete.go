package credhub

import (
	"net/http"
	"net/url"
)

// Deletes a credential by name
func (ch *CredHub) Delete(name string) error {
	query := url.Values{}
	query.Set("name", name)
	resp, err := ch.Request(http.MethodDelete, "/api/v1/data", query, nil)

	if err == nil {
		defer resp.Body.Close()
	}

	return err
}
