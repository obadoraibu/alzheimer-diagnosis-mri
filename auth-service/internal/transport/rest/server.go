package rest

import (
	"golang.org/x/net/context"
	"net/http"
)

type URLShortenerServer struct {
	httpServer *http.Server
}

func NewServer() *URLShortenerServer {
	return new(URLShortenerServer)
}

func (s *URLShortenerServer) Start(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	return s.httpServer.ListenAndServe()
}

func (s *URLShortenerServer) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
