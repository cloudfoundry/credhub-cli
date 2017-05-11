package commands

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

func init() {
	CredHub.Token = func() {
		cfg := config.ReadConfig()
		if cfg.AccessToken != "" {
			fmt.Println("Bearer " + cfg.AccessToken)
		} else {
			fmt.Fprint(os.Stderr, "You are not logged in to a CredHub server")
		}
		os.Exit(0)
	}
}
