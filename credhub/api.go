// CredHub provides methods to interact with CredHub server
package credhub

type Config struct {
	ApiURL             string
	AuthURL            string
	InsecureSkipVerify bool
	CaCerts            []string
	AccessToken        string
	RefreshToken       string
}

type CredHub struct {
	*Config
}

type Server struct {
	ApiURL             string
	InsecureSkipVerify bool
	CaCerts            []string
}

func New() *CredHub {
	panic("Not implemented")
}

func NewWithConfig(config *Config) *CredHub {
	panic("Not implemented")
}

func (ch *CredHub) Target() error {
	panic("Not implemented")
}
