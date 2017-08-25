package auth

import "net/http"

type NoopStrategy struct {
	*http.Client
}

var _ Strategy = new(NoopStrategy)
