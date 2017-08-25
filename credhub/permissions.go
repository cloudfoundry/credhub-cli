package credhub

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/permissions"

// Returns the permissions of a credential.
func (ch *CredHub) GetPermissions(credName string) ([]permissions.Permission, error) {
	panic("Not implemented")
}

// Adds permissions to a credential.
func (ch *CredHub) AddPermissions(credName string, perms []permissions.Permission) ([]permissions.Permission, error) {
	panic("Not implemented")
}

// Deletes a permission for an actor on a credential.
func (ch *CredHub) DeletePermissions(credName string, actor string) error {
	panic("Not implemented")
}
