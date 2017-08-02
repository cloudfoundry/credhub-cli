package credhub

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/permissions"

func (ch CredHub) GetPermissions(credName string) ([]permissions.Permission, error) {
	panic("Not implemented")
}

func (ch CredHub) AddPermissions(credName string, perms []permissions.Permission) ([]permissions.Permission, error) {
	panic("Not implemented")
}

func (ch CredHub) DeletePermissions(credName string, actor string) error {
	panic("Not implemented")
}
