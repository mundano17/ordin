// Package cli is meant to build the flags, cli, commands, etc
package cmd

import (
	"context"
	"fmt"
	"strings"

	"ordin/rules"
	"ordin/tui/dryrun"

	"github.com/urfave/cli/v3"
)

func CliRunner() *cli.Command {
	cmd := &cli.Command{
		Name:    "ordin",
		Usage:   "sort your filez",
		Version: "0.1",

		Commands: []*cli.Command{
			{
				Name:    "check",
				Usage:   "to validate rules config",
				Aliases: []string{"c"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					val := strings.TrimSpace(cmd.Args().First())
					if val != "" {
						err := rules.ValidateFile(val)
						if err != nil {
							return err
						}

						fmt.Println("Good to go!")
						return nil
					}
					return fmt.Errorf("invalid argument")
				},
			},

			{
				Name:    "dryrun",
				Usage:   "to run and get a log",
				Aliases: []string{"drun"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rulesDest := strings.TrimSpace(cmd.Args().Get(0))
					workingDest := strings.TrimSpace(cmd.Args().Get(1))
					if rulesDest == "" {
						return fmt.Errorf("invalid first argument")
					}
					if workingDest == "" {
						return fmt.Errorf("invalid second argument")
					}

					sortedRules, err := rules.ValidateAndSort(rulesDest)
					if err != nil {
						return err
					}
					paths, err := rules.Planner(sortedRules, workingDest)
					if err != nil {
						return err
					}
					dryrun.InitializeDryRun(paths)

					if err != nil {
						return err
					}
					return nil
				},
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Printf("Ordin Version: %s\n", cmd.Version)
			return nil
		},
	}
	return cmd
}
