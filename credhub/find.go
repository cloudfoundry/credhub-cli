package credhub

import "github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"

func (ch CredHub) FindByPartialName(nameLike string) ([]credentials.Base, error) {
	panic("Not implemented")
}

func (ch CredHub) FindByPath(path string) ([]credentials.Base, error) {
	panic("Not implemented")
}

// FIXME Should Path be in credentials?
func (ch CredHub) ShowAllPaths() ([]credentials.Path, error) {
	panic("Not implemented")
}
