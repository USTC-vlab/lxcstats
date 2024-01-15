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
		fmt.Fprintf(w, "  %s [-cp] [-mp] [-ip]\n", os.Args[0])
		fmt.Fprintln(w)
		flag.PrintDefaults()
	}
}

func main() {
	var cpuPressure,
		memoryPressure,
		ioPressure bool
	flag.BoolVar(&cpuPressure, "cp", false, "list LXC with highest CPU pressures")
	flag.BoolVar(&memoryPressure, "mp", false, "list LXC with highest memory pressures")
	flag.BoolVar(&ioPressure, "ip", false, "list LXC with highest I/O pressures")
	flag.Parse()

	showPressure := cpuPressure || memoryPressure || ioPressure
	if showPressure {
		if cpuPressure {
			pressureMain(CPU)
		}
		if memoryPressure {
			pressureMain(MEMORY)
		}
		if ioPressure {
			pressureMain(IO)
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
