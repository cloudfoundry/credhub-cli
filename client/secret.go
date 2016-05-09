package client

type SecretRequest struct {
	Values map[string]string `json:"values"`
}

type SecretResponse struct {
	Values map[string]string  `json:"values"`
}
