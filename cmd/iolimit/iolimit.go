package iolimit

import (
	"errors"
	"fmt"
	"strings"

	"github.com/USTC-vlab/vct/pkg/cgroup"
	"github.com/USTC-vlab/vct/pkg/pve"
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
	flags.Int64VarP(&iops.Rbps, "rbps", "R", 0, "set read bandwidth limit")
	flags.Int64VarP(&iops.Wbps, "wbps", "W", 0, "set write bandwidth limit")
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
	pveStorage, err := pve.GetStorage()
	if err != nil {
		return err
	}

	config, err := pve.GetLXCConfig(ctid)
	if err != nil {
		return err
	}
	rootfs := config["rootfs"]
	rootfsIdent := strings.SplitN(rootfs, ",", 2)[0]
	rootfsParts := strings.Split(rootfsIdent, ":")
	if len(rootfsParts) != 2 {
		return fmt.Errorf("invalid rootfs %s", rootfsIdent)
	}

	major, minor, err := pve.GetBlockDevForStorage(rootfsParts[0], rootfsParts[1], pveStorage)
	if err != nil {
		return err
	}
	iopsline := cgroup.IOPSLine{
		Major: major,
		Minor: minor,
		IOPS:  iops,
	}
	return cgroup.SetIOPSForLXC(ctid, iopsline)
}
