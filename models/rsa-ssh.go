package models

type RsaSsh struct {
	PublicKey  string `json:"public_key,omitempty" mapstructure:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty" mapstructure:"private_key,omitempty"`
}
