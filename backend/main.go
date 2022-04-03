package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/seth2810/go-collab/cmd"
)

func main() {
	ctx, cancelFn := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP,
	)

	defer cancelFn()

	cmd.ExecuteContext(ctx)
}
