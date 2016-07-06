package commands

import (
	"io/ioutil"

	cmcli_errors "github.com/pivotal-cf/cm-cli/errors"
)

func ReadFile(filename string) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", cmcli_errors.NewFileLoadError()
	}
	return string(dat), nil
}
