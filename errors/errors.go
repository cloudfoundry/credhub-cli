package errors

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/pivotal-cf/cm-cli/models"
)

func NewNetworkError() error {
	return errors.New("No response received for the command. Please validate that you are targeting an active credential manager with `cm api` and retry your request.")
}

func NewResponseError() error {
	return errors.New("An error occurred when processing the response. Please validate your input and retry your request.")
}

func NewInvalidTargetError() error {
	return errors.New("The targeted API does not appear to be valid. Please validate the API address and retry your request.")
}

func NewNoTargetUrlError() error {
	return errors.New("Credential Manager API location is not set. Please target the location of your credential manager with `cm api` to continue.")
}

func NewFileLoadError() error {
	return errors.New("A referenced file could not be opened. Please validate the provided filenames and permissions, then retry your request.")
}

func NewCombinationOfParametersError() error {
	return errors.New("The combination of parameters in the request is not allowed. Please validate your input and retry your request.")
}

func NewAuthorizationError() error {
	return errors.New("The provided username and password combination are incorrect. Please validate your input and retry your request.")
}

func NewAuthorizationParametersError() error {
	return errors.New("The combination of parameters in the request is not allowed. Please validate your input and retry your request.")
}

func ParseError(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	serverError := models.ServerError{}
	decoder.Decode(&serverError)

	return errors.New(serverError.Error)
}
