// Package cli is meant to build the flags, cli, commands, etc
package cli

import (
	"context"
	"fmt"
	"strings"

	"ordin/m/core/cli/dryrun"
	"ordin/m/core/ruleengine"

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
						_, err := ruleengine.CheckPipeline(val)
						if err != nil {
							return err
						} else {
							fmt.Println("Good to go!")
							return nil
						}
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
					worlingDest := strings.TrimSpace(cmd.Args().Get(1))
					if rulesDest == "" {
						return fmt.Errorf("invalid first argument")
					}
					if worlingDest == "" {
						return fmt.Errorf("invalid second argument")
					}

					sortedRules, err := ruleengine.CheckSort(rulesDest)
					if err != nil {
						return err
					}
					paths, err := ruleengine.Plan(sortedRules, worlingDest)
					if err != nil {
						return err
					}
					dryrun.DryRunTUIInit(paths)
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
