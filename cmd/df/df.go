package df

import (
	"fmt"
	"log"
	"sort"
	"syscall"

	"github.com/USTC-vlab/lxcstats/pkg/cgroup"
	"github.com/USTC-vlab/lxcstats/pkg/util"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

func getInitPid(id string) (int, error) {
	f, err := cgroup.OpenLXC(id, "ns/init.scope/cgroup.procs")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	var pid int
	_, err = fmt.Fscanf(f, "%d", &pid)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

type idAndUsedRootFSSpace struct {
	id    string
	used  uint64
	total uint64
}

func fstatfsMain() error {
	ids, err := cgroup.ListLXC()
	if err != nil {
		return err
	}

	dfRes := make([]idAndUsedRootFSSpace, 0, len(ids))
	for _, id := range ids {
		initPid, err := getInitPid(id)
		if err != nil {
			log.Printf("get init pid error for %s: %v", id, err)
			continue
		}
		rootFSPath := fmt.Sprintf("/proc/%d/root", initPid)
		var statfs syscall.Statfs_t
		err = syscall.Statfs(rootFSPath, &statfs)
		if err != nil {
			log.Printf("statfs error for %s: %v", rootFSPath, err)
			continue
		}
		avail := statfs.Bavail * uint64(statfs.Bsize)
		total := statfs.Blocks * uint64(statfs.Bsize)
		used := total - avail
		dfRes = append(dfRes, idAndUsedRootFSSpace{id: id, used: used, total: total})
	}
	sort.Slice(dfRes, func(i, j int) bool {
		return (dfRes[i].total - dfRes[i].used) < (dfRes[j].total - dfRes[j].used)
	})

	lines := []string{"ID | RootFS used | total"}
	for i, p := range dfRes {
		if i >= 5 {
			break
		}
		line := fmt.Sprintf("%s | %s | %s", p.id, util.FormatSize(p.used), util.FormatSize(p.total))
		lines = append(lines, line)
	}
	fmt.Printf("Top stats of rootfs space\n")
	fmt.Println(columnize.SimpleFormat(lines))
	return nil
}

func MakeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "df",
		Short: "List LXC with highest root filesystem space usage",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fstatfsMain()
		},
	}
}
