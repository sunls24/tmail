package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tmail/internal/smtpd"
)

func main() {
	if err := run(); err != nil {
		slog.Error("SMTP server stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := smtpd.LoadConfig()
	if err != nil {
		return err
	}
	server, err := smtpd.NewServer(cfg)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		slog.Info("SMTP server started", "addr", cfg.Addr)
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
