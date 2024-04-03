package iostat

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/USTC-vlab/vct/pkg/cgroup"
	"github.com/USTC-vlab/vct/pkg/util"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

type IOSingle struct {
	Rbytes, Wbytes, Dbytes uint64
	Rios, Wios, Dios       uint64
}

var zeroIOSingle = IOSingle{}

func (s IOSingle) Zero() bool {
	return s == zeroIOSingle
}

func (l IOSingle) Diff(r IOSingle) IOSingle {
	return IOSingle{
		Rbytes: util.SafeSub(l.Rbytes, r.Rbytes),
		Wbytes: util.SafeSub(l.Wbytes, r.Wbytes),
		Dbytes: util.SafeSub(l.Dbytes, r.Dbytes),
		Rios:   util.SafeSub(l.Rios, r.Rios),
		Wios:   util.SafeSub(l.Wios, r.Wios),
		Dios:   util.SafeSub(l.Dios, r.Dios),
	}
}

type IOStat map[string]IOSingle

func GetIOStat(vmid cgroup.VMID) (IOStat, error) {
	f, err := cgroup.OpenVM(vmid, "io.stat")
	if err != nil {
		return nil, fmt.Errorf("open io.stat for %s: %w", vmid, err)
	}
	defer f.Close()
	stats := make(IOStat)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		name := fields[0]
		var single IOSingle
		for _, s := range fields[1:] {
			parts := strings.Split(s, "=")
			if len(parts) != 2 {
				// error?
				continue
			}
			v, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				continue
			}
			switch parts[0] {
			case "rbytes":
				single.Rbytes = v
			case "wbytes":
				single.Wbytes = v
			case "dbytes":
				single.Dbytes = v
			case "rios":
				single.Rios = v
			case "wios":
				single.Wios = v
			case "dios":
				single.Dios = v
			}
		}
		stats[name] = single
	}
	return stats, nil
}

func iostatMain(diskPath string) error {
	devInfo, err := os.Stat(diskPath)
	if err != nil {
		return err
	}
	stat, ok := devInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("failed to get device number for %s", diskPath)
	}
	major, minor := util.GetDeviceNumbers(stat.Rdev)
	matchString := fmt.Sprintf("%d:%d", major, minor)

	// Enable IO subtree for qemu, if exists
	if err := cgroup.EnableIOForQemu(); err != nil {
		if os.IsNotExist(err) {
			log.Println("qemu.slice not found (not running any KVM), skipping")
		} else {
			log.Printf("EnableIOForQemu error: %v", err)
		}
	}

	cachedStats := make(map[cgroup.VMID]IOSingle)
	for t := range time.NewTicker(1 * time.Second).C {
		vmids, err := cgroup.ListVM()
		if err != nil {
			log.Fatal(err)
		}
		lines := []string{"ID | Type | Rios | Wios | Rbytes | Wbytes"}
		newStats := make(map[cgroup.VMID]IOSingle)
		for _, vmid := range vmids {
			stats, err := GetIOStat(vmid)
			if err != nil {
				log.Printf("GetIOStat error for %s: %v", vmid, err)
				continue
			}
			stat := stats[matchString]
			newStats[vmid] = stat
			oldStat, ok := cachedStats[vmid]
			if ok {
				diff := stat.Diff(oldStat)
				if !diff.Zero() {
					line := fmt.Sprintf("%s | %s | %d | %d | %s | %s", vmid.Id, vmid.Type,
						diff.Rios, diff.Wios, util.FormatSize(diff.Rbytes), util.FormatSize(diff.Wbytes))
					lines = append(lines, line)
				}
			}
		}
		fmt.Println(t.Format("2006-01-02 15:04:05"))
		fmt.Println(columnize.SimpleFormat(lines))
		fmt.Println()
		cachedStats = newStats
	}
	return nil
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iostat DISK",
		Short: "Show I/O statistics for disks",
		Long:  "Show I/O statistics for DISK",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return iostatMain(args[0])
		},
	}
	return cmd
}
