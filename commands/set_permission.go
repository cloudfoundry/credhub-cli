package commands

type SetPermissionCommand struct {
	Actor      string `short:"a" long:"actor" required:"yes" description:"Name of the actor to grant permissions for"`
	Path       string `short:"p" long:"path" required:"yes" description:"Name of path to grant permissions for"`
	Operations string `short:"o" long:"operations" required:"yes" description:"Operations to actor is granted permissions for one path"`
	OutputJSON bool   `short:"j" long:"output-json" description:"Return response in JSON format"`
	ClientCommand
}

func (c *SetPermissionCommand) Execute([]string) error {
	return nil
}
