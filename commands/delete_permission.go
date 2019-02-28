package commands

import (
	"fmt"
)

type DeletePermissionCommand struct {
	Actor string `short:"a" long:"actor" required:"yes" description:"Name of the actor to grant permissions for"`
	Path  string `short:"p" long:"path" required:"yes" description:"Name of path to grant permissions for"`
	ClientCommand
}

func (c *DeletePermissionCommand) DeletePermission(uuid string) error {

	permission, err := c.client.DeletePermission(uuid)
	if err != nil {
		return err
	}

	formatOutput(false, permission)
	return nil
}

func (c *DeletePermissionCommand) Execute([]string) error {
	serverVersion, _ := c.client.ServerVersion()
	isOlderVersion := serverVersion.Segments()[0] < 2
	if isOlderVersion {
		return fmt.Errorf("credhub server version <2.0 not supported")
	}

	permission, err := c.client.GetPermissionByPathActor(c.Path, c.Actor)
	if err != nil {
		return err
	}
	return c.DeletePermission(permission.UUID)
}
