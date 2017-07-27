package auth

type Auth interface {
	Headers() map[string]string
}

type UAA struct {
	AuthURL      string
	AccessToken  string
	RefreshToken string
}

type TLS struct {
	Certificate string
	PrivateKey  string
}

func (ch *UAA) LoginWithPassword(username, password string) error {
	panic("Not implemented")
}

func (ch *UAA) LoginWithClientCredentials(clientName, clientSecret string) error {
	panic("Not implemented")
}

func (ch *UAA) Refresh() error {
	panic("Not implemented")
}

func (ch *UAA) Logout() {
	panic("Not implemented")
}

func NewUaaAuth(username, password string) *Auth {
	panic("Not implemented")
}

func NewTlsAuth(certificate, privateKey string) *Auth {
	panic("Not implemented")
}
