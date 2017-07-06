package models

type Certificate struct {
	Ca          string `json:"ca,omitempty" mapstructure:"ca,omitempty"`
	Certificate string `json:"certificate,omitempty" mapstructure:"certificate,omitempty"`
	PrivateKey  string `json:"private_key,omitempty" mapstructure:"private_key,omitempty"`
	CaName      string `json:"ca_name,omitempty" mapstructure:"ca_name,omitempty"`
}
