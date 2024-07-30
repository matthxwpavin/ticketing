package server

import (
	"context"
	"net/http"
	"time"

	"github.com/matthxwpavin/ticketing/logging/sugar"
)

type Server struct {
	Addr    string
	Handler http.Handler
}

func ListenAndServe(ctx context.Context, addr string, handler http.Handler) error {
	s := &Server{Addr: addr, Handler: handler}
	return s.ListenAndServe(ctx)
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	logger := sugar.FromContext(ctx)
	srv := &http.Server{
		Addr:         s.Addr,
		Handler:      s.Handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	go func() {
		<-ctx.Done()

		shutDownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutDownCtx); err != nil {
			logger.Errorw("server failed to shutdown", "error", err)
		} else {
			logger.Infoln("server is shutdown")
		}
	}()
	logger.Infow("server is listening", "address", s.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorw("server failed to listen", "error", err)
		return err
	}
	return nil
}
