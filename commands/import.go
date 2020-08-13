package commands

import (
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"fmt"
	"strconv"

	"os"

	"reflect"

	"code.cloudfoundry.org/credhub-cli/errors"
	"code.cloudfoundry.org/credhub-cli/models"
)

type ImportCommand struct {
	File       string `short:"f" long:"file" description:"File containing credentials to import" required:"true"`
	ImportJSON bool   `short:"j" long:"import-json" description:"File to import is of type JSON"`
	ClientCommand
}

type CaAndIndex struct {
	Ca    string
	Index int
}

type ErrorInfo struct {
	Successful   int
	Failed       int
	ImportErrors []string
}

func (c *ImportCommand) Execute([]string) error {
	var bulkImport models.CredentialBulkImport
	err := bulkImport.ReadFile(c.File, c.ImportJSON)

	if err != nil {
		return err
	}

	err = c.setCredentials(bulkImport)

	return err
}

func (c *ImportCommand) setCredentials(bulkImport models.CredentialBulkImport) error {
	var (
		name      string
		errorInfo = ErrorInfo{}
	)
	certsWithCaName := make(map[string]CaAndIndex)

	for i, credential := range bulkImport.Credentials {
		switch credentialName := credential["name"].(type) {
		case string:
			name = credentialName
		default:
			name = ""
		}

		var certWithCaName bool
		var caName string
		switch credential["type"].(string) {
		case "ssh":
			if _, ok := credential["value"].(map[string]interface{})["public_key_fingerprint"]; ok {
				delete(credential["value"].(map[string]interface{}), "public_key_fingerprint")
			}
		case "user":
			if _, ok := credential["value"].(map[string]interface{})["password_hash"]; ok {
				delete(credential["value"].(map[string]interface{}), "password_hash")
			}
		case "value":
			switch value := credential["value"].(type) {
			case int:
				credential["value"] = strconv.Itoa(value)
			case float32:
				credential["value"] = strconv.FormatFloat(float64(value), 'f', -1, 32)
			case float64:
				credential["value"] = strconv.FormatFloat(value, 'f', -1, 64)
			}
		case "certificate":
			caName, certWithCaName = credential["value"].(map[string]interface{})["ca_name"].(string)
		}

		if certWithCaName {
			certsWithCaName[name] = CaAndIndex{caName, i}
		} else {
			err := c.setCredentialInCredHub(
				name, credential["type"].(string), credential["value"], credential["metadata"], &errorInfo, i)
			if err != nil {
				return err
			}
		}
	}

	for signedCert := range certsWithCaName {
		err := c.importCert(signedCert, certsWithCaName, bulkImport.Credentials, &errorInfo)
		if err != nil {
			return err
		}
	}

	fmt.Println("Import complete.")
	_, _ = fmt.Fprintf(os.Stdout, "Successfully set: %d\n", errorInfo.Successful)
	_, _ = fmt.Fprintf(os.Stdout, "Failed to set: %d\n", errorInfo.Failed)
	for _, v := range errorInfo.ImportErrors {
		fmt.Println(v)
	}

	if errorInfo.Failed > 0 {
		return errors.NewFailedToImportError()
	}

	return nil
}

func isAuthenticationError(err error) bool {
	return reflect.DeepEqual(err, errors.NewNoApiUrlSetError()) ||
		reflect.DeepEqual(err, errors.NewUnauthenticatedError()) ||
		reflect.DeepEqual(err, errors.NewRefreshError())
}

func (c *ImportCommand) setCredentialInCredHub(name, credType string, value, metadata interface{}, errorInfo *ErrorInfo, index int) error {
	var options []credhub.SetOption

	if metadata != nil {
		var meta credentials.Metadata
		meta = metadata.(map[string]interface{})
		withMetadata := func(s *credhub.SetOptions) error {
			s.Metadata = meta
			return nil
		}

		options = append(options, withMetadata)
	}

	_, err := c.client.SetCredential(name, credType, value, options...)

	if err != nil {
		if isAuthenticationError(err) {
			return err
		}
		failure := fmt.Sprintf("Credential '%s' at index %d could not be set: %v", name, index, err)
		fmt.Println(failure + "\n")
		errorInfo.ImportErrors = append(errorInfo.ImportErrors, " - "+failure)
		errorInfo.Failed++
	} else {
		errorInfo.Successful++
	}
	return nil
}

func (c *ImportCommand) importCert(cert string, certs map[string]CaAndIndex, credentials []map[string]interface{}, errorInfo *ErrorInfo) error {
	caAndIndex, certNotImported := certs[cert]
	if !certNotImported {
		return nil
	}
	_, caNotImported := certs[caAndIndex.Ca]
	if caNotImported {
		err := c.importCert(caAndIndex.Ca, certs, credentials, errorInfo)
		if err != nil {
			return err
		}
	}
	delete(certs, cert)
	credential := credentials[caAndIndex.Index]
	return c.setCredentialInCredHub(cert, credential["type"].(string), credential["value"], credential["metadata"], errorInfo, caAndIndex.Index)
}
