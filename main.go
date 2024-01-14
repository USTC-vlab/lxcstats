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
		fmt.Fprintf(w, "  %s -ip\n", os.Args[0])
		fmt.Fprintf(w, "  %s -cp\n", os.Args[0])
		fmt.Fprintf(w, "  %s -mp\n", os.Args[0])
		fmt.Fprintln(w)
		flag.PrintDefaults()
	}
}

func main() {
	var ioPressure bool
	var cpuPressure bool
	var memoryPressure bool
	flag.BoolVar(&ioPressure, "ip", false, "list LXC with highest I/O pressures")
	flag.BoolVar(&cpuPressure, "cp", false, "list LXC with highest CPU pressures")
	flag.BoolVar(&memoryPressure, "mp", false, "list LXC with highest memory pressures")
	flag.Parse()

	if ioPressure {
		pressureMain(IO)
	} else if cpuPressure {
		pressureMain(CPU)
	} else if memoryPressure {
		pressureMain(MEMORY)
	} else {
		if flag.NArg() != 1 {
			flag.Usage()
			os.Exit(1)
		}
		diskPath := flag.Arg(0)
		ioStatMain(diskPath)
	}
}
