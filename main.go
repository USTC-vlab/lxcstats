package main

import (
	"os"

	"github.com/USTC-vlab/vct/cmd"
)

var version string = "<unknown>"

func main() {
	cmd.Version = version
	if err := cmd.MakeCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
