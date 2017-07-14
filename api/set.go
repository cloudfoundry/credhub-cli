package api

import (
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	"github.com/cloudfoundry-incubator/credhub-cli/errors"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func Set(
	credentialIdentifier string,
	credentialType string,
	noOverwrite bool,
	value string,
	caName string,
	root string,
	certificate string,
	private string,
	public string,
	rootString string,
	certificateString string,
	privateString string,
	publicString string,
	username string,
	password string,
) (models.Printable, error) {

	cfg := config.ReadConfig()
	repository := repositories.NewCredentialRepository(client.NewHttpClient(cfg))

	action := actions.NewAction(repository, &cfg)
	var request *http.Request
	if credentialType == "ssh" || credentialType == "rsa" {
		var err error

		err = setStringFieldFromFile(&public, &publicString)
		if err != nil {
			return nil, err
		}

		err = setStringFieldFromFile(&private, &privateString)
		if err != nil {
			return nil, err
		}

		request = client.NewSetRsaSshRequest(cfg, credentialIdentifier, credentialType, publicString, privateString, !noOverwrite)
	} else if credentialType == "certificate" {
		var err error

		err = setStringFieldFromFile(&root, &rootString)
		if err != nil {
			return nil, err
		}

		err = setStringFieldFromFile(&certificate, &certificateString)
		if err != nil {
			return nil, err
		}

		err = setStringFieldFromFile(&private, &privateString)
		if err != nil {
			return nil, err
		}

		request = client.NewSetCertificateRequest(cfg, credentialIdentifier, rootString, caName, certificateString, privateString, !noOverwrite)
	} else if credentialType == "user" {
		request = client.NewSetUserRequest(cfg, credentialIdentifier, username, password, !noOverwrite)
	} else if credentialType == "password" {
		request = client.NewSetCredentialRequest(cfg, credentialType, credentialIdentifier, password, !noOverwrite)
	} else if credentialType == "json" {
		request = client.NewSetJsonCredentialRequest(cfg, credentialType, credentialIdentifier, value, !noOverwrite)
	} else {
		request = client.NewSetCredentialRequest(cfg, credentialType, credentialIdentifier, value, !noOverwrite)
	}

	credential, err := action.DoAction(request, credentialIdentifier)

	return credential, err
}

func setStringFieldFromFile(fileField, stringField *string) error {
	var err error
	if *fileField != "" {
		if *stringField != "" {
			return errors.NewCombinationOfParametersError()
		}
		*stringField, err = readFile(*fileField)
		if err != nil {
			return err
		}
	}
	return nil
}

func readFile(filename string) (string, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.NewFileLoadError()
	}
	return string(dat), nil
}
