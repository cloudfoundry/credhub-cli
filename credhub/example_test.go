package credhub_test

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/generator"
)

func Example_lookup() {
	ch := credhub.CredHub{}
	path := "/some/path"
	name := "key"

	creds, err := ch.FindByPath(path)
	if err != nil {
		panic(err)
	}

	found := false
	for _, cred := range creds {
		if cred.Name == (path + name) {
			found = true
		}
	}
	if !found {
		panic("key not found")
	}

	cred, err := ch.GetByName(path + name)
	if err != nil {
		panic(err)
	}

	fmt.Println(cred.Value)
}

func Example_create() {
	ch := credhub.CredHub{}

	name := "/some/path/to/cert"

	gen := generator.Certificate{
		CommonName: "pivotal",
		KeyLength:  2048,
	}

	cert, err := ch.GenerateCertificate(name, gen, false)
	if err != nil {
		panic(err)
	}

	dupCert, err := ch.SetCertificate("/some/path/to/dup-cert", cert.Value, false)
	if err != nil {
		panic(err)
	}

	if dupCert.Value.Certificate != cert.Value.Certificate {
		panic("certs don't match")
	}
}

func Example_generate() {
	ch := credhub.CredHub{}

	username := "some-user"
	path := "/some/path/" + username

	user, err := ch.GetUserByName(username)
	if err != nil {
		ch.Delete(path)

		user, err = ch.GenerateUser(path, generator.User{Username: username}, false)
		if err != nil {
			panic(err)
		}
	} else {
		user, err = ch.RegenerateUser(path)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Password: ", user.Value.Password)
}
