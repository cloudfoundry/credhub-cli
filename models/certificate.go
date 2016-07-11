package models

type Certificate struct {
	Ca          string `json:"ca,omitempty"`
	Certificate string `json:"certificate,omitempty"`
	Private     string `json:"private,omitempty"`
}
