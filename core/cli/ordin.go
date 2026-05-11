package cli

import (
	"context"
	"fmt"
	"ordin/m/core/cli/dryrun"
	"ordin/m/core/rule_engine"
	"strings"

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
						_, err := rule_engine.CheckPipeline(val)
						if err != nil {
							return err
						} else {
							fmt.Println("Good to go!")
							return nil
						}
					}
					return fmt.Errorf("Invalid argument")
				},
			},

			{
				Name:    "dryrun",
				Usage:   "to run and get a log",
				Aliases: []string{"drun"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rules_dest := strings.TrimSpace(cmd.Args().Get(0))
					working_dest := strings.TrimSpace(cmd.Args().Get(1))
					if rules_dest == "" {
						return fmt.Errorf("Invalid first argument")
					}
					if working_dest == "" {
						return fmt.Errorf("Invalid second argument")
					}

					sortedRules, err := rule_engine.CheckSort(rules_dest)
					if err != nil {
						return err
					}
					paths, err := rule_engine.Plan(sortedRules, working_dest)
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
