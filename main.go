package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

var cmd = &cli.Command{
	Name:                  "gcp-iam",
	Usage:                 "Query Google Cloud IAM Roles and Permissions",
	Suggest:               true,
	EnableShellCompletion: true,
	HideHelpCommand:       true,
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return cli.ShowAppHelp(cmd)
	},
	Commands: []*cli.Command{
		{
			Name:  "role",
			Usage: "Query IAM Roles",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return cli.ShowAppHelp(cmd)
			},
			Commands: []*cli.Command{
				{
					Name:  "show",
					Usage: "Show IAM role permissions",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println(cmd.FullName(), cmd.Args().First())
						return nil
					},
				},
				{
					Name:  "search",
					Usage: "Search IAM roles",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println(cmd.FullName(), cmd.Args().First())
						return nil
					},
				},
				{
					Name:  "compare",
					Usage: "Compare permissions of 2 IAM roles",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println(cmd.FullName(), cmd.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:  "permission",
			Usage: "Query IAM Permissions",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return cli.ShowAppHelp(cmd)
			},
			Commands: []*cli.Command{
				{
					Name:  "show",
					Usage: "Show IAM roles with permission",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println(cmd.FullName(), cmd.Args().First())
						return nil
					},
				},
				{
					Name:  "search",
					Usage: "Search IAM permissions",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						fmt.Println(cmd.FullName(), cmd.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:  "update",
			Usage: "Update IAM roles and permissions",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println(cmd.FullName(), cmd.Args().First())
				return nil
			},
		},
		{
			Name:  "info",
			Usage: "Show application configuration",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println(cmd.FullName(), cmd.Args().First())
				return nil
			},
		},
	},
}

func main() {
	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatalf("failed to run command: %v", err)
	}
}
