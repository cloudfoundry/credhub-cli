package commands

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/api"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

func init() {
	CredHub.Token = func() {
		cfg := config.ReadConfig()
		if cfg.AccessToken != "" && cfg.AccessToken != "revoked" {

			_, err := api.NewApi(&cfg).Refresh()
			if err != nil {
				fmt.Println("Bearer " + cfg.AccessToken)
			}

			config.WriteConfig(cfg)

			fmt.Println("Bearer " + cfg.AccessToken)
		} else {
			fmt.Fprint(os.Stderr, "You are not currently authenticated. Please log in to continue.")
		}
		os.Exit(0)
	}
}
