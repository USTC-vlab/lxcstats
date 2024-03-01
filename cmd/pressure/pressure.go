package pressure

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

const (
	LXC  = "LXC"
	QEMU = "Qemu"
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

func GetPressure(id, typ string, filename string) (*PSIStats, error) {
	var f *os.File
	var err error
	if typ == LXC {
		f, err = cgroup.OpenLXC(id, filename)
	} else {
		f, err = cgroup.OpenQemu(id, filename)
	}
	if err != nil {
		return nil, fmt.Errorf("open %s for %s (%s): %w", filename, id, typ, err)
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

type idAndPressure struct {
	id       string
	typ      string
	pressure *PSIStats
}

func listPressures(filename string, topN int) error {
	lxcIds, err := cgroup.ListLXC()
	if err != nil {
		log.Printf("ListLXC error (may not have LXC?): %v", err)
		lxcIds = []string{}
	}
	qemuIds, err := cgroup.ListQemu()
	if err != nil {
		log.Printf("ListQemu error (may not have Qemu?): %v", err)
		qemuIds = []string{}
	}

	pressures := make([]idAndPressure, 0, len(lxcIds)+len(qemuIds))
	appendToPressure := func(ids []string, typ string) {
		for _, id := range ids {
			pressure, err := GetPressure(id, typ, filename)
			if err != nil {
				log.Printf("GetPressure error for %s: %v", id, err)
				continue
			}
			pressures = append(pressures, idAndPressure{id, typ, pressure})
		}
	}
	appendToPressure(lxcIds, LXC)
	appendToPressure(qemuIds, QEMU)
	sort.Slice(pressures, func(i, j int) bool {
		return pressures[i].pressure.Some.Avg10 > pressures[j].pressure.Some.Avg10
	})

	lines := []string{"ID | Type | Avg10 | Avg60 | Avg300"}
	for i, p := range pressures {
		if i >= topN {
			break
		}
		line := fmt.Sprintf("%s | %s | %.1f | %.1f | %.1f", p.id, p.typ, p.pressure.Some.Avg10, p.pressure.Some.Avg60, p.pressure.Some.Avg300)
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
