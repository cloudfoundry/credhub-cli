package models

type Certificate struct {
	Root        string `json:"root,omitempty"`
	Certificate string `json:"certificate,omitempty"`
	PrivateKey  string `json:"private_key,omitempty"`
}
