package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const BaseDir = "/sys/fs/cgroup"

type VMID struct {
	Id   string
	Type string
}

func (v VMID) String() string {
	return fmt.Sprintf("%s (%s)", v.Id, v.Type)
}

const (
	LXC  = "LXC"
	QEMU = "Qemu"
)

func OpenVM(vmid VMID, filename string) (*os.File, error) {
	switch vmid.Type {
	case LXC:
		return OpenLXC(vmid.Id, filename)
	case QEMU:
		return OpenQemu(vmid.Id, filename)
	default:
		return nil, fmt.Errorf("unknown vm type %s", vmid.Type)
	}
}

func OpenLXC(id string, filename string) (*os.File, error) {
	return os.Open(GetFilenameLXC(id, filename))
}

func OpenQemu(id string, filename string) (*os.File, error) {
	return os.Open(GetFilenameQemu(id, filename))
}

func GetFilenameVM(vmid VMID, filename string) string {
	switch vmid.Type {
	case LXC:
		return GetFilenameLXC(vmid.Id, filename)
	case QEMU:
		return GetFilenameQemu(vmid.Id, filename)
	default:
		panic("unknown vm type " + vmid.Type)
	}
}

func GetFilenameLXC(id string, filename string) string {
	return filepath.Join(BaseDir, "lxc", id, filename)
}

func GetFilenameQemu(id string, filename string) string {
	return filepath.Join(BaseDir, "qemu.slice", id+".scope", filename)
}

func EnableIOForQemu() error {
	subtreeControlFilename := filepath.Join(BaseDir, "qemu.slice", "cgroup.subtree_control")
	return os.WriteFile(subtreeControlFilename, []byte("+io"), 0)
}

func ListVM() ([]VMID, error) {
	lxc, err := ListLXC()
	if err != nil {
		return nil, err
	}
	qemu, err := ListQemu()
	if err != nil {
		return nil, err
	}
	ids := make([]VMID, 0, len(lxc)+len(qemu))
	for _, id := range lxc {
		ids = append(ids, VMID{id, LXC})
	}
	for _, id := range qemu {
		ids = append(ids, VMID{id, QEMU})
	}
	return ids, nil
}

func ListLXC() ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(BaseDir, "lxc"))
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
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

func ListQemu() ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(BaseDir, "qemu.slice"))
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			id := strings.TrimSuffix(entry.Name(), ".scope")
			ids = append(ids, id)
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
