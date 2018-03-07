package commands

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

type CurlCommand struct {
	Path string `short:"p" long:"path" description:"The server endpoint to make the request against"`
	//Id               string `long:"id" description:"ID of the credential to retrieve"`
	//NumberOfVersions int    `long:"versions" description:"Number of versions of the credential to retrieve"`
	//OutputJson       bool   `short:"j" long:"output-json" description:"Return response in JSON format"`
	//Key              string `short:"k" long:"key" description:"Return only the specified field of the requested credential"`
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

	response, err := credhubClient.Request("GET", cmd.Path, nil, nil)
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
