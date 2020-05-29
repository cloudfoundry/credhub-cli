package commands

import (
	"code.cloudfoundry.org/credhub-cli/errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type DeleteCommand struct {
	CredentialIdentifier string `short:"n" long:"name" description:"Name of the credential to delete"`
	CredentialPath       string `short:"p" long:"path" description:"Name of the credentials to delete"`
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
	failedCredentials, credentialsCount, err := c.deleteByPath(c.CredentialPath)
	if err != nil {
		return err
	}

	failedCredentialsCount := len(failedCredentials)
	if failedCredentialsCount == 0 {
		fmt.Printf("\nAll %v out of %v credentials under the provided path are successfully deleted.", credentialsCount, credentialsCount)
		return nil
	}

	failureMessage := fmt.Sprintf("%v out of %v credentials under the provided path are successfully deleted. The following credentials failed to delete:", failedCredentialsCount, credentialsCount)
	fmt.Fprintln(os.Stderr, failureMessage)

	s, _ := yaml.Marshal(failedCredentials)
	fmt.Fprint(os.Stderr, string(s))

	return errors.NewBulkDeleteFailureError()
}

type DeleteFailedCredential struct {
	Name string
	Err  string
}

func (c *DeleteCommand) deleteByPath(path string) ([]DeleteFailedCredential, int, error) {
	results, err := c.client.FindByPath(path)
	if err != nil {
		return []DeleteFailedCredential{}, 0, err
	}

	var failedCredentials []DeleteFailedCredential
	for _, cred := range results.Credentials {
		err = c.client.Delete(cred.Name)

		if err != nil {
			failedCredentials = append(failedCredentials, DeleteFailedCredential{
				cred.Name,
				err.Error(),
			})
		}
	}
	return failedCredentials, len(results.Credentials), nil
}
