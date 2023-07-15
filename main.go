package main

import (
	"github.com/simonwep/genesis/commands"
	"github.com/simonwep/genesis/core"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "genesis",
		Usage: "A tinsy tiny server for all your json blob needs",
		Authors: []*cli.Author{
			{Name: "Simon Reinisch", Email: "simon@reinisch.io"},
		},
		Commands: []*cli.Command{
			{
				Name:   "start",
				Usage:  "Start the api",
				Action: commands.Start,
			},
			{
				Name:  "user",
				Usage: "Manage users",
				Subcommands: []*cli.Command{
					{
						Name:      "ls",
						Usage:     "Lists all users",
						UsageText: "genesis user ls",
						Action:    commands.ListUsers,
					},
					{
						Name:      "rm",
						Usage:     "Removes a user",
						UsageText: "genesis user rm [username]",
						Action:    commands.RemoveUser,
					},
					{
						Name:      "add",
						Usage:     "Adds a user",
						UsageText: "genesis user add [username] [password]",
						Action:    commands.AddUser,
					},
					{
						Name:      "update",
						Usage:     "Updates a user",
						UsageText: "genesis user update [options] [username]",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "password",
								Usage: "Sets a new password",
							},
							&cli.StringFlag{
								Name:  "name",
								Usage: "Changes the name",
							},
						},
						Action: commands.UpdateUser,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		core.Logger.Fatal("failed to start cli", zap.Error(err))
	}
}
