package commands

import (
	"fmt"
	"strings"
)

type CreatePermissionCommand struct {
	Actor      string `short:"a" long:"actor" description:"Name of the actor to grant permissions for"`
	Path       string `short:"p" long:"path" description:"Name of path to grant permissions for"`
	Operations string `short:"o" long:"operations" description:"Operations to actor is granted permissions for on path"`
	ClientCommand
}

func (c *CreatePermissionCommand) Execute([]string) error {
	ops := strings.Split(c.Operations, ",")
	trimmedOps := make([]string, len(ops))
	for i, v := range ops {
		trimmedOps[i] = strings.TrimSpace(v)
	}
	permission, _ := c.client.AddPermission(c.Path, c.Actor, trimmedOps)
	fmt.Println(permission.Actor)

	return nil
}
