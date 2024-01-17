package iostat

import "github.com/spf13/cobra"

func main(cmd *cobra.Command, args []string) {
	iostatMain(args[0])
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iostat disk",
		Short: "Show I/O statistics for a disk",
		Args:  cobra.ExactArgs(1),
		Run:   main,
	}
	return cmd
}
