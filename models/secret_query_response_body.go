package models

import (
	"strings"
	"github.com/pivotal-cf/credhub-cli/util"
)

type SecretQueryResponseBody struct {
	Credentials []SecretQueryCredential `json:"credentials"`
}

type SecretQueryCredential struct {
	Name      string `json:"name,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func (secretQueryResponseBody SecretQueryResponseBody) String() string {
	lines := []string{}
	longestNameLength := 0
	for _, credential := range secretQueryResponseBody.Credentials {
		nameLength := len(credential.Name)
		if nameLength > longestNameLength {
			longestNameLength = nameLength
		}
	}
	lines = append(lines, util.BuildLineWithLength("Name", "Updated Date", longestNameLength))
	for _, credential := range secretQueryResponseBody.Credentials {
		lines = append(lines, util.BuildLineWithLength(credential.Name, credential.UpdatedAt, longestNameLength))
	}
	return strings.Join(lines, "\n")
}
