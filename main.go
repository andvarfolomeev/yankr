package main

import (
	"fmt"
	"os"

	"github.com/andvarfolomeev/yankr/cmd"
	"github.com/andvarfolomeev/yankr/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	app := cmd.BuildApp(cfg)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
