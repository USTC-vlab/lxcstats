package pve

import "os/exec"

const PctPath = "/usr/sbin/pct"

func Pct(args ...string) *exec.Cmd {
	return exec.Command(PctPath, args...)
}

func Stop(vmid string) *exec.Cmd {
	return Pct("stop", vmid)
}
