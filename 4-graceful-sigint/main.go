//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"os"
	"os/signal"
)

func main() {
	// Create a process
	proc := MockProcess{}

	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		done := make(chan struct{})

		go func() {
			proc.Stop()
			close(done)
		}()

		select {
		case <-done:
			os.Exit(0)
		case <-sig:
			os.Exit(1)
		}
	}()

	// Run the process (blocking)
	proc.Run()
}
