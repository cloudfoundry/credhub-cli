package models

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Credential struct {
	Name  string      `yaml:"name"`
	Type  string      `yaml:"type"`
	Value interface{} `yaml:"value"`
}

type CredentialBulkImport struct {
	Credentials []Credential `yaml:"credentials"`
}

func (credentialBulkImport *CredentialBulkImport) ReadFile(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return credentialBulkImport.ReadBytes(data)
}

func (credentialBulkImport *CredentialBulkImport) ReadBytes(data []byte) error {
	return yaml.Unmarshal(data, credentialBulkImport)
}
