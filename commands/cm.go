package commands

type CMCommand struct {
	Api ApiCommand `command:"api" description:"Set the API server to use"`

	Get      GetCommand      `command:"get" description:"Get a secret value"`
	Set      SetCommand      `command:"set" description:"Set a secret value"`
	Generate GenerateCommand `command:"generate" description:"Generate a secret value"`
	Delete   DeleteCommand   `command:"delete" description:"Delete a secret value"`
	CaSet    CaSetCommand    `command:"ca-set" description:"Set a root certificate"`
	Version  func()          `short:"v" long:"version" description:"Version of Credential Manager"`
}

var CM CMCommand
