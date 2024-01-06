package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ryanuber/columnize"
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
		Rbytes: SafeSub(l.Rbytes, r.Rbytes),
		Wbytes: SafeSub(l.Wbytes, r.Wbytes),
		Dbytes: SafeSub(l.Dbytes, r.Dbytes),
		Rios:   SafeSub(l.Rios, r.Rios),
		Wios:   SafeSub(l.Wios, r.Wios),
		Dios:   SafeSub(l.Dios, r.Dios),
	}
}

type IOStat map[string]IOSingle

func GetIOStat(id string) (IOStat, error) {
	f, err := os.Open(filepath.Join(BaseDir, id, "io.stat"))
	if err != nil {
		return nil, fmt.Errorf("open io.stat for %s: %w", id, err)
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

func ioStatMain(diskPath string) {
	devInfo, err := os.Stat(diskPath)
	if err != nil {
		log.Fatal(err)
	}
	stat, ok := devInfo.Sys().(*syscall.Stat_t)
	if !ok {
		log.Fatal("Failed to get device id")
	}
	devId := stat.Rdev
	major := (devId >> 8) & 0xfff
	minor := (devId & 0xff) | ((devId >> 12) & 0xfff00)
	matchString := fmt.Sprintf("%d:%d", major, minor)

	cachedStats := make(map[string]IOSingle)
	for t := range time.NewTicker(1 * time.Second).C {
		containersDir, err := os.ReadDir(BaseDir)
		if err != nil {
			log.Fatal(err)
		}
		lines := []string{"ID | Rios | Wios | Rbytes | Wbytes"}
		newStats := make(map[string]IOSingle)
		for _, containerDir := range containersDir {
			if !containerDir.IsDir() {
				continue
			}
			id := containerDir.Name()
			stats, err := GetIOStat(id)
			if err != nil {
				log.Printf("GetIOStat error for %s: %v", id, err)
				continue
			}
			stat := stats[matchString]
			newStats[id] = stat
			oldStat, ok := cachedStats[id]
			if ok {
				diff := stat.Diff(oldStat)
				if !diff.Zero() {
					line := fmt.Sprintf("%s | %d | %d | %s | %s", id,
						diff.Rios, diff.Wios, FormatSize(diff.Rbytes), FormatSize(diff.Wbytes))
					lines = append(lines, line)
				}
			}
		}
		fmt.Println(t.Format("2006-01-02 15:04:05"))
		fmt.Println(columnize.SimpleFormat(lines))
		fmt.Println()
		cachedStats = newStats
	}
}
