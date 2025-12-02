package commands

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/credhub-cli/errors"
	"go.yaml.in/yaml/v3"
)

type DeleteCommand struct {
	CredentialIdentifier string `short:"n" long:"name" description:"Name of the credential to delete"`
	CredentialPath       string `short:"p" long:"path" description:"Path of the credentials to delete"`
	Quiet                bool   `short:"q" long:"quiet" description:"Disable real-time status of delete by path"`
	ClientCommand
}

func (c *DeleteCommand) Execute([]string) error {
	if c.CredentialIdentifier != "" {
		return c.handleDeleteByName()
	} else if c.CredentialPath != "" {
		return c.handleDeleteByPath()
	}

	return errors.NewMissingDeleteParametersError()
}

func (c *DeleteCommand) handleDeleteByName() error {
	err := c.client.Delete(c.CredentialIdentifier)

	if err == nil {
		fmt.Println("Credential successfully deleted")
	}

	return err
}

func (c *DeleteCommand) handleDeleteByPath() error {
	failedCredentials, credentialsCount, err := c.deleteByPath(c.CredentialPath, c.Quiet)
	if err != nil {
		return err
	}

	failedCredentialsCount := len(failedCredentials)
	if failedCredentialsCount == 0 {
		if c.Quiet {
			fmt.Printf("All %v out of %v credentials under the provided path are successfully deleted.\n", credentialsCount, credentialsCount)
		}
		return nil
	}

	if c.Quiet {
		fmt.Printf("%v out of %v credentials under the provided path are successfully deleted.\n", credentialsCount-failedCredentialsCount, credentialsCount)
	}

	failureMessage := fmt.Sprintf("%v out of %v credentials under the provided path failed to delete. The following credentials failed to delete:", failedCredentialsCount, credentialsCount)
	fmt.Fprintln(os.Stderr, failureMessage)

	s, _ := yaml.Marshal(failedCredentials)
	fmt.Fprint(os.Stderr, string(s))

	return errors.NewBulkDeleteFailureError()
}

type DeleteFailedCredential struct {
	Name string
	Err  string
}

func (c *DeleteCommand) deleteByPath(path string, quiet bool) ([]DeleteFailedCredential, int, error) {
	results, err := c.client.FindByPath(path)
	if err != nil {
		return []DeleteFailedCredential{}, 0, err
	}

	var totalCount = len(results.Credentials)
	var failedCredentials []DeleteFailedCredential
	for index, cred := range results.Credentials {
		err = c.client.Delete(cred.Name)

		if err != nil {
			failedCredentials = append(failedCredentials, DeleteFailedCredential{
				cred.Name,
				err.Error(),
			})
		}

		if !quiet {
			succeeded := index + 1 - len(failedCredentials)
			fmt.Printf("\033[2K\r%v out of %v credentials under the provided path are successfully deleted.\n", succeeded, totalCount)
		}
	}
	return failedCredentials, totalCount, nil
}
