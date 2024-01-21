package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
)

const BaseDir = "/sys/fs/cgroup"

func Open(filename string) (*os.File, error) {
	return os.Open(filepath.Join(BaseDir, filename))
}

func OpenLXC(id string, filename string) (*os.File, error) {
	return os.Open(GetFilenameLXC(id, filename))
}

func GetFilenameLXC(id string, filename string) string {
	return filepath.Join(BaseDir, "lxc", id, filename)
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

func KillLXC(id string) error {
	return os.WriteFile(GetFilenameLXC(id, "cgroup.kill"), []byte("1"), 0)
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
