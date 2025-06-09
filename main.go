package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func createContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout >= 0 {
		return context.WithTimeout(ctx, timeout)
	}

	return context.WithCancel(ctx)
}

func main() {
	log.SetFlags(0)

	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", -1, "used to specify a timeout on the actual command")

	var graceful time.Duration
	flag.DurationVar(&graceful, "graceful", -1, "used to specify a graceful shutdown time")

	flag.Parse()

	cmdName := flag.Arg(0)
	args := flag.Args()[1:]

	ctx, cancel := createContext(context.Background(), timeout)
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	cmd := exec.CommandContext(ctx, cmdName, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Start()

	doneCh := make(chan error)
	go func() {
		doneCh <- cmd.Wait()
		close(doneCh)
	}()

	select {
	case err := <-doneCh:
		if err != nil {
			os.Exit(1)
		}
	case <-ch:
		break
	}

	cmd.Process.Signal(os.Interrupt)

	select {
	case <-doneCh:
		return

	case <-time.After(graceful):
		break
	}

	cmd.Process.Kill()
}
