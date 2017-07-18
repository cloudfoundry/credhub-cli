package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/config"
)

type Api struct {
	Config *config.Config
}

func NewApi(cfg *config.Config) *Api {
	if cfg == nil {
		cfg = &config.Config{}
	}

	return &Api{cfg}
}
