package cmd

import (
	"github.com/USTC-vlab/vct/cmd/df"
	"github.com/USTC-vlab/vct/cmd/findpid"
	"github.com/USTC-vlab/vct/cmd/iolimit"
	"github.com/USTC-vlab/vct/cmd/iostat"
	"github.com/USTC-vlab/vct/cmd/killall"
	"github.com/USTC-vlab/vct/cmd/nsenter"
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
		Short: "Vlab's Control Tool",
		Long:  "Vlab's Control Tool, a versatile tool for managing containers and virtual machines on Proxmox VE",
		Args:  cobra.NoArgs,
	}
	cmd.CompletionOptions.HiddenDefaultCmd = true
	cmd.AddCommand(
		df.MakeCmd(),
		findpid.MakeCmd(),
		iostat.MakeCmd(),
		iolimit.MakeCmd(),
		killall.MakeCmd(),
		nsenter.MakeCmd(),
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
