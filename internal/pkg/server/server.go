package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type Server struct {
	http.Server
}

func NewServer(addr string, router http.Handler, tlsConfig *tls.Config) *Server {
	server := Server{
		http.Server{
			Addr: addr,
			Handler: router,
			TLSConfig: tlsConfig,
		},
	}
	return &server
}

func (s *Server) ListenAndServe() {
	fmt.Printf("opening server on http://localhost%v\n", s.Addr)
	s.Server.ListenAndServe()
}