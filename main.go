package main

import (
	"os"

	"github.com/zeynepyz/k6beat/cmd"

	// Make sure all your modules and metricsets are linked in this file
	_ "github.com/zeynepyz/k6beat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
