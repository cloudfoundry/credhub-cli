package commands

type CMCommand struct {
	Api ApiCommand `command:"api" alias:"a" description:"Set the CredHub API target to be used for subsequent commands"`

	Get        GetCommand        `command:"get" alias:"g" description:"Get a credential value"`
	Set        SetCommand        `command:"set" alias:"s" description:"Set a credential with a provided value"`
	Generate   GenerateCommand   `command:"generate" alias:"n" description:"Set a credential with a generated value"`
	Regenerate RegenerateCommand `command:"regenerate" alias:"r" description:"Set a credential with a generated value using the same attributes as the stored value"`
	Delete     DeleteCommand     `command:"delete" alias:"d" description:"Delete a credential value"`
	Login      LoginCommand      `command:"login" alias:"l" description:"Authenticate user with CredHub"`
	Logout     LogoutCommand     `command:"logout" alias:"o" description:"Discard authenticated user session"`
	Find       FindCommand       `command:"find" alias:"f" description:"Find stored credentials based on query parameters"`
	Version    func()            `long:"version" description:"Version of CLI and targeted CredHub API"`
}

var CM CMCommand
