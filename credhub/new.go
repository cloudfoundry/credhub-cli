package credhub

import "net/url"

// New creates a new CredHub API client with the provided server credentials and authentication method.
// See the auth package for supported authentication methods.
func New(addr string, options ...Option) (*CredHub, error) {
	baseURL, err := url.Parse(addr)

	if err != nil {
		return nil, err
	}

	credhub := &CredHub{
		ApiURL:  addr,
		baseURL: baseURL,
	}

	for _, option := range options {
		if err := option(credhub); err != nil {
			return nil, err
		}
	}

	if credhub.baseURL.Scheme == "https" {
		credhub.defaultClient = httpsClient(credhub.insecureSkipVerify, credhub.caCerts)
	} else {
		credhub.defaultClient = httpClient()
	}

	if credhub.Auth != nil {
		return credhub, nil
	}

	if credhub.authBuilder == nil {
		credhub.Auth = credhub.defaultClient
		return credhub, nil
	}

	credhub.Auth, err = credhub.authBuilder(credhub)

	if err != nil {
		return nil, err
	}

	return credhub, nil
}
