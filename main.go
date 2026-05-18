package main

import (
	"context"
	"log"
	"ordin/cmd"
	"os"
)

func main() {
	cmd := cmd.CliRunner()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
