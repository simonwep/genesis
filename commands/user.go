package commands

import (
	"fmt"
	"github.com/simonwep/genisis/core"
	"github.com/urfave/cli/v2"
)

func ListUsers(ctx *cli.Context) error {
	fmt.Printf("//%v//", ctx.Args().Get(0))
	return nil
}

func RemoveUser(ctx *cli.Context) error {
	return core.DeleteUser(ctx.Args().Get(0))
}

func AddUser(ctx *cli.Context) error {
	fmt.Printf("//%v//", ctx.Args().Get(0))
	return nil
}

func UpdateUser(ctx *cli.Context) error {
	fmt.Printf("//%v//", ctx.Args().Get(0))
	return nil
}
