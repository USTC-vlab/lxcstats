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
		fmt.Fprintf(w, "  %s -p\n", os.Args[0])
		fmt.Fprintln(w)
		flag.PrintDefaults()
	}
}

func main() {
	var pressureMode bool
	flag.BoolVar(&pressureMode, "p", false, "list LXC with highest I/O pressures instead")
	flag.Parse()

	if pressureMode {
		ioPressureMain()
		return
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	diskPath := flag.Arg(0)
	ioStatMain(diskPath)
}
