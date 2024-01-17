package cmd

import (
	"github.com/USTC-vlab/vct/cmd/df"
	"github.com/USTC-vlab/vct/cmd/findpid"
	"github.com/USTC-vlab/vct/cmd/iostat"
	"github.com/USTC-vlab/vct/cmd/pressure"
	"github.com/spf13/cobra"
)

var Version string

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vct",
		Short: "Vlab's Container Tool",
		Long:  "Vlab's Container Tool, a versatile tool for managing containers on Proxmox VE",
		Args:  cobra.NoArgs,
	}
	cmd.CompletionOptions.HiddenDefaultCmd = true
	cmd.AddCommand(
		df.MakeCmd(),
		findpid.MakeCmd(),
		iostat.MakeCmd(),
		pressure.MakeCmd(),
	)
	pVersion := cmd.Flags().BoolP("version", "v", false, "show version")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if *pVersion {
			cmd.Println(cmd.Name(), Version)
		} else {
			cmd.Help()
		}
	}
	return cmd
}
