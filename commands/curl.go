package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
)

type CurlCommand struct {
	Path          string `short:"p" long:"path" description:"The server endpoint to make the request against"`
	Fail          bool   `short:"f" long:"fail" description:"Fail silently (no output at all) on HTTP errors"`
	Method        string `short:"X" description:"HTTP method (default: GET)"`
	Data          string `short:"d" description:"HTTP data to include in the request body"`
	IncludeHeader bool   `short:"i" description:"Include the response headers in the output"`
	ClientCommand
}

func (c *CurlCommand) Execute([]string) error {
	if c.Path == "" {
		return errors.New("A path must be provided. Please update and retry your request.")
	}

	query := url.Values{}
	u, err := url.Parse(c.Path)
	if err != nil {
		return err
	}

	if u.RawQuery != "" {
		query, err = url.ParseQuery(u.RawQuery)
		if err != nil {
			return err
		}
	}

	var dat map[string]interface{}
	if c.Data != "" {
		if err := json.Unmarshal([]byte(c.Data), &dat); err != nil {
			return err
		}
	}

	response, err := c.client.Request(c.Method, u.Path, query, dat, false)
	if err != nil {
		return err
	}

	if c.IncludeHeader {
		headers := new(bytes.Buffer)
		err = response.Header.Write(headers)

		if err != nil {
			return err
		}

		fmt.Print(response.Proto)
		fmt.Print(" ")
		fmt.Println(response.StatusCode)
		fmt.Println(headers.String())
	}

	if c.Fail && response.StatusCode >= 400 && response.StatusCode < 600 {
		os.Exit(22)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var bod interface{}
	err = json.Unmarshal(body, &bod)
	if err != nil {
		return err
	}

	formattedBody, err := json.MarshalIndent(bod, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(formattedBody))

	return err
}
