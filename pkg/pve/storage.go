package pve

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"slices"
	"strings"
	"syscall"
)

const StorageConf = "/etc/pve/storage.cfg"

type PVEStorage struct {
	Type string
	Name string
	Attr map[string]string
}

func newPVEStorage() PVEStorage {
	return PVEStorage{
		Attr: make(map[string]string),
	}
}

func parseStorage(r io.Reader) []PVEStorage {
	items := make([]PVEStorage, 0)
	item := newPVEStorage()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			// Invalid line?
			continue
		}
		if strings.HasSuffix(fields[0], ":") {
			if item.Type != "" {
				items = append(items, item)
				item = newPVEStorage()
			}
			item.Type = strings.TrimSuffix(fields[0], ":")
			item.Name = fields[1]
		} else {
			item.Attr[fields[0]] = fields[1]
		}
	}
	if item.Type != "" {
		items = append(items, item)
	}
	return items
}

func GetStorage() ([]PVEStorage, error) {
	f, err := os.Open(StorageConf)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseStorage(f), nil
}

func devMajorMinor(device uint64) (major, minor uint64) {
	major = (device >> 8) & 0xfff
	minor = (device & 0xff) | ((device >> 12) & 0xfff00)
	return
}

func getBlockDevForDir(dir string) (uint64, uint64, error) {
	fileInfo, err := os.Lstat(dir)
	if err != nil {
		return 0, 0, err
	}
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, fmt.Errorf("failed to get backing device for %s: %w", dir, err)
	}
	major, minor := devMajorMinor(stat.Dev)
	return major, minor, nil
}

func getBlockDevForLVM(vgname, lvname string) (uint64, uint64, error) {
	fileInfo, err := os.Stat(fmt.Sprintf("/dev/%s/%s", vgname, lvname))
	if err != nil {
		return 0, 0, err
	}
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, fmt.Errorf("failed to get device number for %s/%s: %w", vgname, lvname, err)
	}
	major, minor := devMajorMinor(stat.Rdev)
	return major, minor, nil
}

// Finds the block device for the given storage and name.
// The "aux" parameter is used to help determine the type of the storage.
func GetBlockDevForStorage(storage, name string, aux []PVEStorage) (uint64, uint64, error) {
	i := slices.IndexFunc(aux, func(s PVEStorage) bool {
		return s.Name == name
	})
	if i == -1 {
		return 0, 0, fs.ErrNotExist
	}
	info := aux[i]
	switch info.Type {
	case "dir":
		return getBlockDevForDir(storage)
	case "lvm", "lvmthin":
		vgname := info.Attr["vgname"]
		return getBlockDevForLVM(vgname, name)
	default:
		return 0, 0, fmt.Errorf("unsupported storage type %s", info.Type)
	}
}
