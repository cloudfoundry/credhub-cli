package commands

type CMCommand struct {
	ApiURL string `long:"api" description:"The Credential Manager API Location"`

	SetSecret SetSecretCommand `command:"set-secret" alias:"ss" description:"Sets secret value"`
}

var CM CMCommand