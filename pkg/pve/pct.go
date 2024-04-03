package pve

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	EtcPve  = "/etc/pve"
	PctPath = "/usr/sbin/pct"
)

func PctCmd(args ...string) *exec.Cmd {
	return exec.Command(PctPath, args...)
}

func StopCmd(vmid string) *exec.Cmd {
	return PctCmd("stop", vmid)
}

func parseConfig(r io.Reader) map[string]string {
	config := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if comment, ok := strings.CutPrefix(line, "#"); ok {
			config["#"] += comment + "\n"
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			// Invalid line?
			continue
		}
		config[strings.TrimSuffix(fields[0], ":")] = fields[1]
	}
	return config
}

func GetConfig(typ, vmid string) (map[string]string, error) {
	if typ == "qemu" {
		typ = "qemu-server"
	}
	f, err := os.Open(fmt.Sprintf("%s/%s/%s.conf", EtcPve, typ, vmid))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseConfig(f), nil
}

func GetLXCConfig(vmid string) (map[string]string, error) {
	return GetConfig("lxc", vmid)
}

func GetQemuConfig(vmid string) (map[string]string, error) {
	return GetConfig("qemu", vmid)
}
