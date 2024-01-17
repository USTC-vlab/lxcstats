package main

import (
	"log"

	"github.com/USTC-vlab/vct/cmd"
)

func main() {
	if err := cmd.MakeCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
