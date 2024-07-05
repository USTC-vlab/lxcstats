package top

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/USTC-vlab/vct/cmd/findpid"
	"github.com/prometheus/procfs"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

func runE(cmd *cobra.Command, args []string) error {
	fs, err := procfs.NewDefaultFS()
	if err != nil {
		return err
	}

	count := 10
	if len(args) > 0 {
		count, err = strconv.Atoi(args[0])
		if err != nil {
			return err
		}
	}

	// Get all process information
	procs, err := fs.AllProcs()
	if err != nil {
		return err
	}
	procs = slices.DeleteFunc(procs, func(p procfs.Proc) bool {
		cgroups, err := p.Cgroups()
		if err != nil {
			cmd.PrintErrln(err)
			return true
		}
		if len(cgroups) == 0 {
			return true
		}
		return !strings.HasPrefix(cgroups[0].Path, "/lxc/")
	})

	// Collect their stats
	stats := make([]procfs.ProcStat, 0, len(procs))
	for _, p := range procs {
		stat, err := p.Stat()
		if err != nil {
			cmd.PrintErrln(err)
			continue
		}
		stats = append(stats, stat)
	}
	// Sort stats by utime+stime, descending
	slices.SortFunc(stats, func(a, b procfs.ProcStat) int {
		switch {
		case a.UTime+a.STime > b.UTime+b.STime:
			return -1
		case a.UTime+a.STime < b.UTime+b.STime:
			return 1
		default:
			return 0
		}
	})

	// Print the top N processes
	lines := make([]string, 1, 1+count)
	lines[0] = "PID | CT | Name | Time"
	for _, stat := range stats[:count] {
		ctid, err := findpid.FindLXCForPid(strconv.Itoa(stat.PID))
		if err != nil {
			cmd.PrintErrln(err)
			continue
		}
		cpuTime := time.Duration(stat.CPUTime() * float64(time.Second)).Truncate(time.Second)
		lines = append(lines, fmt.Sprintf("%d | %s | %s | %s", stat.PID, ctid, stat.Comm, cpuTime))
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), columnize.SimpleFormat(lines))
	return err
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "top [N]",
		Short: "Show top N processes with high total CPU time",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runE,
	}
	return cmd
}
