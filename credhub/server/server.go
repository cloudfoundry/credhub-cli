// CredHub server
package server

// CredHub server
type Server struct {
	// Url to CredHub server
	ApiUrl string
	// CA Certs in PEM format
	CaCerts []string
	// Skip certificate verification
	InsecureSkipVerify bool
}
