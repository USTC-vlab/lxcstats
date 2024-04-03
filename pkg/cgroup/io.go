package cgroup

import (
	"bufio"
	"fmt"
	"io"
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

func ItoaZeroMax(i int64) string {
	if i <= 0 {
		return "max"
	}
	return strconv.FormatInt(i, 10)
}

func (i IOPS) String() string {
	parts := make([]string, 0, 4)
	if i.Rbps != 0 {
		parts = append(parts, "rbps="+ItoaZeroMax(i.Rbps))
	}
	if i.Wbps != 0 {
		parts = append(parts, "wbps="+ItoaZeroMax(i.Wbps))
	}
	if i.Riops != 0 {
		parts = append(parts, "riops="+ItoaZeroMax(i.Riops))
	}
	if i.Wiops != 0 {
		parts = append(parts, "wiops="+ItoaZeroMax(i.Wiops))
	}
	return strings.Join(parts, " ")
}

func (i IOPS) IsZero() bool {
	return i.Rbps == 0 && i.Wbps == 0 && i.Riops == 0 && i.Wiops == 0
}

func ParseIOPS(line string) (i IOPS) {
	parts := strings.Fields(line)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "rbps":
			i.Rbps, _ = strconv.ParseInt(kv[1], 10, 64)
		case "wbps":
			i.Wbps, _ = strconv.ParseInt(kv[1], 10, 64)
		case "riops":
			i.Riops, _ = strconv.ParseInt(kv[1], 10, 64)
		case "wiops":
			i.Wiops, _ = strconv.ParseInt(kv[1], 10, 64)
		}
	}
	return
}

// IOPSLine represents a line in the io.max file.
type IOPSLine struct {
	Major, Minor uint64
	IOPS
}

func (l IOPSLine) String() string {
	return strconv.FormatUint(l.Major, 10) + ":" + strconv.FormatUint(l.Minor, 10) + " " + l.IOPS.String()
}

func parseIOPSLines(r io.Reader) ([]IOPSLine, error) {
	scanner := bufio.NewScanner(r)
	iopss := make([]IOPSLine, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		var major, minor uint64
		_, err := fmt.Sscanf(fields[0], "%d:%d", &major, &minor)
		if err != nil {
			return iopss, err
		}
		iops := ParseIOPS(strings.Join(fields[1:], " "))
		iopss = append(iopss, IOPSLine{
			Major: major,
			Minor: minor,
			IOPS:  iops,
		})
	}
	return iopss, nil
}

func GetIOPSForLXC(id string) ([]IOPSLine, error) {
	filename := GetFilenameLXC(id, "io.max")
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseIOPSLines(f)
}

func SetIOPSForLXC(id string, iops IOPSLine) error {
	filename := GetFilenameLXC(id, "io.max")
	return os.WriteFile(filename, []byte(iops.String()+"\n"), 0)
}
