package commands

import "errors"

func NewNetworkError() error {
	return errors.New("No response received for the command. Please validate that you are targeting an active credential manager with `cm api` and retry your request.")
}
