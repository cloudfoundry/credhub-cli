package models

type CaParameters struct {
	Certificate string `json:"certificate,omitempty"`
	PrivateKey  string `json:"private_key,omitempty"`
}
