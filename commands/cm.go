package commands

type CMCommand struct {
	Api ApiCommand `command:"api" alias:"a" description:"Set the API server to use"`

	Get        GetCommand        `command:"get" alias:"g" description:"Get a secret value"`
	Set        SetCommand        `command:"set" alias:"s" description:"Set a secret value"`
	Generate   GenerateCommand   `command:"generate" alias:"n" description:"Generate a secret value"`
	Regenerate RegenerateCommand `command:"regenerate" description:"Regenerate a secret value"`
	Delete     DeleteCommand     `command:"delete" alias:"d" description:"Delete a secret value"`
	CaSet      CaSetCommand      `command:"ca-set" alias:"cs" description:"Set a certificate authority for generating signed certificates"`
	CaGet      CaGetCommand      `command:"ca-get" alias:"cg" description:"Get a certificate authority"`
	CaGenerate CaGenerateCommand `command:"ca-generate" alias:"cn" description:"Configures a certificate authority with a generated key pair."`
	Login      LoginCommand      `command:"login" alias:"l" description:"Authenticates user with CredHub"`
	Logout     LogoutCommand     `command:"logout" alias:"o" description:"Discard authenticated user session."`
	Find       FindCommand       `command:"find" alias:"f" description:"Find existing credentials with query parameters"`
	Version    func()            `long:"version" description:"Version of Credential Manager"`
}

var CM CMCommand
