package main

import (
	"os"

	"github.com/getopends/opends/pkg/cmd"
)

func main() {
	cmd := cmd.RootCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
