package pve

import "os/exec"

const PctPath = "/usr/sbin/pct"

func PctCmd(args ...string) *exec.Cmd {
	return exec.Command(PctPath, args...)
}

func StopCmd(vmid string) *exec.Cmd {
	return PctCmd("stop", vmid)
}
