package credhub

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub/server"
)

func (c *CredHub) Info() (*server.Info, error) {
	response, err := c.request("GET", "/info", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	info := &server.Info{}
	decoder := json.NewDecoder(response.Body)

	if err = decoder.Decode(&info); err != nil {
		return nil, err
	}

	return info, nil
}
