package main

import (
	"os"

	"github.com/USTC-vlab/vct/cmd"
)

func main() {
	if err := cmd.MakeCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
