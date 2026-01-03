package http

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// Serve runs an HTTP server until the context is canceled.
func Serve(ctx context.Context, addr string, handler http.Handler) error {
	if handler == nil {
		return errors.New("handler is required")
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
