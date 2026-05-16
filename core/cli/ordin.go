// Package cli is meant to build the flags, cli, commands, etc
package cli

import (
	"context"
	"fmt"
	"strings"

	"ordin/m/core/cli/dryrun"
	"ordin/m/core/ruleengine"
	"ordin/m/core/ruleengine/executor"

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
					workingDest := strings.TrimSpace(cmd.Args().Get(1))
					if rulesDest == "" {
						return fmt.Errorf("invalid first argument")
					}
					if workingDest == "" {
						return fmt.Errorf("invalid second argument")
					}

					sortedRules, err := ruleengine.CheckSort(rulesDest)
					if err != nil {
						return err
					}
					paths, err := ruleengine.Plan(sortedRules, workingDest)
					if err != nil {
						return err
					}
					finalPaths := dryrun.DryRunTUIInit(paths)
					err = dryrun.SavePlan("plan.json", finalPaths)
					if err != nil {
						return err
					}
					return nil
				},
			},

			{
				Name:    "run",
				Usage:   "to run",
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

					sortedRules, err := ruleengine.CheckSort(rulesDest)
					if err != nil {
						return err
					}
					paths, err := ruleengine.Plan(sortedRules, workingDest)
					if err != nil {
						return err
					}
					finalPaths := dryrun.DryRunTUIInit(paths)
					finalStats, err := executor.Executor(finalPaths)
					if err != nil {
						return err
					}
					fmt.Printf("Successful Operations: %v\nFailed Operations: %v\nRun Complete!", finalStats.Successful.Load(), finalStats.Error.Load())
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
