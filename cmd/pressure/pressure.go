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

func GetPressure(id string, filename string) (*PSIStats, error) {
	f, err := cgroup.OpenLXC(id, filename)
	if err != nil {
		return nil, fmt.Errorf("open %s for %s: %w", filename, id, err)
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
	pressure *PSIStats
}

func listPressures(filename string) error {
	ids, err := cgroup.ListLXC()
	if err != nil {
		log.Fatal(err)
	}

	pressures := make([]idAndPressure, 0, len(ids))
	for _, id := range ids {
		pressure, err := GetPressure(id, filename)
		if err != nil {
			log.Printf("GetPressure error for %s: %v", id, err)
			continue
		}
		pressures = append(pressures, idAndPressure{id, pressure})
	}
	sort.Slice(pressures, func(i, j int) bool {
		return pressures[i].pressure.Some.Avg10 > pressures[j].pressure.Some.Avg10
	})

	lines := []string{"ID | Avg10 | Avg60 | Avg300"}
	for i, p := range pressures {
		if i >= 5 {
			break
		}
		line := fmt.Sprintf("%s | %.1f | %.1f | %.1f", p.id, p.pressure.Some.Avg10, p.pressure.Some.Avg60, p.pressure.Some.Avg300)
		lines = append(lines, line)
	}
	fmt.Printf("Top stats from %s\n", filename)
	fmt.Println(columnize.SimpleFormat(lines))
	return nil
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pressure [-c | --cpu] [-m | --memory] [-i | --io]",
		Short:   "List LXC with highest pressures",
		Aliases: []string{"p"},
		Args:    cobra.NoArgs,
	}
	flags := cmd.Flags()
	flags.SortFlags = false
	pCPU := flags.BoolP("cpu", "c", false, "list LXC with highest CPU pressures")
	pMemory := flags.BoolP("memory", "m", false, "list LXC with highest memory pressures")
	pIO := flags.BoolP("io", "i", false, "list LXC with highest I/O pressures")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		enabled := []bool{*pCPU, *pMemory, *pIO}
		filename := []string{CPU, MEMORY, IO}
		for i, e := range enabled {
			if !e {
				continue
			}
			if err := listPressures(filename[i]); err != nil {
				return err
			}
		}
		return nil
	}
	return cmd
}
