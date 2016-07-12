package models

type Certificate struct {
	Root          string `json:"root,omitempty"`
	Certificate   string `json:"certificate,omitempty"`
	Private       string `json:"private,omitempty"`
}
