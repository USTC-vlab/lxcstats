package killall

import (
	"strconv"
	"sync"

	"github.com/USTC-vlab/vct/pkg/cgroup"
	"github.com/USTC-vlab/vct/pkg/pve"
	"github.com/spf13/cobra"
)

func killWorker(ch <-chan string, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := range ch {
		cmd := pve.Stop(id)
		if err := cmd.Run(); err != nil {
			errCh <- err
		}
	}
}

func killallMain(n int, minID int) error {
	ch := make(chan string)
	chErr := make(chan error)
	wg := &sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go killWorker(ch, chErr, wg)
	}
	defer func() {
		close(ch)
		wg.Wait()
		close(chErr)
	}()
	// discard errors for now
	go func() {
		for range chErr {
		}
	}()

	ids, err := cgroup.ListLXC()
	if err != nil {
		return err
	}
	for _, id := range ids {
		numID, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		if numID < minID {
			continue
		}
		ch <- id
	}
	return nil
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "killall",
		Short: "Kill all running containers",
		Args:  cobra.NoArgs,
	}
	flags := cmd.Flags()
	pN := flags.IntP("n", "n", 5, "max number of parallel killing containers")
	pS := flags.IntP("min", "m", 1000, "minimum ID of containers to kill")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return killallMain(*pN, *pS)
	}
	return cmd
}
