package main

import (
	"os"

	"github.com/gocastsian/roham/cmd/filer/command"
)

func main() {
	if err := command.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
