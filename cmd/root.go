package cmd

import (
	"github.com/USTC-vlab/vct/cmd/df"
	"github.com/USTC-vlab/vct/cmd/findpid"
	"github.com/USTC-vlab/vct/cmd/iostat"
	"github.com/USTC-vlab/vct/cmd/killall"
	"github.com/USTC-vlab/vct/cmd/pressure"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:    "version",
	Short:  "Show version",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(cmd.Root().Name(), version)
	},
}

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
		killall.MakeCmd(),
		pressure.MakeCmd(),
		versionCmd,
	)
	pVersion := cmd.Flags().BoolP("version", "v", false, "show version")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if *pVersion {
			versionCmd.Run(versionCmd, args)
		} else {
			cmd.Help()
		}
	}
	return cmd
}
