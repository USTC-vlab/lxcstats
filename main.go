package main

import (
	"fmt"
	"os"

	"github.com/USTC-vlab/vct/cmd"
)

func main() {
	if err := cmd.MakeCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
