package main

import (
	"os"

	"github.com/uesugitorachiyo/ao-blueprint/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		os.Exit(1)
	}
}
