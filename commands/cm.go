package commands

type CMCommand struct {
	Api ApiCommand `command:"api" description:"Set the API server to use"`

	Get     GetCommand     `command:"get" description:"Get a secret value"`
	Set     SetCommand     `command:"set" description:"Set a secret value"`
	Delete  DeleteCommand  `command:"delete" description:"Delete a secret value"`
	Version VersionCommand `command:"version" description:"Display the CLI version"`
}

var CM CMCommand
