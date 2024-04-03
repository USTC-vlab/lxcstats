package iolimit

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/USTC-vlab/vct/pkg/cgroup"
	"github.com/USTC-vlab/vct/pkg/pve"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iolimit [--riops COUNT] [--wiops COUNT] [--rbps BYTES] [--wbps BYTES] CTID...",
		Short: "Set I/O limits for an LXC container",
		Long: `Set I/O limits for an LXC container.
If no setting is given, print the current limits.
Note that zero means "don't change". Use -1 to remove an existing limit (set to "max").`,
	}
	flags := cmd.Flags()
	var iops cgroup.IOPS
	flags.Int64VarP(&iops.Rbps, "rbps", "R", 0, "set read bandwidth limit")
	flags.Int64VarP(&iops.Wbps, "wbps", "W", 0, "set write bandwidth limit")
	flags.Int64VarP(&iops.Riops, "riops", "r", 0, "set read IOPS limit")
	flags.Int64VarP(&iops.Wiops, "wiops", "w", 0, "set write IOPS limit")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		errs := make([]error, 0, len(args))
		for _, arg := range args {
			err := iolimitMain(cmd, arg, iops)
			if err != nil {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	}
	return cmd
}

func showIOLimit(w io.Writer, ctid string) error {
	iopss, err := cgroup.GetIOPSForLXC(ctid)
	if err != nil {
		return err
	}
	if len(iopss) == 0 {
		fmt.Fprintf(w, "No I/O limits for %s\n", ctid)
		return nil
	}

	fmt.Fprintf(w, "I/O limits for %s:\n", ctid)
	lines := []string{"Device | Rbps | Wbps | Riops | Wiops"}
	for _, iops := range iopss {
		line := fmt.Sprintf("%d:%d | %s | %s | %s | %s",
			iops.Major, iops.Minor,
			cgroup.ItoaZeroMax(iops.Rbps), cgroup.ItoaZeroMax(iops.Wbps),
			cgroup.ItoaZeroMax(iops.Riops), cgroup.ItoaZeroMax(iops.Wiops))
		lines = append(lines, line)
	}
	fmt.Fprintln(w, columnize.SimpleFormat(lines))
	return nil
}

func setIOLimit(ctid string, iops cgroup.IOPS) error {

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

func iolimitMain(cmd *cobra.Command, ctid string, iops cgroup.IOPS) error {
	if iops.IsZero() {
		return showIOLimit(cmd.OutOrStdout(), ctid)
	}
	return setIOLimit(ctid, iops)
}
