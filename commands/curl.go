package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

type CurlCommand struct {
	Path   string `short:"p" long:"path" description:"The server endpoint to make the request against"`
	Method string `short:"X" description:"HTTP method (default: GET)"`
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

	response, err := credhubClient.Request(cmd.Method, u.Path, query, nil, false)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return err
}
