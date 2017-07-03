package models

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Credential struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

type ImportFile struct {
	Credentials []Credential `yaml:"credentials"`
}

func (importFile *ImportFile) ReadFile(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	return importFile.ReadBytes(data)
}

func (importFile *ImportFile) ReadBytes(data []byte) error {
	return yaml.Unmarshal(data, importFile)
}
