package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
)

const BaseDir = "/sys/fs/cgroup"

func Open(pathComp ...string) (*os.File, error) {
	pathComp = append([]string{BaseDir}, pathComp...)
	return os.Open(filepath.Join(pathComp...))
}

func OpenLXC(id string, pathComp ...string) (*os.File, error) {
	pathComp = append([]string{BaseDir, "lxc", id}, pathComp...)
	return os.Open(filepath.Join(pathComp...))
}

func ListLXC() ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(BaseDir, "lxc"))
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			ids = append(ids, entry.Name())
		}
	}
	return ids, nil
}

func GetLXCInitPid(id string) (int, error) {
	f, err := OpenLXC(id, "ns/init.scope/cgroup.procs")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	var pid int
	_, err = fmt.Fscanf(f, "%d", &pid)
	if err != nil {
		return 0, err
	}
	return pid, nil
}
