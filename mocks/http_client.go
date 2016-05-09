package mocks

import "encoding/json"

type HttpClient struct {
	PutCall struct {
		Args struct {
			Route        string
			RequestData  interface{}
			ResponseData interface{}
		}
		Return struct {
			Error error
		}

		ResponseJSON string
	}
}

func (c *HttpClient) Put(route string, requestData interface{}, responseData interface{}) error {
	c.PutCall.Args.Route = route
	c.PutCall.Args.RequestData = requestData
	c.PutCall.Args.ResponseData = responseData
	err := json.Unmarshal([]byte(c.PutCall.ResponseJSON), responseData)
	if err != nil {
		panic(err)
	}
	return c.PutCall.Return.Error
}
