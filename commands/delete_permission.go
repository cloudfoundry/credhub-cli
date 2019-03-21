package commands

type DeletePermissionCommand struct {
	Actor string `short:"a" long:"actor" required:"yes" description:"Name of the actor to grant permissions for"`
	Path  string `short:"p" long:"path" required:"yes" description:"Name of path to grant permissions for"`
	ClientCommand
}

func (c *DeletePermissionCommand) Execute([]string) error {
	return nil
}
