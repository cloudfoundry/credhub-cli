package commands

import (
	"encoding/json"
	"fmt"
	"github.com/pivotal-cf/cm-cli/client"
)

func PrintResponse(responseBytes []byte) {
	secret := new(client.Secret)

	json.Unmarshal(responseBytes, &secret)

	fmt.Println(fmt.Sprintf("Name:	%s\nValue:	%s", secret.Name, secret.Value))
}
