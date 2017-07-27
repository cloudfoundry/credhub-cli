package credentials

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"

type Path struct {
	Path string
}

type Base struct {
	Name             string
	VersionCreatedAt string
}

type Credential struct {
	Base
	Id    string
	Type  string
	Value interface{}
}

type Value struct {
	Credential
	Value values.Value
}

type JSON struct {
	Credential
	Value values.JSON
}

type Password struct {
	Credential
	Value values.Password
}

type User struct {
	Credential
	Value struct {
		values.User
		PasswordHash string
	}
}

type Certificate struct {
	Credential
	Value values.Certificate
}

type RSA struct {
	Credential
	Value values.RSA
}

type SSH struct {
	Credential
	Value values.SSH
}
