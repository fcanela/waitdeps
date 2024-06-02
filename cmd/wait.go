package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fcanela/waitdeps/internal/tcpcheck"
	"github.com/spf13/cobra"
)

var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Wait until the dependencies enumerated in the flags are available",
	Long: `Checks every provided service and quits when all the services are
available or after a timeout. To perform the same kind of check to multiple
dependencies, repeat the same flag or provide a comma separated list of them.

For example, to wait until db:5432 and 10.1.2.3:25 accepts TCP connections:
  waitdeps wait --tcp-connect db:5432 --tcp-connect 10.1.2.3:25

You can also use URI format. Users, passwords and paths are ignored, port is
automatically resolved to the default one when not provided:
  waitdeps wait --tcp-connect postgres://db --tcp-connect smtp://10.1.2.3

The timeout (or maximum wait time) is expressed as a number and a unit:
  waitdeps wait --timeout 500ms --tcp-connect postgres://db
  waitdeps wait --timeout 10s --tcp-connect postgres://db
  waitdeps wait --timeout 1m --tcp-connect postgres://db
  waitdeps wait --timeout 2h --tcp-connect postgres://db
`,
	Run: func(cmd *cobra.Command, args []string) {
		timeoutFlag, err := cmd.Flags().GetDuration("timeout")
		if err != nil {
			log.Fatalf("unable to parse --timeout flag: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeoutFlag)
		defer cancel()
		var wg sync.WaitGroup

		tcpConnectFlag, err := cmd.Flags().GetStringSlice("tcp-connect")
		if err != nil {
			fmt.Printf("unable to parse --tcp-connect flag: %v\n", err)
			os.Exit(1)
		}
		if len(tcpConnectFlag) > 0 {
			tcpcheck.Check(ctx, &wg, tcpConnectFlag)
		}

		wg.Wait()
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("some dependencies are still unavailable after the timeout")
			os.Exit(1)
		}

		fmt.Println("all dependencies are available")
	},
}

func init() {
	rootCmd.AddCommand(waitCmd)

	waitCmd.Flags().Duration("timeout", 30*time.Second, "Total wait time for all the dependencies to become available. Use number with units expressions like 200ms, 30s, 5m or 1h")
	waitCmd.Flags().StringSlice("tcp-connect", []string{}, "Performs a TCP connect check. You can use expressions like 127.0.0.1:22 or full URIs like https://news.ycombinator.com")
}
