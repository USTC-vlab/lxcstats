package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ryanuber/columnize"
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

func GetIOPressure(id string) (*PSIStats, error) {
	f, err := os.Open(filepath.Join(BaseDir, id, "io.pressure"))
	if err != nil {
		return nil, fmt.Errorf("open io.pressure for %s: %w", id, err)
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

func ioPressureMain() {
	containersDir, err := os.ReadDir(BaseDir)
	if err != nil {
		log.Fatal(err)
	}

	pressures := make([]idAndPressure, 0, len(containersDir))
	for _, containerDir := range containersDir {
		if !containerDir.IsDir() {
			continue
		}
		id := containerDir.Name()
		pressure, err := GetIOPressure(id)
		if err != nil {
			log.Printf("GetIOPressure error for %s: %v", id, err)
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
	fmt.Println(columnize.SimpleFormat(lines))
}
