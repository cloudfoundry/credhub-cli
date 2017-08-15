package credhub

import "encoding/json"

type Info struct {
	App struct {
		Name    string
		Version string
	}
	AuthServer struct {
		Url string
	} `json:"auth-server"`
}

func (c *CredHub) Info() (*Info, error) {
	response, err := c.request("GET", "/info", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	info := &Info{}
	decoder := json.NewDecoder(response.Body)

	if err = decoder.Decode(&info); err != nil {
		return nil, err
	}

	return info, nil
}
