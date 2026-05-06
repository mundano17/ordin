package main

import (
	"context"
	"log"
	"ordin/m/core/cli"
	"os"
)

func main() {
	cmd := cli.CliRunner()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
