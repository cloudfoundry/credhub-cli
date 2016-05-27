package client

import (
	"net/http"
	"time"
)

const TIMEOUT_SECS = 30

//go:generate counterfeiter . HttpClient

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

func NewHttpClient() HttpClient {
	return &http.Client{Timeout: time.Second * TIMEOUT_SECS}
}
