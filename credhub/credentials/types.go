package credentials

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"

// Base fields of a credential
type Base struct {
	Name             string
	VersionCreatedAt string
}

// A generic credential
type Credential struct {
	Base
	Id    string
	Type  string
	Value interface{}
}

// A Value type credential
type Value struct {
	Credential
	Value values.Value
}

// A JSON type credential
type JSON struct {
	Credential
	Value values.JSON
}

// A Password type credential
type Password struct {
	Credential
	Value values.Password
}

// A User type credential
type User struct {
	Credential
	Value struct {
		values.User
		PasswordHash string
	}
}

// A Certificate type credential
type Certificate struct {
	Credential
	Value values.Certificate
}

// An RSA type credential
type RSA struct {
	Credential
	Value values.RSA
}

// An SSH type credential
type SSH struct {
	Credential
	Value values.SSH
}
