package client

type SecretRequest struct {
	Value string `json:"value"`
}

type SecretResponse struct {
	Value string `json:"value"`
}
