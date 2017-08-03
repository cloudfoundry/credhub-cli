package auth

import "net/http"

type MutualTls struct {
	Certificate string
}

func (a *MutualTls) Do(http.Request) (http.Response, error) {
	panic("Not implemented")
}

// Constructs a func that will produce a TlsConfig using the Tls certificate
func MutualTlsCertificate(certificate string) AuthOption {
	panic("Not implemented")
}
