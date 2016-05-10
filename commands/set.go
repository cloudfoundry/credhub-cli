package commands

import (
	"encoding/json"
	"fmt"

	"github.com/pivotal-cf/cm-cli/client"
)

type SetCommand struct {
	SecretIdentifier string `short:"i" long:"identifier" description:"The Identifier of the Secret"`
	SecretContent    string `short:"s" long:"secret" description:"The Content of the Secret"`
}

func (cmd SetCommand) Execute([]string) error {
	secret := client.SecretRequest{Value: cmd.SecretContent}

	secretJson, _ := json.MarshalIndent(secret, "", "  ")

	fmt.Println(string(secretJson))
	return nil
}
