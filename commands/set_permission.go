package commands

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SetPermissionCommand struct {
	Actor      string `short:"a" long:"actor" required:"yes" description:"Name of the actor to grant permissions for"`
	Path       string `short:"p" long:"path" required:"yes" description:"Name of path to grant permissions for"`
	Operations string `short:"o" long:"operations" required:"yes" description:"Operations to actor is granted permissions for one path"`
	ClientCommand
}

func ParseOperations(operations string) []string {
	ops := strings.Split(operations, ",")
	trimmedOps := make([]string, len(ops))
	for i, v := range ops {
		trimmedOps[i] = strings.TrimSpace(v)
	}
	return trimmedOps
}

func (c *SetPermissionCommand) addPermission() error {
	ops := ParseOperations(c.Operations)
	permission, err := c.client.AddPermission(c.Path, c.Actor, ops)

	if err != nil {
		return err
	}
	permissionJson, _ := json.Marshal(permission)
	fmt.Println(string(permissionJson))
	return nil
}

func (c *SetPermissionCommand) updatePermission(uuid string) error {
	ops := ParseOperations(c.Operations)
	permission, err := c.client.UpdatePermission(uuid, c.Path, c.Actor, ops)
	if err != nil {
		return err
	}
	permissionJson, _ := json.Marshal(permission)
	fmt.Println(string(permissionJson))
	return nil
}

func (c *SetPermissionCommand) Execute([]string) error {
	serverVersion, _ := c.client.ServerVersion()
	isOlderVersion := serverVersion.Segments()[0] < 2
	if isOlderVersion {
		return fmt.Errorf("credhub server version <2.0 not supported")
	}

	permission, err := c.client.GetPermissionByPathActor(c.Path, c.Actor)
	if err != nil {
		return err
	}

	if permission.UUID == "" {
		return c.addPermission()
	} else {
		return c.updatePermission(permission.UUID)
	}

}
