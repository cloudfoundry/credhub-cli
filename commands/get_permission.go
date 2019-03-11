package commands

import (
	"fmt"
)

type GetPermissionCommand struct {
	Actor      string `short:"a" long:"actor" required:"yes" description:"Name of the actor to grant permissions for"`
	Path       string `short:"p" long:"path" required:"yes" description:"Name of path to grant permissions for"`
	OutputJSON bool   `short:"j" long:"output-json" description:"Return response in JSON format"`
	ClientCommand
}

func (c *GetPermissionCommand) GetPermission() error {
	permission, err := c.client.GetPermissionByPathActor(c.Path, c.Actor)
	if err != nil {
		return err
	}

	formatOutput(c.OutputJSON, permission)
	return err
}

func (c *GetPermissionCommand) Execute([]string) error {
	serverVersion, _ := c.client.ServerVersion()
	isOlderVersion := serverVersion.Segments()[0] < 2
	if isOlderVersion {
		return fmt.Errorf("credhub server version <2.0 not supported")
	}
	return c.GetPermission()

}
