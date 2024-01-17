package cmd

import (
	"github.com/USTC-vlab/vct/cmd/df"
	"github.com/USTC-vlab/vct/cmd/findpid"
	"github.com/USTC-vlab/vct/cmd/iostat"
	"github.com/USTC-vlab/vct/cmd/pressure"
	"github.com/spf13/cobra"
)

func showHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vct",
		Short: "Vlab's Container Tool",
		Long:  "Vlab's Container Tool, a versatile tool for managing containers on Proxmox VE",
		Args:  cobra.NoArgs,
		Run:   showHelp,
	}
	cmd.CompletionOptions.HiddenDefaultCmd = true
	cmd.AddCommand(df.MakeCmd())
	cmd.AddCommand(findpid.MakeCmd())
	cmd.AddCommand(iostat.MakeCmd())
	cmd.AddCommand(pressure.MakeCmd())
	return cmd
}
