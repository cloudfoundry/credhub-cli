package test

import (
	"io/ioutil"
	"os"
)

func CleanEnv() {
	os.Unsetenv("CREDHUB_SECRET")
	os.Unsetenv("CREDHUB_CLIENT")
}

func CreateTempDir(prefix string) string {
	name, err := ioutil.TempDir("", prefix)
	if err != nil {
		panic(err)
	}
	return name
}

func CreateCredentialFile(dir, filename string, contents string, perm os.FileMode) string {
	path := dir + "/" + filename
	err := ioutil.WriteFile(path, []byte(contents), perm)
	if err != nil {
		panic(err)
	}
	return path
}
