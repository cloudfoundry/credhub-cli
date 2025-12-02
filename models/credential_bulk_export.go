package models

import (
	"encoding/json"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"go.yaml.in/yaml/v3"
)

type exportCredential struct {
	Name     string
	Type     string
	Value    interface{}
	Metadata credentials.Metadata `json:",omitempty" yaml:",omitempty"`
}

type exportCredentials struct {
	Credentials []exportCredential
}

type CredentialBulkExport struct {
	Bytes []byte
}

func ExportCredentials(credentials []credentials.Credential, outputJSON bool) (*CredentialBulkExport, error) {
	exportCreds := exportCredentials{make([]exportCredential, len(credentials))}

	for i, credential := range credentials {
		exportCreds.Credentials[i] = exportCredential{
			credential.Name,
			credential.Type,
			credential.Value,
			credential.Metadata,
		}
	}

	result, err := yaml.Marshal(exportCreds)
	if outputJSON {
		result, err = json.Marshal(exportCreds)
	}

	if err != nil {
		return nil, err
	}

	return &CredentialBulkExport{result}, nil
}

func (credentialBulkExport *CredentialBulkExport) String() string {
	return string(credentialBulkExport.Bytes)
}
