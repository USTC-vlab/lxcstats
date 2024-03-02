package pressure

import (
	"bufio"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/USTC-vlab/vct/pkg/cgroup"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

const (
	IO     = "io.pressure"
	CPU    = "cpu.pressure"
	MEMORY = "memory.pressure"
)

const pressureLineFormat = "avg10=%f avg60=%f avg300=%f total=%d"

type PSILine struct {
	Avg10  float64
	Avg60  float64
	Avg300 float64
	Total  uint64
}

type PSIStats struct {
	Some *PSILine
	Full *PSILine
}

func GetPressure(vmid cgroup.VMID, filename string) (*PSIStats, error) {
	f, err := cgroup.OpenVM(vmid, filename)
	if err != nil {
		return nil, fmt.Errorf("open %s for %s: %w", filename, vmid, err)
	}
	defer f.Close()
	stats := &PSIStats{
		Some: &PSILine{},
		Full: &PSILine{},
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		prefix := strings.Fields(l)[0]
		var psi *PSILine
		switch prefix {
		case "some":
			psi = stats.Some
		case "full":
			psi = stats.Full
		default:
			continue
		}
		_, err := fmt.Sscanf(l, fmt.Sprintf("%s %s", prefix, pressureLineFormat), &psi.Avg10, &psi.Avg60, &psi.Avg300, &psi.Total)
		if err != nil {
			return nil, err
		}
	}
	return stats, nil
}

type vmidAndPressure struct {
	vmid     cgroup.VMID
	pressure *PSIStats
}

func listPressures(filename string, topN int) error {
	vmids, err := cgroup.ListVM()
	if err != nil {
		return err
	}

	pressures := make([]vmidAndPressure, 0, len(vmids))
	for _, vmid := range vmids {
		pressure, err := GetPressure(vmid, filename)
		if err != nil {
			log.Printf("GetPressure error for %s: %v", vmid, err)
			continue
		}
		pressures = append(pressures, vmidAndPressure{vmid, pressure})
	}

	sort.Slice(pressures, func(i, j int) bool {
		return pressures[i].pressure.Some.Avg10 > pressures[j].pressure.Some.Avg10
	})

	lines := []string{"ID | Type | Avg10 | Avg60 | Avg300"}
	for i, p := range pressures {
		if i >= topN {
			break
		}
		line := fmt.Sprintf("%s | %s | %.1f | %.1f | %.1f", p.vmid.Id, p.vmid.Type, p.pressure.Some.Avg10, p.pressure.Some.Avg60, p.pressure.Some.Avg300)
		lines = append(lines, line)
	}
	fmt.Printf("Top %d containers/VMs with %s\n", topN, filename)
	fmt.Println(columnize.SimpleFormat(lines))
	return nil
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pressure [-c | --cpu] [-m | --memory] [-i | --io]",
		Short:   "List LXC and Qemu with highest pressures",
		Aliases: []string{"p"},
		Args:    cobra.NoArgs,
	}
	flags := cmd.Flags()
	flags.SortFlags = false
	pCPU := flags.BoolP("cpu", "c", false, "list LXC and Qemu with highest CPU pressures")
	pMemory := flags.BoolP("memory", "m", false, "list LXC and Qemu with highest memory pressures")
	pIO := flags.BoolP("io", "i", false, "list LXC and Qemu with highest I/O pressures")
	pN := flags.IntP("count", "n", 5, "number of containers & VMs to show")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		enabled := []bool{*pCPU, *pMemory, *pIO}
		filename := []string{CPU, MEMORY, IO}
		anyEnabled := false
		for i, e := range enabled {
			if !e {
				continue
			}
			anyEnabled = true
			if err := listPressures(filename[i], *pN); err != nil {
				return err
			}
		}
		if !anyEnabled {
			return fmt.Errorf("no pressure type specified")
		}
		return nil
	}
	return cmd
}
