package commands

import "errors"

func NewNetworkError() error {
	return errors.New("No response received for the command. Please validate that you are targeting an active credential manager with `cm api` and retry your request.")
}

func NewSecretNotFoundError() error {
	return errors.New("Secret not found. Please validate your input and retry your request.")
}

func NewInvalidTargetError() error {
	return errors.New("The targeted API does not appear to be valid. Please validate the API address and retry your request.")
}

func NewInvalidStatusError() error {
	return errors.New("Unable to perform the request. Please validate your input and retry your request.")
}
