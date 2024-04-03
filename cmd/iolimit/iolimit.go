package iolimit

import (
	"errors"

	"github.com/USTC-vlab/vct/pkg/cgroup"
	"github.com/spf13/cobra"
)

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iolimit [--riops] [--wiops] [--rbps] [--wbps] CTID...",
		Short: "Set I/O limits for an LXC container",
		Long:  "Set I/O limits for an LXC container. Use -1 to remove existing limit.",
		Args:  cobra.MinimumNArgs(1),
	}
	flags := cmd.Flags()
	var iops cgroup.IOPS
	flags.Int64VarP(&iops.Rbps, "rbps", "", 0, "set read bandwidth limit")
	flags.Int64VarP(&iops.Wbps, "wbps", "", 0, "set write bandwidth limit")
	flags.Int64VarP(&iops.Riops, "riops", "r", 0, "set read IOPS limit")
	flags.Int64VarP(&iops.Wiops, "wiops", "w", 0, "set write IOPS limit")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		errs := make([]error, 0, len(args))
		for _, arg := range args {
			err := iolimitMain(arg, iops)
			if err != nil {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	}
	return cmd
}

func iolimitMain(ctid string, iops cgroup.IOPS) error {
	return cgroup.SetIOPSForLXC(ctid, cgroup.IOPSLine{IOPS: iops})
}
