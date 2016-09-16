package models

import (
	"fmt"
	"strconv"
	"strings"
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
	lines = append(lines, buildLineWithLength("Name", "Updated Date", longestNameLength))
	for _, credential := range secretQueryResponseBody.Credentials {
		lines = append(lines, buildLineWithLength(credential.Name, credential.UpdatedAt, longestNameLength))
	}
	return strings.Join(lines, "\n")
}

func buildLineWithLength(leftItem string, rightItem string, leftItemSize int) string {
	return fmt.Sprintf("%-" + strconv.Itoa(leftItemSize) + "s   %s", leftItem, rightItem)
}
