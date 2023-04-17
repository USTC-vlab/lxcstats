package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

type IO struct {
	TotalRead  uint64
	TotalWrite uint64
}

var cachedIO = make(map[string]IO)

func main() {
	regex := regexp.MustCompile("rios=([0-9]+) wios=([0-9]+)")
	diskPtr := flag.String("disk", "", "Disk to monitor")
	flag.Parse()
	if *diskPtr == "" {
		log.Fatal("No disk specified")
	}
	// Get disk device id
	devInfo, err := os.Stat(*diskPtr)
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
	for {
		// list dirs in /sys/fs/cgroup/lxc
		containersDir, err := os.ReadDir("/sys/fs/cgroup/lxc")
		if err != nil {
			log.Fatal(err)
		}
		for _, containerDir := range containersDir {
			if !containerDir.IsDir() {
				continue
			}
			// read from io.stat
			id := containerDir.Name()
			ioStat, err := os.Open("/sys/fs/cgroup/lxc/" + id + "/io.stat")
			if err != nil {
				log.Printf("Failed to open io.stat for %s: %v", id, err)
				continue
			}
			scanner := bufio.NewScanner(ioStat)
			for scanner.Scan() {
				line := scanner.Text()
				if line[:len(matchString)] == matchString {
					matches := regex.FindStringSubmatch(line)
					rios, err := strconv.ParseUint(matches[1], 10, 64)
					if err != nil {
						log.Printf("Failed to parse rios for %s: %v", id, err)
						continue
					}
					wios, err := strconv.ParseUint(matches[2], 10, 64)
					if err != nil {
						log.Printf("Failed to parse wios for %s: %v", id, err)
						continue
					}
					// get latest total io
					latestIO, ok := cachedIO[id]
					if ok {
						if rios < latestIO.TotalRead || wios < latestIO.TotalWrite {
							// do nothing
						} else {
							readIops := rios - latestIO.TotalRead
							writeIops := wios - latestIO.TotalWrite
							if readIops > 0 || writeIops > 0 {
								fmt.Printf("%s: %d read iops, %d write iops\n", id, readIops, writeIops)
							}
						}
					}
					cachedIO[id] = IO{rios, wios}
					break
				}
			}
		}
		fmt.Println()
		time.Sleep(1 * time.Second)
	}
}
