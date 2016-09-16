package commands

type CMCommand struct {
	Api ApiCommand `command:"api" description:"Set the API server to use"`

	Get        GetCommand        `command:"get" description:"Get a secret value"`
	Set        SetCommand        `command:"set" description:"Set a secret value"`
	Generate   GenerateCommand   `command:"generate" description:"Generate a secret value"`
	Delete     DeleteCommand     `command:"delete" description:"Delete a secret value"`
	CaSet      CaSetCommand      `command:"ca-set" description:"Set a certificate authority for generating signed certificates"`
	CaGet      CaGetCommand      `command:"ca-get" description:"Get a certificate authority"`
	CaGenerate CaGenerateCommand `command:"ca-generate" description:"Configures a certificate authority with a generated key pair."`
	Login      LoginCommand      `command:"login" description:"Authenticates user with CredHub"`
	Logout     LogoutCommand     `command:"logout" description:"Discard authenticated user session."`
	Find       FindCommand       `command:"find" description:"Find existing credentials with query parameters"`
	Version    func()            `long:"version" description:"Version of Credential Manager"`
}

var CM CMCommand
