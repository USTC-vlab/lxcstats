package df

import (
	"fmt"
	"io"
	"log"
	"sort"
	"syscall"

	"github.com/USTC-vlab/vct/pkg/cgroup"
	"github.com/USTC-vlab/vct/pkg/util"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

type idAndUsedRootFSSpace struct {
	id    string
	used  uint64
	total uint64
}

func fstatfsMain(w io.Writer, topN int) error {
	ids, err := cgroup.ListLXC()
	if err != nil {
		return err
	}

	dfRes := make([]idAndUsedRootFSSpace, 0, len(ids))
	for _, id := range ids {
		initPid, err := cgroup.GetLXCInitPid(id)
		if err != nil {
			log.Printf("get init pid error for %s: %v", id, err)
			continue
		}
		rootFSPath := fmt.Sprintf("/proc/%d/root", initPid)
		var statfs syscall.Statfs_t
		if err = syscall.Statfs(rootFSPath, &statfs); err != nil {
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

	lines := []string{"ID | RootFS used | total | %"}
	for i, p := range dfRes {
		if i >= topN {
			break
		}
		line := fmt.Sprintf("%s | %s | %s | %.1f", p.id,
			util.FormatSize(p.used),
			util.FormatSize(p.total),
			float64(p.used)/float64(p.total)*100)
		lines = append(lines, line)
	}
	fmt.Fprintln(w, columnize.SimpleFormat(lines))
	return nil
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "df",
		Short: "List LXC with highest root filesystem space usage",
		Args:  cobra.NoArgs,
	}
	pN := cmd.Flags().IntP("count", "n", 5, "number of containers to show")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return fstatfsMain(cmd.OutOrStdout(), *pN)
	}
	return cmd
}
