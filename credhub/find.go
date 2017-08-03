package credhub

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"

// Retrieves a list of stored credential names which contain the search.
func (ch *CredHub) FindByPartialName(nameLike string) ([]credentials.Base, error) {
	panic("Not implemented")
}

// Retrieves a list of stored credential names which are within the specified path.
func (ch *CredHub) FindByPath(path string) ([]credentials.Base, error) {
	panic("Not implemented")
}

// Retrieves a list of all paths which contain credentials.
func (ch *CredHub) ShowAllPaths() ([]credentials.Path, error) {
	panic("Not implemented")
}
