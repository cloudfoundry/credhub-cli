package models

type RsaSsh struct {
	PublicKey  string `json:"public_key,omitempty" yaml:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty" yaml:"private_key,omitempty"`
}
