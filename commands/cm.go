package commands

type CMCommand struct {
	ApiURL string `long:"api" description:"Credential Manager API URL" default:"https://pivotal-credential-manager.cfapps.io/"`

	Get GetCommand `command:"get" description:"Get a secret value"`
	Set SetCommand `command:"set" description:"Set a secret value"`
}

var CM CMCommand
