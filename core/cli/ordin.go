package main

import (
	"context"
	"fmt"
	"log"
	"ordin/m/core/rule_engine"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

func main() {

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
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Printf("Ordin Version: %s\n", cmd.Version)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
