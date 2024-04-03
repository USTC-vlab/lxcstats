package pve

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testStorageCfg = `dir: local
	path /var/lib/vz
	content vztmpl,iso,images

lvmthin: local-lvm
	thinpool data
	vgname pve
	content rootdir,images

dir: nfs-template
	path /mnt/vz
	content vztmpl,iso,images
	shared 1
`

var testStorage = []PVEStorage{
	{
		Type: "dir",
		Name: "local",
		Attr: map[string]string{
			"path":    "/var/lib/vz",
			"content": "vztmpl,iso,images",
		},
	},
	{
		Type: "lvmthin",
		Name: "local-lvm",
		Attr: map[string]string{
			"thinpool": "data",
			"vgname":   "pve",
			"content":  "rootdir,images",
		},
	},
	{
		Type: "dir",
		Name: "nfs-template",
		Attr: map[string]string{
			"path":    "/mnt/vz",
			"content": "vztmpl,iso,images",
			"shared":  "1",
		},
	},
}

func TestParseStorage(t *testing.T) {
	as := assert.New(t)
	as.Equal(testStorage, parseStorage(strings.NewReader(testStorageCfg)))
}
