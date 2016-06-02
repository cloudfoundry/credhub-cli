package models

type Certificate struct {
	Ca      string `json:"ca,omitempty"`
	Public  string `json:"public,omitempty"`
	Private string `json:"private,omitempty"`
}
