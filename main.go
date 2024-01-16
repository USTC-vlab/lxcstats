package main

import (
	"flag"
	"fmt"
	"os"
)

const BaseDir = "/sys/fs/cgroup/lxc"

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(w, "  %s <disk>\n", os.Args[0])
		fmt.Fprintf(w, "  %s [-cp] [-mp] [-ip] [-df]\n", os.Args[0])
		fmt.Fprintln(w)
		flag.PrintDefaults()
	}
}

func main() {
	var cpuPressure,
		memoryPressure,
		ioPressure,
		rootFSSpace bool
	flag.BoolVar(&cpuPressure, "cp", false, "list LXC with highest CPU pressures")
	flag.BoolVar(&memoryPressure, "mp", false, "list LXC with highest memory pressures")
	flag.BoolVar(&ioPressure, "ip", false, "list LXC with highest I/O pressures")
	flag.BoolVar(&rootFSSpace, "df", false, "list LXC with highest root filesystem space usage")
	flag.Parse()

	showStats := cpuPressure || memoryPressure || ioPressure || rootFSSpace
	if showStats {
		if cpuPressure {
			pressureMain(CPU)
		}
		if memoryPressure {
			pressureMain(MEMORY)
		}
		if ioPressure {
			pressureMain(IO)
		}
		if rootFSSpace {
			fstatfsMain()
		}
	} else {
		if flag.NArg() != 1 {
			flag.Usage()
			os.Exit(1)
		}
		diskPath := flag.Arg(0)
		ioStatMain(diskPath)
	}
}
