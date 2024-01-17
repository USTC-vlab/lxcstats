package findpid

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func findLXCForPid(pid string) (string, error) {
	filename := fmt.Sprintf("/proc/%s/cgroup", pid)
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	if !s.Scan() {
		return "", fmt.Errorf("empty cgroup file %s", filename)
	}
	line := s.Text()
	remainder, ok := strings.CutPrefix(line, "0::/lxc/")
	if !ok {
		// does not belong to an LXC container
		return "", nil
	}
	parts := strings.SplitN(remainder, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid cgroup file %s (read %q)", filename, line)
	}
	return parts[0], nil
}

func runE(cmd *cobra.Command, args []string) error {
	w := cmd.OutOrStdout()
	for _, pid := range args {
		id, err := findLXCForPid(pid)
		if err != nil {
			return err
		}
		if id == "" {
			id = "<none>"
		}
		fmt.Fprintf(w, "%s: %s\n", pid, id)
	}
	return nil
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "findpid PID...",
		Short: "Find container ID by PID",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runE,
	}
	return cmd
}
