package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

type CurlCommand struct {
	Path          string `short:"p" long:"path" description:"The server endpoint to make the request against"`
	Method        string `short:"X" description:"HTTP method (default: GET)"`
	Data          string `short:"d" description:"HTTP data to include in the request body"`
	IncludeHeader bool   `short:"i" description:"Include the response headers in the output"`
}

func (cmd CurlCommand) Execute([]string) error {
	var err error

	cfg := config.ReadConfig()

	credhubClient, err := initializeCredhubClient(cfg)
	if err != nil {
		return err
	}

	if cmd.Path == "" {
		return errors.New("A path must be provided. Please update and retry your request.")
	}

	query := url.Values{}
	u, err := url.Parse(cmd.Path)
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
	if cmd.Data != "" {
		if err := json.Unmarshal([]byte(cmd.Data), &dat); err != nil {
			return err
		}
	}

	response, err := credhubClient.Request(cmd.Method, u.Path, query, dat, false)
	if err != nil {
		return err
	}
	if cmd.IncludeHeader {
		headers := new(bytes.Buffer)
		err = response.Header.Write(headers)

		if err != nil {
			return err
		}

		fmt.Print(response.Proto)
		fmt.Print(" ")
		fmt.Println(response.StatusCode)
		fmt.Println(response.TransferEncoding)
		fmt.Println(headers.String())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var bod map[string]interface{}
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
