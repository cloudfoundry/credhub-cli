package auth

import "net/http"

type MutualTls struct {
	Certificate string
}

// It will:
//   Use its ApiClient to complete the request
//   And returns the api response
func (c *MutualTls) Client() http.Client {
	panic("Not implemented")
}

// Constructs a func that will produce a TlsConfig using the Tls certificate
func MutualTlsCertificate(certificate string) AuthOption {
	panic("Not implemented")
}
