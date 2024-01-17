package cmd

import (
	"github.com/USTC-vlab/lxcstats/cmd/df"
	"github.com/USTC-vlab/lxcstats/cmd/iostat"
	"github.com/USTC-vlab/lxcstats/cmd/pressure"
	"github.com/spf13/cobra"
)

func showHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vct",
		Short: "",
		Args:  cobra.NoArgs,
		Run:   showHelp,
	}
	cmd.AddCommand(df.MakeCmd())
	cmd.AddCommand(iostat.MakeCmd())
	cmd.AddCommand(pressure.MakeCmd())
	return cmd
}
