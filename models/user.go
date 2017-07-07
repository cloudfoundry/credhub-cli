package models

type User struct {
	Username string `json:"username,omitempty" mapstructure:"username,omitempty"`
	Password string `json:"password,omitempty" mapstructure:"password,omitempty"`
}
