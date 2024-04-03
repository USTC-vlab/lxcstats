package cgroup

import (
	"os"
	"strconv"
	"strings"
)

// IOPS represents I/O limits for a single device.
// Zero means "unset" or "untouched".
// -1 represents no limit and will be translated to "max".
type IOPS struct {
	Rbps, Wbps, Riops, Wiops int64
}

func itoaZeroMax(i int64) string {
	if i <= 0 {
		return "max"
	}
	return strconv.FormatInt(i, 10)
}

func (i IOPS) String() string {
	parts := make([]string, 0, 4)
	if i.Rbps != 0 {
		parts = append(parts, "rbps="+itoaZeroMax(i.Rbps))
	}
	if i.Wbps != 0 {
		parts = append(parts, "wbps="+itoaZeroMax(i.Wbps))
	}
	if i.Riops != 0 {
		parts = append(parts, "riops="+itoaZeroMax(i.Riops))
	}
	if i.Wiops != 0 {
		parts = append(parts, "wiops="+itoaZeroMax(i.Wiops))
	}
	return strings.Join(parts, " ")
}

// IOPSLine represents a line in the io.max file.
type IOPSLine struct {
	Major, Minor uint64
	IOPS         IOPS
}

func (l IOPSLine) String() string {
	return strconv.FormatUint(l.Major, 10) + ":" + strconv.FormatUint(l.Minor, 10) + " " + l.IOPS.String()
}

func SetIOPSForLXC(id string, iops IOPSLine) error {
	filename := GetFilenameLXC(id, "io.max")
	return os.WriteFile(filename, []byte(iops.String()+"\n"), 0)
}
