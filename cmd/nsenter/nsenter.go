package nsenter

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/USTC-vlab/vct/pkg/cgroup"
	"github.com/spf13/cobra"
)

func nsenter(id string) error {
	f, err := cgroup.OpenLXC(id, "ns/init.scope/cgroup.procs")
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	if !s.Scan() {
		return fmt.Errorf("empty cgroup file %s", f.Name())
	}
	args := []string{"nsenter", "-a", "-t", s.Text()}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	return syscall.Exec(path, args, os.Environ())
}

func runE(cmd *cobra.Command, args []string) error {
	idStr := args[0]
	_, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	return nsenter(idStr)
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nsenter PID...",
		Short: "Enter container by ID",
		Args:  cobra.ExactArgs(1),
		RunE:  runE,
	}
	return cmd
}
