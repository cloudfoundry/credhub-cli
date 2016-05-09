package client

type CMClient struct {
	HttpClient httpClient
}

type Secret struct {
	Values map[string]string `json:"values"`
}

//go:generate counterfeiter . httpClient

type httpClient interface {
	Put(route string, requestData interface{}, responseData interface{}) error
}

func (c *CMClient) SetSecrets(secretName string, kvs map[string]string) (Secret, error) {
	requestData := Secret{Values: kvs}
	responseData := Secret{Values: make(map[string]string)}
	route := "/api/secret/" + secretName
	err := c.HttpClient.Put(route, requestData, &responseData)
	return responseData, err
}
