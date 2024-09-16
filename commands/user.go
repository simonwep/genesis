package commands

import (
	"fmt"
	"github.com/simonwep/genesis/core"
	"github.com/urfave/cli/v2"
	"strings"
)

func ListUsers(_ *cli.Context) error {
	if users, err := core.GetAllUsers(); err != nil {
		return err
	} else {

		// Print users
		for _, user := range users {
			fmt.Printf("Name: %v, Admin: %v\n", user.Name, user.Admin)
		}
	}

	return nil
}

func RemoveUser(ctx *cli.Context) error {
	return core.DeleteUser(ctx.Args().Get(0))
}

func AddUser(ctx *cli.Context) error {
	username, password := ctx.Args().Get(0), ctx.Args().Get(1)
	admin := strings.HasSuffix(username, "!")

	return core.CreateUser(core.User{
		Name:     strings.TrimSuffix(username, "!"),
		Admin:    admin,
		Password: password,
	})
}

func UpdateUser(ctx *cli.Context) error {
	fmt.Printf("//%v//", ctx.Args().Get(0))
	return nil
}
