package commands

type GetPermissionCommand struct {
	Actor      string `short:"a" long:"actor" required:"yes" description:"Name of the actor to grant permissions for"`
	Path       string `short:"p" long:"path" required:"yes" description:"Name of path to grant permissions for"`
	OutputJSON bool   `short:"j" long:"output-json" description:"Return response in JSON format"`
	ClientCommand
}

func (c *GetPermissionCommand) Execute([]string) error {
	return nil
}
