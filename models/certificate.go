package models

type Certificate struct {
	Ca          string `json:"ca,omitempty"`
	Certificate string `json:"certificate,omitempty"`
	PrivateKey  string `json:"private_key,omitempty"`
}
