package main

import (
	"os"

	"github.com/nycon/qcli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
