package cleanup

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		for {
			s := <-signal_chan
			switch s {
			// kill -SIGHUP XXXX
			case syscall.SIGHUP:
				fmt.Println("hungup")

			// kill -SIGINT XXXX or Ctrl+c
			case syscall.SIGINT:
				fmt.Println("Warikomi")

			// kill -SIGTERM XXXX
			case syscall.SIGTERM:
				fmt.Println("force stop")

			// kill -SIGQUIT XXXX
			case syscall.SIGQUIT:
				fmt.Println("stop and core dump")

			default:
				fmt.Println("Unknown signal.")
			}

			cleanup()
			os.Exit(1)
		}
	}()
}
