package client

import "net/http"

//go:generate counterfeiter . HttpClient

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}
