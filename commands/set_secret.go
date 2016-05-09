package commands

import (
	"fmt"
	"github.com/pivotal-cf/cm-cli/client"
	"encoding/json"
)

type SetSecretCommand struct {
	SecretName string `short:"s" long:"secret" description:"The Secret Name"`
	KeyValues map[string]string `short:"k" long:"key-value" description:"The Key-Value pairs"`
}


func (cmd SetSecretCommand) Execute([]string) error {
	secret := client.SecretRequest{Values: cmd.KeyValues}
	secretJson, _ := json.MarshalIndent(secret, "", "  ")

	fmt.Println(string(secretJson))
	return nil
}
