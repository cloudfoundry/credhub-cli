package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func PrintResponse(responseBody io.Reader) {
	responseBuffer := new(bytes.Buffer)
	responseBuffer.ReadFrom(responseBody)

	jsonBuffer := new(bytes.Buffer)
	json.Indent(jsonBuffer, responseBuffer.Bytes(), "", "  ")

	fmt.Println(jsonBuffer)
}
