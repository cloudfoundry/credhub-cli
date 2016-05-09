package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	BaseURL       string
	SkipTLSVerify bool
}

func (c *HttpClient) Put(route string, requestData interface{}, responseData interface{}) error {
	reqBody, err := jsonBody(requestData)
	if err != nil {
		return err
	}

	req, err := c.jsonRequest(route, reqBody)
	if err != nil {
		return err
	}

	resp, err := c.client().Do(req)
	if err != nil {
		return err
	}

	err = handleResponse(resp, responseData)
	if err != nil {
		return fmt.Errorf("Error performing PUT:\n%s", err.Error())
	}
	return nil
}

func jsonBody(requestData interface{}) (*bytes.Reader, error) {
	var bodyBytes []byte
	var err error
	if requestData != nil {
		bodyBytes, err = json.Marshal(requestData)
		if err != nil {
			return nil, err // not covered by tests
		}
	}
	return bytes.NewReader(bodyBytes), nil
}

func (c *HttpClient) jsonRequest(route string, reqBody *bytes.Reader) (*http.Request, error) {
	req, err := c.request("PUT", route, reqBody)
	if err != nil {
		return nil, err
	}

	if reqBody.Len() > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func handleResponse(resp *http.Response, responseData interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Status code %d.\nBody: %s", resp.StatusCode, body)
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, responseData)
		if err != nil {
			return fmt.Errorf("server returned malformed JSON: %s", err)
		}
	}

	return nil
}

func (c *HttpClient) client() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.SkipTLSVerify},
	}
	return &http.Client{Transport: tr}
}

func (c *HttpClient) request(reqVerb string, route string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(reqVerb, c.BaseURL+route, body)
	if err != nil {
		return nil, err
	}

	return req, nil
}
