package cctx

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func SignalContext(ctx context.Context, name string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("%s — listening for shutdown signal", name)
		<-sigs
		log.Printf("%s — shutdown signal received", name)
		signal.Stop(sigs)
		close(sigs)
		cancel()
	}()

	return ctx
}
