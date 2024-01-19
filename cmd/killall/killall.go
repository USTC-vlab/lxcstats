package killall

import (
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/USTC-vlab/vct/pkg/cgroup"
	"github.com/USTC-vlab/vct/pkg/pve"
	"github.com/spf13/cobra"
)

type idAndError struct {
	id  string
	err error
}

func killWorker(ch <-chan string, errCh chan<- idAndError, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := range ch {
		cmd := pve.Stop(id)
		if err := cmd.Run(); err != nil {
			errCh <- idAndError{id, err}
		}
	}
}

func killallMain(out io.Writer, n int, minID int) (err error) {
	ch := make(chan string)
	chErr := make(chan idAndError)
	wg := &sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go killWorker(ch, chErr, wg)
	}
	hasError := false
	go func() {
		for e := range chErr {
			hasError = true
			fmt.Fprintf(out, "error killing %s: %v\n", e.id, e.err)
		}
	}()
	defer func() {
		close(ch)
		wg.Wait()
		close(chErr)
		if hasError {
			err = fmt.Errorf("some containers failed to stop")
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
	return
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
		return killallMain(cmd.OutOrStderr(), *pN, *pS)
	}
	return cmd
}
