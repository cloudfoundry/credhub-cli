package commands

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
)

type LogoutCommand struct {
}

func (cmd LogoutCommand) Execute([]string) error {
	api.Logout()
	fmt.Println("Logout Successful")
	return nil
}
